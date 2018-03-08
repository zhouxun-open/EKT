package blockchain_manager

import (
	"errors"
	"sync"

	"github.com/EducationEKT/EKT/io/ekt8/blockchain"
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
)

var MainBlockChain *blockchain.BlockChain

func init() {
	MainBlockChain = &blockchain.BlockChain{blockchain.BackboneChainId, blockchain.InitStatus, sync.RWMutex{}, blockchain.BackboneConsensus}
	MainBlockChain.SyncBlockChain()
}

type Engine struct {
	blockChain  *blockchain.BlockChain
	Pack        chan bool
	Transaction chan *common.Transaction
	Status      int
}

func (engine *Engine) NewTransaction(transaction *common.Transaction) error {
	if engine.Status == 100 {
		//block, err := engine.blockChain.CurrentBlock()
		//if err != nil {
		//	return err
		//}
		//Transaction(block, transaction)
	}
	return errors.New("Wait next block")
}
