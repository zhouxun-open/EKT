package db

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
)

type EKTDB interface {
	Set(key, value []byte)
	Get(Key, value []byte)
	Delete(Key []byte)
}

func init() {
	db, _ := leveldb.OpenFile("db.leveldb", nil)
	defer db.Close()
	db.Put([]byte("Hello"), []byte("World"), nil)
	bts, _ := db.Get([]byte("Hello"), nil)
	fmt.Println(string(bts))
	fmt.Println(db.Has([]byte("Hello"), nil))
	db.Delete([]byte("Hello"), nil)
	bts, _ = db.Get([]byte("Hello"), nil)
	fmt.Println("=====", string(bts))
	if err := db.Delete([]byte("Hello"), nil); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(db.Has([]byte("Hello"), nil))
}
