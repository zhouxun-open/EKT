package blockchain_manager

import (
	"errors"
	"sync"

	"github.com/EducationEKT/EKT/io/ekt8/blockchain"
	"github.com/EducationEKT/EKT/io/ekt8/consensus"
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
)

var MainBlockChain *blockchain.BlockChain
var MainBlockChainConsensus consensus.DPOSConsensus

var blockchainManager *BlockchainManager

type BlockchainManager struct {
	Blockchains map[string]blockchain.BlockChain
	Consensuses map[string]consensus.Consensus
}

func init() {
	blockchainManager = &BlockchainManager{
		Blockchains: make(map[string]blockchain.BlockChain),
		Consensuses: make(map[string]consensus.Consensus),
	}
	MainBlockChain = &blockchain.BlockChain{blockchain.BackboneChainId, blockchain.InitStatus, sync.RWMutex{},
		blockchain.BackboneConsensus, 1e6, []byte("FFFFFF")}
	MainBlockChainConsensus = consensus.DPOSConsensus{}
	MainBlockChainConsensus.ManageBlockChain(MainBlockChain)
	go MainBlockChainConsensus.Run()
}

func GetManagerInst() *BlockchainManager {
	return blockchainManager
}

func GetMainChain() *blockchain.BlockChain {
	return MainBlockChain
}

func GetMainChainConsensus() consensus.DPOSConsensus {
	return MainBlockChainConsensus
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
