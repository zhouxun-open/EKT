package blockchain

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/EducationEKT/EKT/io/ekt8/crypto"
)

type Block struct {
	Height       int64  `json:"height"`
	StatRoot     []byte `json:"statRoot"`
	TxRoot       []byte `json:"txRoot"`
	EventRoot    []byte `json:"eventRoot"`
	Nonce        int64  `json:"nonce"`
	PreviousHash []byte `json:"previousHash"`
	CurrentHash  []byte `json:"currentHash"`
}

func (block *Block) String() string {
	return fmt.Sprintf(`{"height": %d, "statRoot": "%s", "txRoot": "%s", "eventRoot": "%s", "nonce": %d, "previousHash": "%s"}`,
		block.Height, block.StatRoot, block.TxRoot, block.Nonce, block.PreviousHash)
}

func (block *Block) Bytes() []byte {
	return []byte(block.String())
}

func (block *Block) Hash() []byte {
	return crypto.Sha3_256(block.Hash())
}

func (block *Block) NewNonce() {
	block.Nonce++
}

func (block Block) Validate() error {
	if !bytes.Equal(block.Hash(), block.CurrentHash) {
		return errors.New("Invalid Hash")
	}
	return nil
}
