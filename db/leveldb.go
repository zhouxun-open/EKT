package db

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDB struct {
	DB *leveldb.DB
}

func NewLevelDB(filePath string) *LevelDB {
	levelDB, err := leveldb.OpenFile(filePath, nil)
	if err != nil {
		panic(err)
	}
	return &LevelDB{DB: levelDB}
}

func (levelDB LevelDB) Set(key, value []byte) error {
	return levelDB.DB.Put(key, value, nil)
}

func (levelDB LevelDB) Get(key []byte) ([]byte, error) {
	return levelDB.DB.Get(key, nil)
}

func (levelDB LevelDB) Delete(key []byte) error {
	return levelDB.DB.Delete(key, nil)
}
