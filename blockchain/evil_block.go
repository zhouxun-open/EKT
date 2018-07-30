package blockchain

import (
	"encoding/json"
)

type EvilBlock struct {
	Block1 *Block `json:"block1"`
	Block2 *Block `json:"block2"`
}

func NewEvilBlock(block1, block2 *Block) *EvilBlock {
	return &EvilBlock{
		Block1: block1,
		Block2: block2,
	}
}

func (evil EvilBlock) Bytes() []byte {
	data, _ := json.Marshal(evil)
	return data
}
