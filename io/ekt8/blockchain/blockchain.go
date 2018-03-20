package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/EducationEKT/EKT/io/ekt8/MPTPlus"
	"github.com/EducationEKT/EKT/io/ekt8/blockchain_manager"
	"github.com/EducationEKT/EKT/io/ekt8/consensus"
	"github.com/EducationEKT/EKT/io/ekt8/db"
	"github.com/EducationEKT/EKT/io/ekt8/p2p"
)

var BackboneChainId []byte = [32]byte{31: byte(1 & 0xFF)}[:]

const (
	CurrentBlockKey       = "CurrentBlock"
	BackboneConsensus     = consensus.DPOS
	BackboneBlockInterval = 3
	InitStatus            = 0
	OpenStatus            = 100
	CaculateHashStatus    = 150
)

type BlockChain struct {
	ChainId    []byte
	Consensus  consensus.ConsensusType
	Locker     sync.RWMutex
	status     int // 100 正在计算MTProot, 150停止计算root,开始计算block Hash
	Fee        int64
	Difficulty []byte
}

func (blockchain *BlockChain) GenesisBlock(peers p2p.Peers) {}

func (blockchain *BlockChain) SyncBlockChain() error {
	blockchain.Locker.Lock()
	defer blockchain.Locker.Unlock()
	blockchain.status = OpenStatus
	return nil
}

func (blockchain *BlockChain) GetStatus() int {
	blockchain.Locker.RLock()
	defer blockchain.Locker.RUnlock()
	return blockchain.status
}

func (blockchain *BlockChain) NewBlock(block Block) error {
	blockchain.Locker.Lock()
	defer blockchain.Locker.Unlock()
	if err := block.Validate(); err != nil {
		return err
	}
	db.GetDBInst().Set(block.Hash(), block.Bytes())
	newBlock := &Block{
		Height:       block.Height + 1,
		Nonce:        0,
		Fee:          blockchain.Fee,
		TotalFee:     0,
		PreviousHash: block.Hash(),
		Locker:       sync.RWMutex{},
		StatTree:     MPTPlus.MTP_Tree(db.GetDBInst(), block.StatTree.Root),
		TxTree:       MPTPlus.NewMTP(db.GetDBInst()),
		EventTree:    MPTPlus.NewMTP(db.GetDBInst()),
		Round:        consensus.NextRound(block.Round, block.Hash()),
	}
	newBlock.UpdateMPTPlusRoot()
	block.UpdateMPTPlusRoot()
	// TODO refact block的产生和交易模块
	return db.GetDBInst().Set(blockchain.CurrentBlockKey(), block.Hash())
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

func (blockchain *BlockChain) CurrentBlock() (*Block, error) {
	var err error = nil
	var block *Block
	if currentBlock == nil {
		blockValue, err := db.GetDBInst().Get(blockchain.CurrentBlockKey())
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(blockValue, &block)
	}
	currentBlock := &Block{
		Height:       block.Height + 1,
		Nonce:        0,
		Fee:          blockchain.Fee,
		TotalFee:     0,
		PreviousHash: block.Hash(),
		Locker:       sync.RWMutex{},
		StatTree:     MPTPlus.MTP_Tree(db.GetDBInst(), block.StatTree.Root),
		TxTree:       MPTPlus.NewMTP(db.GetDBInst()),
		EventTree:    MPTPlus.NewMTP(db.GetDBInst()),
		Round:        consensus.NextRound(block.Round, block.Hash()),
	}
	return currentBlock, err
}

func (blockchain *BlockChain) CurrentBlockKey() []byte {
	buffer := bytes.Buffer{}
	buffer.WriteString(CurrentBlockKey)
	buffer.Write(blockchain.ChainId)
	return buffer.Bytes()
}

func (blockchain *BlockChain) WaitAndPack() {
	time.Sleep(BackboneBlockInterval * time.Second)
	blockchain.Pack()
}

func (blockchain *BlockChain) Pack() {
	block, _ := blockchain.CurrentBlock()
	block.Locker.Lock()
	defer block.Locker.Unlock()
	start := time.Now().Nanosecond()
	for ; !bytes.HasPrefix(block.CaculateHash(), []byte("FFFFFF")); block.NewNonce() {
	}
	end := time.Now().Nanosecond()
	fmt.Printf(`\ndifficulty="FFFFFF", cost=%d\n`, (end-start)/1e6)
	blockchain_manager.MainBlockChainConsensus.BlockBorn(block)
	//db.GetDBInst().Set(block.Hash(), block.Bytes())
	//blockchain.consensus.NewBlock(*block, ConsensusCb)
}

func ConsensusCb(blockChain BlockChain, block Block, result bool) {
	if result {
		blockChain.NewBlock(block)
	}
}
