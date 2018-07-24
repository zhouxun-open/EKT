package MPTPlus

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/EducationEKT/EKT/crypto"
)

// Merkle Trie Plus树是一个安全的自校验的字典树的升级,每个节点都带有自己的路径值,叶子节点的
// 儿子节点存储的是Value的Hash值,根据Hash可以在levelDB上获取自己的Value
// key=strings.Join(pathValues, "") value=db.Get(leafNode.Sons[0].Hash)

// !importent if strings.Contains(string(key1), string(key2)) {
// 		panic("invalid key")
// }

func (mtp *MTP) GetInterfaceValue(key []byte, v interface{}) error {
	value, err := mtp.GetValue(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(value, v)
}

func (mtp *MTP) GetValue(key []byte) (value []byte, err error) {
	return mtp.GetValueByKey(key)
}

func (mtp *MTP) GetValueByKey(key []byte) (value []byte, err error) {
	mtp.Lock.RLock()
	defer mtp.Lock.RUnlock()
	parentHashes, prefixs, err := mtp.FindParents(key)
	if err != nil {
		return nil, err
	}
	if bytes.EqualFold(bytes.Join(prefixs, nil), key) {
		leaf, err := mtp.GetNode(parentHashes[len(parentHashes)-1])
		if err != nil {
			return nil, err
		} else {
			return mtp.DB.Get(leaf.Sons[0].Hash)
		}
	} else {
		return nil, errors.New("key not exist")
	}
}

func (mtp *MTP) ContainsKey(key []byte) bool {
	_, prefixs, err := mtp.FindParents(key)
	if err != nil {
		return false
	}
	prefix := bytes.Join(prefixs, nil)
	if bytes.EqualFold(prefix, key) {
		return true
	}
	return false
}

func (mtp *MTP) Update(key, value []byte, parentHashes [][]byte, prefixs [][]byte) error {
	leafNode, err := mtp.GetNode(parentHashes[len(parentHashes)-1])
	valueHash, err := mtp.SaveValue(value)
	if err != nil {
		return err
	}
	leafNode.DeleteSon(nil)
	leafNode.AddSon(valueHash, nil)
	newHash, err := mtp.SaveNode(*leafNode)
	pathValue := leafNode.PathValue
	for i := len(parentHashes) - 2; i >= 0; i-- {
		node, _ := mtp.GetNode(parentHashes[i])
		node.DeleteSon(pathValue)
		node.AddSon(newHash, pathValue)
		newHash, _ = mtp.SaveNode(*node)
		pathValue = node.PathValue
	}
	rootNode, _ := mtp.GetNode(mtp.Root)
	rootNode.DeleteSon(pathValue)
	rootNode.AddSon(newHash, pathValue)
	mtp.Root, err = mtp.SaveNode(*rootNode)
	return err
}

/**
*把Key和Value插入到root对应的树上
*
*首先搜索到要插入的节点,插入之后向上回溯寻找自己的Parent节点更新,直至root节点
 */
func (mtp *MTP) MustInsert(key, value []byte) error {
	mtp.Lock.Lock()
	defer mtp.Lock.Unlock()
	// 遍历字典树
	parentHashes, prefixs, err := mtp.FindParents(key)
	if err != nil {
		return err
	}
	// 包含key， 更新value
	if bytes.EqualFold(bytes.Join(prefixs, nil), key) {
		return mtp.Update(key, value, parentHashes, prefixs)
	} else {
		// 重新插入key和value
		return mtp.TryInsert(key, value, parentHashes, prefixs)
	}
}

func (mtp *MTP) TryInsert(key, value []byte, parentHashes, prefixs [][]byte) error {
	hash, err := mtp.SaveValue(value)
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
	leafNodeHash, err := mtp.SaveNode(leafNode)
	if err != nil {
		return err
	}

	oldPrefix_, newPrefix_, newHash_ := leafNode.PathValue, leafNode.PathValue, leafNodeHash
	for i := len(parentHashes) - 1; i >= 0; i-- {
		currentHash := parentHashes[i]
		currentNode, _ := mtp.GetNode(currentHash)
		if len(currentNode.PathValue) > len(prefixs[i]) {
			oldPrefix_ = currentNode.PathValue
			newNode := TrieNode{Root: false, Leaf: false, PathValue: prefixs[i]}
			currentNode.PathValue = currentNode.PathValue[len(prefixs[i]):]
			newCHash, _ := mtp.SaveNode(*currentNode)
			newNode.AddSon(newCHash, currentNode.PathValue)
			newNode.AddSon(newHash_, newPrefix_)
			newHash_, _ = mtp.SaveNode(newNode)
			newPrefix_ = newNode.PathValue
		} else {
			currentNode.DeleteSon(oldPrefix_)
			currentNode.AddSon(newHash_, newPrefix_)
			newHash_, _ = mtp.SaveNode(*currentNode)
			oldPrefix_, newPrefix_ = currentNode.PathValue, currentNode.PathValue
		}
	}

	rootNode, _ := mtp.GetNode(mtp.Root)
	rootNode.DeleteSon(oldPrefix_)
	rootNode.AddSon(newHash_, newPrefix_)
	mtp.Root, _ = mtp.SaveNode(*rootNode)

	return nil
}

func (mtp *MTP) FindParents(key []byte) (parentHashes [][]byte, prefixs [][]byte, err error) {
	left, currentHash := key, mtp.Root
	var node *TrieNode
	for {
		node, err = mtp.GetNode(currentHash)
		if nil != err || nil == node.Sons {
			return
		}
		exist := false
		for _, son := range node.Sons {
			if length := PrefixLength(left, son.PathValue); length > 0 {
				parentHashes = append(parentHashes, son.Hash)
				prefixs = append(prefixs, left[:length])
				if len(son.PathValue) > length || len(left) == length {
					return
				}
				left = left[length:]
				currentHash = son.Hash
				exist = true
				break
			}
		}
		if !exist {
			return
		}
	}

	return
}

func (mtp *MTP) GetNode(hash []byte) (*TrieNode, error) {
	data, err := mtp.DB.Get(hash)
	if err != nil || len(data) == 0 {
		return nil, err
	}
	var node TrieNode
	err = json.Unmarshal(data, &node)
	return &node, err
}

func (mtp *MTP) SaveNode(node TrieNode) (nodeHash []byte, err error) {
	data, err := json.Marshal(node)
	if err != nil {
		return nil, err
	}
	return mtp.SaveValue(data)
}

func (mtp *MTP) SaveValue(value []byte) ([]byte, error) {
	hash := crypto.Sha3_256(value)
	return hash, mtp.DB.Set(hash, value)
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
