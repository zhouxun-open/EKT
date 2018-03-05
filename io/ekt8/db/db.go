package db

import "github.com/EducationEKT/EKT/io/ekt8/db/leveldb"

type EKTDB interface {
	Set(key, value []byte) error
	Get(Key []byte) ([]byte, error)
	Delete(Key []byte) error
}

var DB EKTDB

func InitEKTDB(filePath string) error {
	db, err := leveldb.NewLevelDB(filePath)
	if err != nil {
		return err
	}
	DB = db
	return nil
}
