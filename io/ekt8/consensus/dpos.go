package consensus

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"xserver/x_http/x_resp"

	"github.com/EducationEKT/EKT/io/ekt8/MPTPlus"
	"github.com/EducationEKT/EKT/io/ekt8/blockchain"
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"github.com/EducationEKT/EKT/io/ekt8/crypto"
	"github.com/EducationEKT/EKT/io/ekt8/i_consensus"
	"github.com/EducationEKT/EKT/io/ekt8/p2p"
	"github.com/EducationEKT/EKT/io/ekt8/util"
)

type DPOSConsensus struct {
	Round      i_consensus.Round
	Blockchain *blockchain.BlockChain
}

func (dpos DPOSConsensus) ValidateBlock(header *blockchain.Block) {
	peer := header.Round.Peers[header.Round.CurrentIndex]
	body, err := getBlockBody(peer, header.Height)
	if err != nil || body.Height != header.Height {
		// TODO vote false
	}
	// TODO validate body

	//TODO vote true
}

//从网络层转发过来的交易,进入打包流程
func (dpos DPOSConsensus) NewTransaction(tx common.Transaction) {
	dpos.Blockchain.Locker.Lock()
	defer dpos.Blockchain.Locker.Unlock()
	lastBlock, _ := dpos.Blockchain.LastBlock()
	if dpos.Blockchain.Status == blockchain.OpenStatus {
		var account common.Account
		address, _ := hex.DecodeString(tx.From)
		if err := lastBlock.StatTree.GetInterfaceValue(address, &account); err != nil {
			if account.GetNonce()+1 < tx.Nonce {
			}
		}
	}
}

//
func (dpos DPOSConsensus) BlockBorn(block *blockchain.Block) {
}

func (dpos DPOSConsensus) Run() {
	//获取21个节点的集合
	peers := dpos.GetCurrentDPOSPeers()
	//
	dpos.Round = i_consensus.Round{CurrentIndex: -1, Peers: peers, Random: -1}
	//获取当前的待验证block header
	block := dpos.CurrentBlock()
	if block == nil {
		block = &blockchain.Block{}
	}
	//验证block是否合法
	if err := crypto.Validate(block.Bytes(), block.CaculateHash()); err != nil {
		panic(err)
	}
	//异步在全局添加区块到区块链
	dpos.SyncBlock(block)
}

func (dpos DPOSConsensus) pullBlock() {
	for {
		peers := dpos.Blockchain.CurrentBlock.Round.Peers
		for _, peer := range peers {
			block, _ := CurrentBlock(peer)
			if dpos.Blockchain.CurrentBlock.Height < block.Height {
			}
		}
	}
}

//从其他节点得到待验证block header
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

//同步区块链
func (dpos DPOSConsensus) SyncBlockChain() {
	lastBlock, err := dpos.Blockchain.LastBlock()
	if err != nil {
		lastBlock = nil
	}
	peerLast := dpos.CurrentBlock()
	if peerLast.Height > lastBlock.Height {
		dpos.Blockchain.NewBlock(peerLast)
	}
}

//根据区块header同步body
func (dpos DPOSConsensus) SyncBlock(block *blockchain.Block) {
	MPTPlus.SyncDB(block.StatRoot, dpos.Round.Peers, false)
}

//获取当前的peers
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

//向指定节点获取最新区块
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

func getBlockBody(peer p2p.Peer, height int64) (*blockchain.BlockBody, error) {
	url := fmt.Sprintf(`http://%s:%d/block/api/body?height=%d`, peer.Address, peer.Port, height)
	body, err := util.HttpGet(url)
	if err != nil {
		return nil, err
	}
	var resp x_resp.XRespBody
	err = json.Unmarshal(body, &resp)
	data, err := json.Marshal(resp.Result)
	if err == nil {
		var blockBody blockchain.BlockBody
		err = json.Unmarshal(data, &blockBody)
		return &blockBody, err
	}
	return nil, err
}
