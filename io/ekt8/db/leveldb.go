package db

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDB struct {
	DB *leveldb.DB
}

func NewLevelDB(filePath string) (db *LevelDB, err error) {
	levelDB, err := leveldb.OpenFile(filePath, nil)
	if err != nil {
		return
	}
	db = &LevelDB{DB: levelDB}
	return
}

func (levelDB LevelDB) Set(key, value []byte) error {
	db := levelDB.DB
	return db.Put(key, value, nil)
}

func (levelDB LevelDB) Get(key []byte) ([]byte, error) {
	db := levelDB.DB
	return db.Get(key, nil)
}

func (levelDB LevelDB) Delete(key []byte) error {
	return levelDB.DB.Delete(key, nil)
}
