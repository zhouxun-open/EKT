package blockchain

import (
	"bytes"
	"encoding/json"
	"errors"
	"sync"

	"github.com/EducationEKT/EKT/io/ekt8/consensus"
	"github.com/EducationEKT/EKT/io/ekt8/db"
)

var BackboneChainId []byte = [32]byte{31: byte(1 & 0xFF)}[:]

const (
	CurrentBlock       = "CurrentBlock"
	BackboneConsensus  = consensus.DPOS
	InitStatus         = 0
	OpenStatus         = 100
	CaculateHashStatus = 150
)

type BlockChain struct {
	ChainId   []byte
	Consensus consensus.ConsensusType
	Locker    sync.RWMutex
	Status    int // 100 正在计算MTProot, 150停止计算root,开始计算block Hash
}

func (this *BlockChain) SyncBlockChain() error {
	this.Locker.Lock()
	defer this.Locker.Unlock()
	this.Status = OpenStatus
	return nil
}

func (this *BlockChain) NewBlock(block Block) error {
	this.Locker.Lock()
	defer this.Locker.Unlock()
	if err := block.Validate(); err != nil {
		return err
	}
	lastBlock, err := this.CurrentBlock()
	if err != nil {
		return err
	}
	if lastBlock.Height > block.Height {
		return errors.New("height exist")
	}
	err = db.GetDBInst().Set(block.CurrentHash, block.Hash())
	if err != nil {
		return err
	}
	value, err := json.Marshal(block)
	if err != nil {
		return err
	}
	return db.GetDBInst().Set(this.CurrentBlockKey(), value)
}

func (this BlockChain) CurrentBlock() (*Block, error) {
	var err error = nil
	if currentBlock == nil {
		blockValue, err := db.GetDBInst().Get(this.CurrentBlockKey())
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(blockValue, &currentBlock)
	}
	return currentBlock, err
}

func (this BlockChain) CurrentBlockKey() []byte {
	buffer := bytes.Buffer{}
	buffer.WriteString(CurrentBlock)
	buffer.Write(this.ChainId)
	return buffer.Bytes()
}
