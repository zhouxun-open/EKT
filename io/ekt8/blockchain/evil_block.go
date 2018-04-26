package blockchain

import (
	"encoding/hex"
	"encoding/json"
)

type EvilBlock struct {
	Block1          *Block `json:"block1"`
	Block2          *Block `json:"block2"`
	Block1Signature string `json:"signature1"`
	Block2Signature string `json:"signature2"`
}

func NewEvilBlock(block1, block2 *Block) *EvilBlock {
	return &EvilBlock{
		Block1:          block1,
		Block2:          block2,
		Block1Signature: GetBlockRecordInst().Signatures[hex.EncodeToString(block1.Hash())],
		Block2Signature: GetBlockRecordInst().Signatures[hex.EncodeToString(block2.Hash())],
	}
}

func (evil EvilBlock) Bytes() []byte {
	data, _ := json.Marshal(evil)
	return data
}
