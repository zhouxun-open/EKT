package crypto

import (
	"bytes"
	"errors"

	"github.com/ethereum/go-ethereum/crypto/sha3"
)

func Sha3_256(data []byte) []byte {
	result := sha3.Sum256(data)
	return result[:]
}

func Validate(data, hash []byte) error {
	if !bytes.Equal(Sha3_256(data), hash) {
		return errors.New("Error Hash")
	}
	return nil
}
