package consensus

import (
	"encoding/json"
	"fmt"

	"github.com/EducationEKT/EKT/io/ekt8/blockchain"
	"github.com/EducationEKT/EKT/io/ekt8/blockchain_manager"
	"github.com/EducationEKT/EKT/io/ekt8/p2p"
	"github.com/EducationEKT/EKT/io/ekt8/util"
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

func NextRound(round *Round, CurrentHash []byte) *Round {
	if round.CurrentIndex == len(round.Peers)-1 {
		bytes := CurrentHash[22:]
		Random := util.BytesToInt(bytes)
		round = &Round{
			CurrentIndex: -1,
			Peers:        round.Peers,
			Random:       Random,
		}
	} else {
		round.CurrentIndex++
	}
	return round
}

func (dpos DPOSConsensus) NewBlock(block blockchain.Block, cb ConsensusCallBack) {
	cb(*blockchain_manager.MainBlockChain, block, true)
}

func (dpos DPOSConsensus) ManageBlockChain(blockchain *blockchain.BlockChain) {
	dpos.blockchain = blockchain
}

func (dpos DPOSConsensus) Run() {
	if dpos.blockchain == nil {
		return
	}
	peers := dpos.GetCurrentDPOSPeers()
	dpos.Round = Round{CurrentIndex: -1, Peers: peers, Random: -1}
	height := dpos.CurrentHeight()
	dpos.SyncHeightStatTree(height)
	//TODO tasks
}

func (dpos DPOSConsensus) CurrentHeight() int64 {
	var currentHeight int64 = 0
	block, err := dpos.blockchain.CurrentBlock()
	if err == nil && block != nil {
		currentHeight = block.Height
	}
	heights := make(map[int64]int64)
	for _, peer := range dpos.Round.Peers {
		peerHeight, _ := peer.CurrentHeight()
		num, exist := heights[peerHeight]
		if exist && num+1 >= int64(len(dpos.Round.Peers))/2 {
			currentHeight = peerHeight
			break
		} else {
			if exist {
				heights[peerHeight] = num + 1
			} else {
				heights[peerHeight] = 1
			}
		}
	}
	var height, num int64 = 0, 0
	if currentHeight <= 0 {
		for _, height = range heights {
			if heights[height] > num {
				num = heights[height]
			}
		}
	}
	return height
}

func (dpos DPOSConsensus) SyncHeightStatTree(heigth int64) {
	//TODO
}

func (dpos DPOSConsensus) GetCurrentDPOSPeers() p2p.Peers {
	return p2p.MainChainDPosNode
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

func (round Round) String() string {
	peers, _ := json.Marshal(round.Peers)
	return fmt.Sprintf(`{"peers": %s, "random": %d}`, string(peers), round.Random)
}
