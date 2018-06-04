package blockchain

import (
	"encoding/hex"
	"sync"
)

var BlockRecorder *BlockRecord

func init() {
	BlockRecorder = NewBlockRecorder()
}

type BlockRecord struct {
	Status1    sync.Map
	Status     map[string]int
	Blocks     map[string]*Block
	Signatures map[string]string
}

func NewBlockRecorder() *BlockRecord {
	return &BlockRecord{
		Blocks:     make(map[string]*Block),
		Status:     make(map[string]int),
		Signatures: make(map[string]string),
	}
}

func GetBlockRecordInst() *BlockRecord {
	return BlockRecorder
}

func (record BlockRecord) Record(block *Block, sign []byte) {
	record.Blocks[hex.EncodeToString(block.Hash())] = block
	record.Signatures[hex.EncodeToString(block.Hash())] = hex.EncodeToString(sign)
}
