package leveldb

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDB struct {
	DB *leveldb.DB
}

func NewLevelDB(filePath string) (db LevelDB, err error) {
	levelDB, err := leveldb.OpenFile(filePath, nil)
	if err != nil {
		return
	}
	db = LevelDB{levelDB}
	return
}

func (levelDB LevelDB) Set(key, value []byte) error {
	return levelDB.DB.Put(key, value, nil)
}

func (levelDB LevelDB) Get(key, value []byte) error {
	value, err := levelDB.DB.Get(key, nil)
	return err
}

func (levelDB LevelDB) Delete(key []byte) error {
	return levelDB.DB.Delete(key, nil)
}
