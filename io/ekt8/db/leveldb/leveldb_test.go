package leveldb

import (
	"fmt"
	"testing"
)

func TestLevelDB(t *testing.T) {
	db, err := NewLevelDB("testLevelDB")
	if err != nil {
		t.Fail()
	}
	var bts []byte
	err = db.Get([]byte("Hello"), bts)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(string(bts))
}
