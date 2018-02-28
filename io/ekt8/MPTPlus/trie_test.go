package MPTPlus

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/EducationEKT/EKT/io/ekt8/db/leveldb"
)

func TestMTP(t *testing.T) {
	db, err := leveldb.NewLevelDB("testTrie11")
	defer db.DB.Close()
	key1 := []byte("HelloWorld")
	value1 := []byte("value1")
	key2 := []byte("HelloX")
	value2 := []byte("value2")
	if err != nil {
		fmt.Println("===", err)
		t.Fail()
	}
	node := TrieNode{
		Root:      true,
		Leaf:      false,
		PathValue: nil,
		Sons:      *new(SortedSon),
	}
	mtp := &MTP{DB: db, Root: nil}
	mtp.Root, err = mtp.SaveNode(node)
	fmt.Println(hex.EncodeToString(mtp.Root))
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	node1, err := mtp.GetNode(mtp.Root)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Println(node1)
	err = mtp.TryInsert(key1, value1)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	err = mtp.TryInsert(key2, value2)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(hex.EncodeToString(mtp.Root))
	value, err := mtp.GetValue(key1)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(string(value))
}
