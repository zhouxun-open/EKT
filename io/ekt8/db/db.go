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
	
}
