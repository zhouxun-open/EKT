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
	value1 := []byte("this is value1 ")
	key2 := []byte("HelloX")
	value2 := []byte("this is value 2")
	key3 := []byte("zhouxun")
	value3 := []byte("this is value 3")
	key4, value4 := []byte("HZhouxun"), []byte("this is value 4")
	key5, value5 := []byte("HZhouWorld"), []byte("this is value 5")
	key6, value6 := []byte("HelloWorld"), []byte("this is value 6")
	key7, value7 := []byte("HelloWorlef"), []byte("this is value 7")
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
	err = mtp.MustInsert(key1, value1)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	err = mtp.MustInsert(key2, value2)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	err = mtp.MustInsert(key3, value3)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	err = mtp.MustInsert(key5, value5)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	err = mtp.MustInsert(key4, value4)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	err = mtp.MustInsert(key6, value6)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	err = mtp.MustInsert(key7, value7)
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
