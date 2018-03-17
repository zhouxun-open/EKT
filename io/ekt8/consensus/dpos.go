package consensus

import (
	"github.com/EducationEKT/EKT/io/ekt8/blockchain"
	"github.com/EducationEKT/EKT/io/ekt8/blockchain_manager"
	"github.com/EducationEKT/EKT/io/ekt8/p2p"
)

type DPOSConsensus struct {
	Round      Round
	blockchain *blockchain.BlockChain
}

type Round struct {
	CurrentIndex int // default -1
	Peers        []p2p.Peer
	Random       int
}

func NextRound(round *Round, Random int) *Round {
	newRound := &Round{
		CurrentIndex: -1,
		Peers:        round.Peers,
		Random:       Random,
	}
	return newRound
}

func (dpos DPOSConsensus) NewBlock(block blockchain.Block, cb ConsensusCallBack) {
	cb(*blockchain_manager.MainBlockChain, block, true)
}

func (dpos DPOSConsensus) ManageBlockChain(blockchain *blockchain.BlockChain) {
	dpos.blockchain = blockchain
}

func (dpos DPOSConsensus) Run() {
	// TODO
}

func (round Round) Len() int {
	return len(round.Peers)
}

func (round Round) Swap(i, j int) {
	round.Peers[i], round.Peers[j] = round.Peers[j], round.Peers[i]
}

func (round Round) Less(i, j int) bool {
	return round.Random%(i+j)%2 == 1
}
