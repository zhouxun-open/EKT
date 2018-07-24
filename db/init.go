package db

var EktDB IKVDatabase

func InitEKTDB(filePath string) {
	EktDB = NewComposedKVDatabase(filePath)
}

func GetDBInst() IKVDatabase {
	return EktDB
}
