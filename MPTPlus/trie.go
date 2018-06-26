package MPTPlus

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/EducationEKT/EKT/crypto"
	"github.com/EducationEKT/EKT/db"
	"github.com/EducationEKT/EKT/p2p"
)

var DB db.LevelDB

// Merkle Trie Plus树是一个安全的自校验的字典树的升级,每个节点都带有自己的路径值,叶子节点的
// 儿子节点存储的是Value的Hash值,根据Hash可以在levelDB上获取自己的Value
// key=strings.Join(pathValues, "") value=db.Get(leafNode.Sons[0].Hash)

// !importent if strings.Contains(string(key1), string(key2)) {
// 		panic("invalid key")
// }

func (this *MTP) GetInterfaceValue(key []byte, v interface{}) error {
	value, err := this.GetValue(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(value, v)
}

func (this *MTP) GetValue(key []byte) (value []byte, err error) {
	this.Lock.RLock()
	defer this.Lock.RUnlock()
	hash := this.Root
	left := key
	var vHash []byte
	find := false
	for {
		node, _ := this.GetNode(hash)
		if !node.Root {
			if PrefixLength(node.PathValue, left) == len(node.PathValue) {
				left = left[len(node.PathValue):]
			} else {
				return nil, errors.New("Not Exist")
			}
		}
		if len(left) == 0 {
			find = true
			vHash = node.Sons[0].Hash
			break
		}
		exist := false
		for _, son := range node.Sons {
			if PrefixLength(son.PathValue, left) > 0 {
				hash = son.Hash
				exist = true
				break
			}
		}
		if !exist {
			return nil, errors.New("Not Exist")
		}
	}
	if find {
		return this.DB.Get(vHash)
	}
	return nil, nil
}

func (this *MTP) ContainsKey(key []byte) bool {
	_, prefixs, err := this.FindParents(key)
	if err != nil {
		return false
	}
	prefix := bytes.Join(prefixs, nil)
	if bytes.Equal(prefix, key) {
		return true
	}
	return false
}

func (this *MTP) Update(key, value []byte) error {
	if !this.ContainsKey(key) {
		return errors.New("Not Exist")
	}
	parentHashes, _, err := this.FindParents(key)
	if err != nil {
		return nil
	}
	leafNode, err := this.GetNode(parentHashes[len(parentHashes)-1])
	valueHash, err := this.SaveValue(value)
	if err != nil {
		return err
	}
	leafNode.DeleteSon([]byte(""))
	leafNode.AddSon(valueHash, []byte(""))
	newHash, err := this.SaveNode(*leafNode)
	pathValue := leafNode.PathValue
	for i := len(parentHashes) - 2; i >= 0; i-- {
		node, _ := this.GetNode(parentHashes[i])
		node.DeleteSon(pathValue)
		node.AddSon(newHash, pathValue)
		newHash, _ = this.SaveNode(*node)
		pathValue = node.PathValue
	}
	rootNode, _ := this.GetNode(this.Root)
	rootNode.DeleteSon(pathValue)
	rootNode.AddSon(newHash, pathValue)
	this.Root, err = this.SaveNode(*rootNode)
	return err
}

func SyncDB(key []byte, peers p2p.Peers, leaf bool) {
	if _, err := db.GetDBInst().Get(key); err != nil {
		for _, peer := range peers {
			value, err := peer.GetDBValue(key)
			if err != nil {
				continue
			}
			if crypto.Validate(value, key) != nil {
				continue
			}
			db.GetDBInst().Set(key, value)
			if !leaf {
				var node TrieNode
				err = json.Unmarshal(value, &node)
				if err != nil {
					continue
				}
				if len(node.Sons) > 0 {
					for _, son := range node.Sons {
						go SyncDB(son.Hash, peers, node.Leaf)
					}
				}
			}
			break
		}
	}
}

/**
*把Key和Value插入到root对应的树上
*
*首先搜索到要插入的节点,插入之后向上回溯寻找自己的Parent节点更新,直至root节点
 */
func (this *MTP) MustInsert(key, value []byte) error {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	if this.ContainsKey(key) {
		return this.Update(key, value)
	}
	return this.TryInsert(key, value)
}

func (this *MTP) TryInsert(key, value []byte) error {
	hash, err := this.SaveValue(value)
	if err != nil {
		return err
	}
	parentHashes, prefixs, err := this.FindParents(key)
	if err != nil {
		return err
	}
	prefix := bytes.Join(prefixs, nil)
	leafNode := TrieNode{
		Sons:      []TrieSonInfo{TrieSonInfo{Hash: hash}},
		Leaf:      true,
		Root:      false,
		PathValue: key[len(prefix):],
	}
	leafNodeHash, err := this.SaveNode(leafNode)
	if err != nil {
		return err
	}

	oldPrefix_, newPrefix_, newHash_ := leafNode.PathValue, leafNode.PathValue, leafNodeHash
	for i := len(parentHashes) - 1; i >= 0; i-- {
		currentHash := parentHashes[i]
		currentNode, _ := this.GetNode(currentHash)
		if len(currentNode.PathValue) > len(prefixs[i]) {
			oldPrefix_ = currentNode.PathValue
			newNode := TrieNode{Root: false, Leaf: false, PathValue: prefixs[i]}
			currentNode.PathValue = currentNode.PathValue[len(prefixs[i]):]
			newCHash, _ := this.SaveNode(*currentNode)
			newNode.AddSon(newCHash, currentNode.PathValue)
			newNode.AddSon(newHash_, newPrefix_)
			newHash_, _ = this.SaveNode(newNode)
			newPrefix_ = newNode.PathValue
		} else {
			currentNode.DeleteSon(oldPrefix_)
			currentNode.AddSon(newHash_, newPrefix_)
			newHash_, _ = this.SaveNode(*currentNode)
			oldPrefix_, newPrefix_ = currentNode.PathValue, currentNode.PathValue
		}
	}

	rootNode, _ := this.GetNode(this.Root)
	rootNode.DeleteSon(oldPrefix_)
	rootNode.AddSon(newHash_, newPrefix_)
	this.Root, _ = this.SaveNode(*rootNode)

	return nil
}

func (this *MTP) FindParents(key []byte) (parentHashes [][]byte, prefixs [][]byte, err error) {
	left, currentHash := key, this.Root
	var node *TrieNode
	for {
		node, err = this.GetNode(currentHash)
		if nil != err || nil == node.Sons {
			return
		}
		exist := false
		for _, son := range node.Sons {
			if length := PrefixLength(left, son.PathValue); length > 0 {
				parentHashes = append(parentHashes, son.Hash)
				prefixs = append(prefixs, left[:length])
				left = left[length:]
				currentHash = son.Hash
				exist = true
				if length < len(son.PathValue) {
					return
				}
				break
			}
		}
		if !exist {
			return
		}
	}

	return
}

func (this *MTP) GetNode(hash []byte) (*TrieNode, error) {
	data, err := this.DB.Get(hash)
	if err != nil || len(data) == 0 {
		return nil, err
	}
	var node TrieNode
	err = json.Unmarshal(data, &node)
	return &node, err
}

func (this *MTP) SaveNode(node TrieNode) (nodeHash []byte, err error) {
	data, err := json.Marshal(node)
	if err != nil {
		return nil, err
	}
	return this.SaveValue(data)
}

func (this *MTP) SaveValue(value []byte) ([]byte, error) {
	hash := crypto.Sha3_256(value)
	return hash, this.DB.Set(hash, value)
}

//返回公共前缀的长度
func PrefixLength(a, b []byte) int {
	length := len(a)
	if len(b) < length {
		length = len(b)
	}
	i := 0
	for ; i < length; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return i
}
