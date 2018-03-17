package consensus

import (
	"github.com/EducationEKT/EKT/io/ekt8/blockchain"
)

const (
	DPOS = 1
	POW  = 2
	POS  = 3
)

type ConsensusType int

type ConsensusCallBack func(blockChain blockchain.BlockChain, block blockchain.Block, result bool)

type Consensus interface {
	ManageBlockChain(chain *blockchain.BlockChain)
	Run()
}
