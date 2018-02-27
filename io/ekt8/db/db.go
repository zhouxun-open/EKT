package db

type EKTDB interface {
	Set(key, value []byte) error
	Get(Key, value []byte) error
	Delete(Key []byte) error
}

func init() {
}
