package MPTPlus

import (
	"bytes"

	"github.com/EducationEKT/EKT/io/ekt8/db"
	"github.com/EducationEKT/EKT/io/ekt8/db/leveldb"
	"github.com/EducationEKT/EKT/io/ekt8/rlp"
)

var DB db.EKTDB

// Merkle Trie Plus树是一个安全的自校验的字典树的升级,每个节点都带有自己的路径值,叶子节点的
// 儿子节点存储的是Value的Hash值,根据Hash可以在levelDB上获取自己的Value
// key=strings.Join(pathValues, "") value=db.Get(leafNode.Sons[0].Hash)

/*
*TrieSonInfo存储的是当前节点的下一个节点的信息
*如果当前节点是叶子节点,则Sons的长度为1,且TrieSonInfo.Hash是当前路径key的结果值的hash
*如果当前节点不是叶子节点,则Sons的长度大于等于1,存储的是子节点的Hash值和PathValue
*/
type TrieSonInfo struct {
	Hash      []byte
	PathValue string
}

/*
*TrieNode存储的是当前节点的一些信息,包括PathValue 是否是叶子节点,子节点信息等等
*strings.Join(pathValue,"")就是用户要存储的key
*/
type TrieNode struct {
	Sons      []TrieSonInfo
	Leaf      bool
	PathValue string
}

func init() {
	db, err := leveldb.NewLevelDB("testTrie")
	if err != nil {
		panic(err)
	} else {
		DB = db
	}
}

/**
*把Key和Value插入到root对应的树上
*
*首先搜索到要插入的节点,插入之后向上回溯寻找自己的Parent节点更新,直至root节点
*/
func Insert(root, key, value []byte) (newRoot []byte, err error) {
	newRoot = root
	var buffer bytes.Buffer
	err = DB.Get(root, buffer.Bytes())
	if err != nil {
		return nil, err
	} else {
		if nil != buffer.Bytes() && len(buffer.Bytes()) > 0 {
			var node TrieNode
			err = rlp.Decode(buffer.Bytes(), &node)
			if err != nil {
				return
			} else {

			}
		}
	}
	//TODO
	return
}

func Get(root []byte, key string) (value []byte, exist bool ){
	return nil, false
}

func Delete(root, hash string) bool {
	//TODO
	return true
}
