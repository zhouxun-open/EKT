package blockchain_manager

import (
	"sync"

	"github.com/EducationEKT/EKT/io/ekt8/blockchain"
	"github.com/EducationEKT/EKT/io/ekt8/consensus"
	"github.com/EducationEKT/EKT/io/ekt8/i_consensus"
	"github.com/EducationEKT/EKT/io/ekt8/tx_pool"
)

var MainBlockChain *blockchain.BlockChain
var MainBlockChainConsensus consensus.DPOSConsensus

var blockchainManager *BlockchainManager

type BlockchainManager struct {
	Blockchains map[string]blockchain.BlockChain
	Consensuses map[string]i_consensus.Consensus
}

func Init() {
	blockchainManager = &BlockchainManager{
		Blockchains: make(map[string]blockchain.BlockChain),
		Consensuses: make(map[string]i_consensus.Consensus),
	}
	MainBlockChain = &blockchain.BlockChain{blockchain.BackboneChainId, blockchain.InitStatus, nil, sync.RWMutex{},
		blockchain.BackboneConsensus, 1e6, []byte("FFFFFF"), tx_pool.NewTxPool(), 0}
	MainBlockChainConsensus = consensus.DPOSConsensus{Blockchain: MainBlockChain}
	//TODO 查看自己subscribe的blockchain
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
