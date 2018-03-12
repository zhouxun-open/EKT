package crypto

import (
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

func Sha3_256(data []byte) []byte {
	result := sha3.Sum256(data)
	return result[:]
}
