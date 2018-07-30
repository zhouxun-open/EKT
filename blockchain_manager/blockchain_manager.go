package blockchain_manager

import (
	"encoding/json"

	"github.com/EducationEKT/EKT/blockchain"
	"github.com/EducationEKT/EKT/consensus"
	"github.com/EducationEKT/EKT/db"
	"github.com/EducationEKT/EKT/i_consensus"
)

const (
	BlockchainManagerDBKey = "BlockchainManagerDBKey"
)

var MainBlockChain *blockchain.BlockChain
var MainBlockChainConsensus *consensus.DPOSConsensus

var blockchainManager *BlockchainManager

type BlockchainManager struct {
	Blockchains map[int64]*blockchain.BlockChain
	Consensuses map[int64]i_consensus.Consensus
}

func Init() {
	blockchainManager = &BlockchainManager{
		Blockchains: make(map[int64]*blockchain.BlockChain),
		Consensuses: make(map[int64]i_consensus.Consensus),
	}
	MainBlockChain = blockchain.NewBlockChain(blockchain.BackboneChainId, blockchain.BackboneConsensus, blockchain.BackboneChainFee, nil, blockchain.BackboneBlockInterval)
	MainBlockChainConsensus = consensus.NewDPoSConsensus(MainBlockChain)
	go MainBlockChainConsensus.StableRun()
	value, err := db.GetDBInst().Get([]byte(BlockchainManagerDBKey))
	if err != nil {
		return
	}
	blockchains := make([]*blockchain.BlockChain, 0)
	err = json.Unmarshal(value, &blockchains)
	if err != nil {
		return
	}
	for _, bc := range blockchains {
		blockchainManager.Blockchains[bc.ChainId] = bc
		switch bc.Consensus {
		case i_consensus.DPOS:
			consensus := consensus.NewDPoSConsensus(bc)
			blockchainManager.Consensuses[bc.ChainId] = consensus
			go consensus.StableRun()
		default:
			consensus := consensus.NewDPoSConsensus(bc)
			blockchainManager.Consensuses[bc.ChainId] = consensus
			go consensus.StableRun()
		}
	}
}

func GetMainChain() *blockchain.BlockChain {
	return MainBlockChain
}

func GetMainChainConsensus() *consensus.DPOSConsensus {
	return MainBlockChainConsensus
}
