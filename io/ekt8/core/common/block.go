package common

import (
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/sha3"
)

type Block struct {
	Height          int64
	PreviousHash    Hash
	CurrentHash     Hash
	Nonce           int64
	StatRoot        []byte
	TransactionRoot []byte
}

func (block *Block) Hash() {
	hash := sha3.New256()
	result := hash.Sum(nil)
	fmt.Println(string(result))
}

func (block *Block) Bytes() []byte {
	return []byte(block.ToString())
}

func (block *Block) ToString() string {
	return fmt.Sprintf(`{"height": %d, "previousHash": "%s", "statRoot": "%s", "transactionRoot": "%s", "nonce": %d}`,
		block.Height, hex.EncodeToString(block.PreviousHash[:]),
		hex.EncodeToString(block.StatRoot),
		hex.EncodeToString(block.TransactionRoot))
}
