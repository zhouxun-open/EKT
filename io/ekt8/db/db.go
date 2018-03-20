package db

var EktDB *LevelDB

func InitEKTDB(filePath string) error {
	db, err := NewLevelDB(filePath)
	if err != nil {
		panic(err)
	}
	EktDB = db
	return err
}

func GetDBInst() *LevelDB {
	return EktDB
}
