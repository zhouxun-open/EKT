package consensus

import (
	"github.com/EducationEKT/EKT/io/ekt8/blockchain"
	"github.com/EducationEKT/EKT/io/ekt8/blockchain_manager"
)

type DPOSConsensus struct {
}

func (this DPOSConsensus) NewBlock(block blockchain.Block, cb ConsensusCallBack) {
	cb(*blockchain_manager.MainBlockChain, block, true)
}
