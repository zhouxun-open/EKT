package db

import (
	"encoding/hex"
	"sync"
)

type MemKVDatabase struct {
	Map sync.Map
}

func NewMemKVDatabase() *MemKVDatabase {
	return &MemKVDatabase{
		Map: sync.Map{},
	}
}

func (db *MemKVDatabase) Set(key, value []byte) error {
	db.Map.Store(hex.EncodeToString(key), value)
	return nil
}

func (db *MemKVDatabase) Get(key []byte) ([]byte, error) {
	result, exist := db.Map.Load(hex.EncodeToString(key))
	if !exist {
		return nil, NoSuchKeyError
	}
	value, ok := result.([]byte)
	if !ok {
		return nil, InvalidTypeError
	}
	return value, nil
}

func (db *MemKVDatabase) Delete(key []byte) error {
	db.Map.Delete(hex.EncodeToString(key))
	return nil
}
