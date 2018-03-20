package db

import (
	"fmt"
	"testing"
)

func TestLevelDB(t *testing.T) {
	db, err := NewLevelDB("testLevelDB")
	if err != nil {
		t.Fail()
	}
	key1 := []byte("HelloWorld")
	value1 := []byte("value1")
	//key2 := []byte("HelloX")
	//value2 := []byte("value2")
	err = db.Set(key1, value1)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	var bts []byte
	bts, err = db.Get(key1)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(string(bts))
}
