package cache

import "github.com/EducationEKT/EKT/io/ekt8/db"

var (
	CurrentBlockKey = []byte("MainBlockchainCurrentBlockKey")
	InitedKey       = []byte("inited")
)

func CurrentBlock() ([]byte, error) {
	return db.GetDBInst().Get(CurrentBlockKey)
}

func SetCurrentBlock(hash []byte) {
	db.GetDBInst().Set(CurrentBlockKey, hash)
}

func IsInited() bool {
	v, err := db.GetDBInst().Get(InitedKey)
	if err != nil || len(v) == 0 {
		return false
	}
	return true
}

func GenesisBlockInit() {
	db.GetDBInst().Set(InitedKey, InitedKey)
}
