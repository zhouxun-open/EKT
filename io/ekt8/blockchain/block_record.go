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
	status *sync.Map
	blocks *sync.Map
}

func NewBlockRecorder() *BlockRecord {
	return &BlockRecord{
		status: &sync.Map{},
		blocks: &sync.Map{},
	}
}

func (recorder BlockRecord) GetStatus(blockHash string) int {
	obj, exist := recorder.status.Load(blockHash)
	if !exist {
		return -1
	}
	return obj.(int)
}

func (recorder BlockRecord) SetStatus(blockHash string, status int) {
	recorder.status.Store(blockHash, status)
}

func (recorder BlockRecord) GetBlock(hash string) *Block {
	obj, exist := recorder.blocks.Load(hash)
	if !exist {
		return nil
	}
	return obj.(*Block)
}

func (recorder BlockRecord) SetBlock(block *Block) {
	recorder.blocks.Store(hex.EncodeToString(block.CurrentHash), block)
}
