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
	Locker     sync.RWMutex
	Status     map[string]int
	Blocks     map[string]*Block
	Signatures map[string]string
}

func NewBlockRecorder() *BlockRecord {
	return &BlockRecord{
		Locker:     sync.RWMutex{},
		Blocks:     make(map[string]*Block),
		Status:     make(map[string]int),
		Signatures: make(map[string]string),
	}
}

func (recorder BlockRecord) GetStatus(blockHash string) int {
	recorder.Locker.RLock()
	defer recorder.Locker.RUnlock()
	if status, exist := recorder.Status[blockHash]; exist {
		return status
	}
	return -1
}

func (recorder BlockRecord) SetStatus(blockHash string, status int) {
	recorder.Locker.Lock()
	defer recorder.Locker.Unlock()
	recorder.Status[blockHash] = status
}

func GetBlockRecordInst() *BlockRecord {
	return BlockRecorder
}

func (record BlockRecord) Record(block *Block, sign []byte) {
	record.Blocks[hex.EncodeToString(block.Hash())] = block
	record.Signatures[hex.EncodeToString(block.Hash())] = hex.EncodeToString(sign)
}
