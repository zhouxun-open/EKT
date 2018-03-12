package db

import "github.com/EducationEKT/EKT/io/ekt8/db/leveldb"

type EKTDB interface {
	Set(key, value []byte) error
	Get(Key []byte) ([]byte, error)
	Delete(Key []byte) error
}

var ektDB EKTDB

func InitEKTDB(filePath string) error {
	db, err := leveldb.NewLevelDB(filePath)
	if err != nil {
		return err
	}
	ektDB = db
	return nil
}

func GetDBInst() EKTDB {
	return ektDB
}
