package consensus

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/EducationEKT/EKT/io/ekt8/MPTPlus"
	"github.com/EducationEKT/EKT/io/ekt8/blockchain"
	"github.com/EducationEKT/EKT/io/ekt8/blockchain_manager"
	"github.com/EducationEKT/EKT/io/ekt8/crypto"
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
	block := dpos.CurrentBlock()
	if err := crypto.Validate(block.Bytes(), block.Hash()); err != nil {
		panic(err)
	}
	dpos.SyncBlock(block)
}

func (dpos DPOSConsensus) CurrentBlock() *blockchain.Block {
	var currentBlock *blockchain.Block = nil
	block, err := dpos.blockchain.CurrentBlock()
	if err == nil && block != nil {
		currentBlock = block
	}
	blocks := make(map[string]int64)
	mapping := make(map[string]*blockchain.Block)
	for _, peer := range dpos.Round.Peers {
		block, err := peer.CurrentBlock()
		if err != nil {
			continue
		}
		mapping[hex.EncodeToString(block.Hash())] = block
		num, exist := blocks[hex.EncodeToString(block.Hash())]
		if exist && num+1 >= int64(len(dpos.Round.Peers))/2 {
			currentBlock = block
			break
		} else {
			if exist {
				blocks[hex.EncodeToString(block.Hash())] = num + 1
			} else {
				blocks[hex.EncodeToString(block.Hash())] = 1
			}
		}
	}
	var maxNum int64 = 0
	var consensusHash string
	if currentBlock == nil {
		for hash, num := range blocks {
			if num > maxNum {
				maxNum, consensusHash = num, hash
			}
		}
	}
	return mapping[consensusHash]
}

func (dpos DPOSConsensus) SyncBlock(block *blockchain.Block) {
	MPTPlus.SyncDB(block.StatRoot, dpos.Round.Peers, false)
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
