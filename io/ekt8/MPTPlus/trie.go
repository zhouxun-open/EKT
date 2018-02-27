package MPTPlus

import (
	"bytes"

	"github.com/EducationEKT/EKT/io/ekt8/db"
	"github.com/EducationEKT/EKT/io/ekt8/db/leveldb"
	"github.com/EducationEKT/EKT/io/ekt8/rlp"
)

var DB db.EKTDB

type TrieSonInfo struct {
	Hash []byte
	Tip  string
}

type TrieNode struct {
	Sons []TrieSonInfo
	Leaf bool
	Str  string
}

func init() {
	db, err := leveldb.NewLevelDB("testTrie")
	if err != nil {
		panic(err)
	} else {
		DB = db
	}
}

func Insert(root, value []byte) (newRoot []byte, err error) {
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

//func Get(root, hash string) string {
//	//TODO
//	return hash
//}

func Contains(root, hash string) {

}

func Delete(root, hash string) bool {
	//TODO
	return true
}
