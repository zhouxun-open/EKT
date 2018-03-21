package consensus

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/EducationEKT/EKT/io/ekt8/MPTPlus"
	"github.com/EducationEKT/EKT/io/ekt8/blockchain"
	"github.com/EducationEKT/EKT/io/ekt8/crypto"
	"github.com/EducationEKT/EKT/io/ekt8/i_consensus"
	"github.com/EducationEKT/EKT/io/ekt8/p2p"
	"github.com/EducationEKT/EKT/io/ekt8/util"
)

type DPOSConsensus struct {
	Round      i_consensus.Round
	Blockchain *blockchain.BlockChain
}

func (dpos DPOSConsensus) BlockBorn(block *blockchain.Block) {
}

func (dpos DPOSConsensus) Run() {
	peers := dpos.GetCurrentDPOSPeers()
	dpos.Round = i_consensus.Round{CurrentIndex: -1, Peers: peers, Random: -1}
	block := dpos.CurrentBlock()
	if err := crypto.Validate(block.Bytes(), block.CaculateHash()); err != nil {
		panic(err)
	}
	dpos.SyncBlock(block)
}

func (dpos DPOSConsensus) CurrentBlock() *blockchain.Block {
	var currentBlock *blockchain.Block = nil
	blocks := make(map[string]int64)
	mapping := make(map[string]*blockchain.Block)
	for _, peer := range dpos.Round.Peers {
		block, err := CurrentBlock(peer)
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

func CurrentHeight(peer p2p.Peer) (int64, error) {
	url := fmt.Sprintf(`http://%s:%d/blocks/api/last`, peer.Address, peer.Port)
	body, err := util.HttpGet(url)
	if err != nil {
		return -1, err
	}
	var block blockchain.Block
	err = json.Unmarshal(body, &block)
	return block.Height, err
}

func CurrentBlock(peer p2p.Peer) (*blockchain.Block, error) {
	url := fmt.Sprintf(`http://%s:%d/blocks/api/last`, peer.Address, peer.Port)
	body, err := util.HttpGet(url)
	if err != nil {
		return nil, err
	}
	var block blockchain.Block
	err = json.Unmarshal(body, &block)
	return &block, err
}
