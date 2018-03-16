package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/EducationEKT/EKT/io/ekt8/MPTPlus"
	"github.com/EducationEKT/EKT/io/ekt8/consensus"
	"github.com/EducationEKT/EKT/io/ekt8/db"
)

var BackboneChainId []byte = [32]byte{31: byte(1 & 0xFF)}[:]

const (
	CurrentBlockKey    = "CurrentBlock"
	BackboneConsensus  = consensus.DPOS
	InitStatus         = 0
	OpenStatus         = 100
	CaculateHashStatus = 150
)

type BlockChain struct {
	ChainId    []byte
	Consensus  consensus.ConsensusType
	Locker     sync.RWMutex
	status     int // 100 正在计算MTProot, 150停止计算root,开始计算block Hash
	Fee        int64
	Difficulty []byte
	consensus  consensus.Consensus
}

func (this *BlockChain) SyncBlockChain() error {
	this.Locker.Lock()
	defer this.Locker.Unlock()
	this.status = OpenStatus
	return nil
}

func (this *BlockChain) GetStatus() int {
	this.Locker.RLock()
	defer this.Locker.RUnlock()
	return this.status
}

func (this *BlockChain) NewBlock(block Block) error {
	this.Locker.Lock()
	defer this.Locker.Unlock()
	if err := block.Validate(); err != nil {
		return err
	}
	db.GetDBInst().Set(block.Hash(), block.Bytes())
	newBlock := &Block{
		Height:       block.Height + 1,
		Nonce:        0,
		Fee:          this.Fee,
		PreviousHash: block.Hash(),
		Locker:       sync.RWMutex{},
		StatTree:     MPTPlus.MTP_Tree(db.GetDBInst(), block.StatTree.Root),
		TxTree:       MPTPlus.NewMTP(db.GetDBInst()),
		EventTree:    MPTPlus.NewMTP(db.GetDBInst()),
	}
	value, _ := json.Marshal(newBlock)
	return db.GetDBInst().Set(this.CurrentBlockKey(), value)
	//lastBlock, err := this.CurrentBlock()
	//if err != nil {
	//	return err
	//}
	//if lastBlock.Height > block.Height {
	//	return errors.New("height exist")
	//}
	//err = db.GetDBInst().Set(block.CurrentHash, block.Hash())
	//if err != nil {
	//	return err
	//}
	//value, err := json.Marshal(block)
	//if err != nil {
	//	return err
	//}
	//return db.GetDBInst().Set(this.CurrentBlockKey(), value)
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
	buffer.WriteString(CurrentBlockKey)
	buffer.Write(this.ChainId)
	return buffer.Bytes()
}

func (this BlockChain) Pack() {
	block, _ := this.CurrentBlock()
	block.Locker.Lock()
	start := time.Now().Nanosecond()
	for ; !bytes.HasPrefix(block.Hash(), []byte("FFFFFF")); block.NewNonce() {
	}
	end := time.Now().Nanosecond()
	fmt.Printf(`difficulty="FFFFFF", cost=%d`, (end-start)/1e6)
	//db.GetDBInst().Set(block.Hash(), block.Bytes())
	this.consensus.NewBlock(*block, ConsensusCb)
}

func ConsensusCb(blockChain BlockChain, block Block, result bool) {
	if result {
		blockChain.NewBlock(block)
	}
}
