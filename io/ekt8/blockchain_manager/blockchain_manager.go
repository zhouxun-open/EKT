package blockchain_manager

import (
	"encoding/hex"
	"encoding/json"

	"github.com/EducationEKT/EKT/io/ekt8/blockchain"
	"github.com/EducationEKT/EKT/io/ekt8/consensus"
	"github.com/EducationEKT/EKT/io/ekt8/db"
	"github.com/EducationEKT/EKT/io/ekt8/i_consensus"
)

const (
	BlockchainManagerDBKey = "BlockchainManagerDBKey"
)

var MainBlockChain *blockchain.BlockChain
var MainBlockChainConsensus consensus.DPOSConsensus

var blockchainManager *BlockchainManager

type BlockchainManager struct {
	Blockchains map[string]*blockchain.BlockChain
	Consensuses map[string]i_consensus.Consensus
}

func Init() {
	blockchainManager = &BlockchainManager{
		Blockchains: make(map[string]*blockchain.BlockChain),
		Consensuses: make(map[string]i_consensus.Consensus),
	}
	MainBlockChain = blockchain.NewBlockChain(blockchain.BackboneChainId, blockchain.BackboneConsensus, blockchain.BackboneChainFee, blockchain.BackboneChainDifficulty, blockchain.BackboneBlockInterval)
	MainBlockChainConsensus = consensus.NewDPoSConsensus(MainBlockChain)
	go MainBlockChainConsensus.Run()
	value, err := db.GetDBInst().Get([]byte(BlockchainManagerDBKey))
	if err != nil {
		return
	}
	blockchains := make([]*blockchain.BlockChain, 0)
	err = json.Unmarshal(value, &blockchains)
	if err != nil {
		return
	}
	for _, blockchain := range blockchains {
		chainId := hex.EncodeToString(blockchain.ChainId)
		blockchainManager.Blockchains[chainId] = blockchain
		switch blockchain.Consensus {
		case i_consensus.DPOS:
			consensus := consensus.NewDPoSConsensus(blockchain)
			blockchainManager.Consensuses[chainId] = consensus
			go consensus.Run()
		default:
			consensus := consensus.NewDPoSConsensus(blockchain)
			blockchainManager.Consensuses[chainId] = consensus
			go consensus.Run()
		}
	}
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
