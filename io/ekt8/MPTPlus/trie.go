package MPTPlus

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/EducationEKT/EKT/io/ekt8/crypto"
	"github.com/EducationEKT/EKT/io/ekt8/db"
	"github.com/EducationEKT/EKT/io/ekt8/rlp"
)

var DB db.EKTDB

// Merkle Trie Plus树是一个安全的自校验的字典树的升级,每个节点都带有自己的路径值,叶子节点的
// 儿子节点存储的是Value的Hash值,根据Hash可以在levelDB上获取自己的Value
// key=strings.Join(pathValues, "") value=db.Get(leafNode.Sons[0].Hash)

func init() {
}

/**
*把Key和Value插入到root对应的树上
*
*首先搜索到要插入的节点,插入之后向上回溯寻找自己的Parent节点更新,直至root节点
 */
func (this *MTP) Insert(key, value []byte) (err error) {
	var left, prefix, searchRoot []byte
	finish := false
	var parents [][]byte
	searchRoot = this.Root
	parents = append(parents, searchRoot)
	for !finish && err == nil {
		finish, prefix, searchRoot, err = Find(searchRoot, prefix, left, this.DB)
		parents = append(parents, searchRoot)
	}
	if err != nil {
		return
	}
	hash, err := this.SaveValue(value)
	if err != nil {
		return
	}
	leafNode := TrieNode{
		Sons:      []TrieSonInfo{TrieSonInfo{Hash: hash}},
		Leaf:      true,
		Root:      false,
		PathValue: key[len(prefix):],
	}
	leafNodeData, err := rlp.Encode(leafNode)
	if err != nil {
		return
	}
	leafNodeHash, err := this.SaveValue(leafNodeData)
	if len(parents) > 0 {
		lastHash := parents[len(parents)-1]
		lastNode, err1 := this.GetNode(lastHash)
		if err1 != nil {
			return err1
		}
		lastNode.AddSon(leafNodeHash, key)
	}
	//TODO
	return
}

func (this *MTP) GetValue(key []byte) (value []byte, err error) {
	hash := this.Root
	left := key
	for {
		node, _ := this.GetNode(hash)
		if !node.Root {
			if PrefixLength(node.PathValue, left) == len(node.PathValue) {
				left = left[len(node.PathValue):]
			} else {
				return nil, errors.New("Not Exist")
			}
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
	return nil, nil
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

	if nil == parentHashes || 0 == len(parentHashes) {
		rootNode, err := this.GetNode(this.Root)
		rootNode.AddSon(leafNodeHash, key)
		root, err := this.SaveNode(*rootNode)
		fmt.Println(hex.EncodeToString(root))
		if err == nil {
			this.Root = root
		}
		return err
	}

	parentIndex := len(parentHashes) - 1
	parentHash := parentHashes[parentIndex]
	parentNode, err := this.GetNode(parentHash)
	if err != nil {
		return nil
	}

	var newNodeHash, oldPrefix, newPrefix []byte
	if len(parentNode.PathValue) > len(prefixs[parentIndex]) {
		var newNode TrieNode
		newNode.Root = false
		newNode.Leaf = false
		parentNode.PathValue = parentNode.PathValue[len(prefixs[parentIndex]):]
		nodeValue, err := rlp.Encode(parentNode)
		if err != nil {
			return err
		}
		oldHash, err := this.SaveValue(nodeValue)
		if err != nil {
			return err
		}
		newNode.PathValue = prefixs[parentIndex]

		newNode.AddSon(leafNodeHash, key[len(prefix):])
		newNode.AddSon(oldHash, parentNode.PathValue)
		newNodeHash, _ = this.SaveNode(newNode)
		oldPrefix = parentNode.PathValue
		newPrefix = newNode.PathValue
	} else {
		parentNode.AddSon(leafNodeHash, key[len(prefix):])
		newNodeHash, _ = this.SaveNode(*parentNode)
		oldPrefix = parentNode.PathValue
		newPrefix = parentNode.PathValue
	}
	for i := len(parentHashes) - 2; i >= 0; i++ {
		node, _ := this.GetNode(parentHashes[i])
		node.DeleteSon(oldPrefix)
		node.AddSon(newNodeHash, newPrefix)
		oldPrefix = node.PathValue
		newPrefix = node.PathValue
		newNodeHash, _ = this.SaveNode(*node)
	}
	this.Root = newNodeHash
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
		for _, son := range node.Sons {
			if length := PrefixLength(left, son.PathValue); length > 0 {
				parentHashes = append(parentHashes, son.Hash)
				prefixs = append(prefixs, left[:length])
				left = left[length:]
				currentHash = son.Hash
				if length < len(son.PathValue) {
					return
				}
				break
			}
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

func Find(root, prefix, left []byte, db db.EKTDB) (finish bool, newPrefix []byte, nextRoot []byte, err error) {
	value, err := db.Get(root)
	if err != nil || len(value) == 0 {
		return
	}
	var node TrieNode
	err = rlp.Decode(value, &node)
	if err != nil {
		return
	}
	if node.Root {
		if len(node.Sons) > 0 {
			for _, sonNode := range node.Sons {
				if PrefixLength(sonNode.PathValue, left) > 0 {
					nextRoot = sonNode.Hash
					break
				}
			}
		}
		return
	}
	prefixLength := PrefixLength(node.PathValue, left)
	buffer := bytes.Buffer{}
	buffer.Write(prefix)
	buffer.Write(node.PathValue[:prefixLength])
	if prefixLength < len(node.PathValue) {
		finish = true
		nextRoot = []byte("")
	} else {
		if len(node.Sons) > 0 {
			for _, sonNode := range node.Sons {
				if PrefixLength(sonNode.PathValue, left[prefixLength:]) > 0 {
					nextRoot = sonNode.Hash
					break
				}
			}
		}
	}
	newPrefix = buffer.Bytes()
	return
}

func Get(root []byte, key string) (value []byte, exist bool) {
	return nil, false
}

func Delete(root, hash string) bool {
	//TODO
	return true
}

//返回公共前缀的长度
func PrefixLength(a, b []byte) int {
	length := len(a)
	if len(b) < length {
		length = len(b)
	}
	for i := 0; i < length; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return 0
}
