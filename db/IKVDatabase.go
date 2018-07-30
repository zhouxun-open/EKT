package db

type IKVDatabase interface {
	Set(key, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
}
