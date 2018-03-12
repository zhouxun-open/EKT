package rlp

import (
	"bytes"

	"github.com/ethereum/go-ethereum/rlp"
)

func Encode(value interface{}) ([]byte, error) {
	w := new(bytes.Buffer)
	err := rlp.Encode(w, value)
	return w.Bytes(), err
}

func Decode(data []byte, value interface{}) error {
	return rlp.DecodeBytes(data, value)
}
