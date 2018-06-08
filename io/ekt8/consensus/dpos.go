package consensus

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"

	"xserver/x_http/x_resp"

	"sync"
	"time"

	"github.com/EducationEKT/EKT/io/ekt8/MPTPlus"
	"github.com/EducationEKT/EKT/io/ekt8/blockchain"
	"github.com/EducationEKT/EKT/io/ekt8/conf"
	"github.com/EducationEKT/EKT/io/ekt8/db"
	"github.com/EducationEKT/EKT/io/ekt8/i_consensus"
	"github.com/EducationEKT/EKT/io/ekt8/log"
	"github.com/EducationEKT/EKT/io/ekt8/p2p"
	"github.com/EducationEKT/EKT/io/ekt8/param"
	"github.com/EducationEKT/EKT/io/ekt8/util"
)

type DPOSConsensus struct {
	Blockchain  *blockchain.BlockChain
	Block       chan blockchain.Block
	Vote        chan blockchain.BlockVote
	VoteResults blockchain.VoteResults
	Locker      sync.RWMutex
	DPoSStatus  int // 0 未开始   100 正在进行中
}

func NewDPoSConsensus(Blockchain *blockchain.BlockChain) *DPOSConsensus {
	return &DPOSConsensus{
		Blockchain:  Blockchain,
		Block:       make(chan blockchain.Block),
		Vote:        make(chan blockchain.BlockVote),
		VoteResults: blockchain.NewVoteResults(),
		Locker:      sync.RWMutex{},
		DPoSStatus:  0,
	}
}

func (dpos DPOSConsensus) Start() {
	for {
		select {
		case block := <-dpos.Block:
			dpos.BlockFromPeer(block)
			//case
		}
	}
}

func (dpos DPOSConsensus) BlockFromPeer(block blockchain.Block) {
	dpos.Locker.Lock()
	defer dpos.Locker.Unlock()
	if int(time.Now().UnixNano()/1e6-block.Timestamp) > int(dpos.Blockchain.BlockInterval/1e6) {
		fmt.Println(time.Now().UnixNano()/1e6, block.Timestamp, dpos.Blockchain.BlockInterval/1e6)
		fmt.Println("Recieved a block packed before 1 second, return.")
	}
	if !dpos.PeerTurn(block.Timestamp, block.Round.Peers[block.Round.CurrentIndex]) {
		fmt.Println("This is not the right node, return false.")
	}
	dpos.Blockchain.BlockFromPeer(block)
}

func (dpos *DPOSConsensus) Run() {
	for {
		defer func() {
			if r := recover(); r != nil {
				log.GetLogInst().LogDebug(`Consensus occured an unknown error, recovered. %v`, r)
				log.GetLogInst().LogCrit(`Consensus occured an unknown error, recovered. %v`, r)
				fmt.Println(r)
			}
		}()
		dpos.RUN()
	}
}

func (dpos DPOSConsensus) DPoSRun() {
	fmt.Println("DPoS started.")
	interval := dpos.Blockchain.BlockInterval / 10
	for {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("A panic occurred.")
				log.GetLogInst().LogDebug("A panic occurred, %v.\n", r)
			}
		}()
		round := &i_consensus.Round{Peers: param.MainChainDPosNode, CurrentIndex: -1}
		if dpos.Blockchain.CurrentHeight > 0 {
			round = dpos.Blockchain.CurrentBlock.Round
		}
		if AliveDPoSPeerCount(round.Peers, false) <= len(round.Peers)/2 {
			fmt.Println("Alive node is less than half, waiting for other DPoS node restart.")
			time.Sleep(3 * time.Second)
			continue
		}
		log.GetLogInst().LogInfo(`Timer tick: is my turn?`)
		if dpos.IsMyTurn() {
			fmt.Printf("This is my turn, current heigth is %d. \n", dpos.Blockchain.CurrentHeight)
			log.GetLogInst().LogInfo("Yes.")
			dpos.Pack(dpos.Blockchain.CurrentHeight)
		} else {
			log.GetLogInst().LogInfo("No, sleeping %d nano second.", interval)
		}
		time.Sleep(interval)
	}
}

func (dpos DPOSConsensus) PeerTurn(packTime int64, peer p2p.Peer) bool {
	fmt.Println("Validating peer has the right to pack block.")
	round := &i_consensus.Round{
		Peers:        param.MainChainDPosNode,
		CurrentIndex: -1,
	}
	dpos.Blockchain.Locker.RLock()
	defer dpos.Blockchain.Locker.RUnlock()
	if dpos.Blockchain.CurrentHeight > 0 {
		round = dpos.Blockchain.CurrentBlock.Round
	} else {
		fmt.Println("Current height is 0, waiting for the first node pack block.")
		if round.Peers[0].Equal(peer) {
			fmt.Println("This is the first node, return true.")
			return true
		} else {
			fmt.Println("This is not the first node, return true.")
			return false
		}
	}
	time, interval := int(packTime-dpos.Blockchain.CurrentBlock.Timestamp), int(dpos.Blockchain.BlockInterval/1e6)
	if time >= interval*round.Len() {
		fmt.Println("More than a round time, waiting for the next node pack block.")
		if round.NextPeerRight(peer, dpos.Blockchain.CurrentBlock.CurrentHash) {
			fmt.Println("This is the next node, return true.")
			return true
		} else {
			fmt.Println("This is not the next node, return false.")
			return false
		}
	} else {
		n := int(time) / int(interval)
		remainder := int(time) % int(interval)
		if remainder > int(interval)/2 {
			n++
		}
		fmt.Printf("Current round is %s \n", round.String())
		if round.CurrentIndex+n >= round.Len() {
			round = round.NewRandom(dpos.Blockchain.CurrentBlock.CurrentHash)
			sort.Sort(round)
		}
		round.CurrentIndex = (round.CurrentIndex + n) % round.Len()
		fmt.Printf("Next round is %s, is my turn? \n", round.String())
		if round.Peers[round.CurrentIndex].Equal(peer) {
			fmt.Println("This is the next node, return true.")
			return true
		} else {
			fmt.Println("This is not the next node, return false.")
			return false
		}
	}
	return false
}

func (dpos DPOSConsensus) IsMyTurn() bool {
	return dpos.PeerTurn(time.Now().UnixNano()/1e6, conf.EKTConfig.Node)
}

func (dpos *DPOSConsensus) RUN() {
	// 从数据库中恢复当前节点已同步的区块
	fmt.Println("Recover data from local database.")
	dpos.RecoverFromDB()
	fmt.Printf("Local data recovered. Current height is %d.\n", dpos.Blockchain.CurrentHeight)

	//获取21个节点的集合
	peers := dpos.GetCurrentDPOSPeers()
WaitingNodes:
	loop := true
	for loop {
		fmt.Println("detecting alive nodes......")
		aliveCount := AliveDPoSPeerCount(peers, true)
		if aliveCount > len(peers)/2 {
			fmt.Printf("Alive node count is %d, starting synchronized block. \n", aliveCount)
			loop = false
		} else {
			if aliveCount == 0 {
				fmt.Println("There is no node alive.")
			}
			fmt.Println("The number of surviving nodes is less than half, waiting for other nodes to restart.")
			time.Sleep(3 * time.Second)
		}
	}
	fmt.Println("Alive node more than half, continue.")

	fmt.Println("Synchronizing blockchain...")
	interval, failCount := 50*time.Millisecond, 0
	for height := dpos.Blockchain.CurrentHeight + 1; ; {
		if dpos.SyncHeight(height) {
			fmt.Printf("Synchronizing block at height %d successed. \n", height)
			height++
			failCount = 0
		} else {
			fmt.Printf("Synchronizing block at height %d failed. \n", height)
			round := &i_consensus.Round{
				Peers:        param.MainChainDPosNode,
				CurrentIndex: -1,
			}
			if dpos.Blockchain.CurrentHeight > 0 {
				round = dpos.Blockchain.CurrentBlock.Round
			}
			if AliveDPoSPeerCount(peers, false) <= len(round.Peers)/2 {
				goto WaitingNodes
			}
			failCount++
			// 如果区块同步失败，会重试三次，三次之后判断当前节点是否是DPoS节点，选择不同的同步策略
			if failCount >= 3 {
				fmt.Println("Fail count more than 3 times.")
				// 如果当前节点是DPoS节点，则不再根据区块高度同步区块，而是通过投票结果来同步区块
				if round.MyIndex() != -1 {
					fmt.Println("This peer is DPoS node, start DPoS thread.")
					dpos.startDPOS()
				}
				interval = 3 * time.Second
			}
		}
		time.Sleep(interval)
	}
}

func (dpos *DPOSConsensus) startDPOS() {
	dpos.Locker.Lock()
	if dpos.DPoSStatus == 100 {
		dpos.Locker.Unlock()
		fmt.Println("Dpos goroutine is already running, return.")
		return
	} else {
		fmt.Printf("Status is %d, starting dpos goroutine.", dpos.DPoSStatus)
		dpos.DPoSStatus = 100
		go dpos.DPoSRun()
		dpos.Locker.Unlock()
	}
}

// 共识向blockchain发送signal进行下一个区块的打包
func (dpos DPOSConsensus) Pack(height int64) {
	bc := dpos.Blockchain
	bc.PackSignal(height)
}

func (dpos DPOSConsensus) BlockMinedCallBack(block *blockchain.Block) {
	fmt.Println("Mined block, sending block to other dpos  peer.")
	fmt.Println(dpos.Blockchain.CurrentBlock.Round)
	for _, peer := range block.Round.Peers {
		url := fmt.Sprintf("http://%s:%d/block/api/newBlock", peer.Address, peer.Port)
		resp, err := util.HttpPost(url, block.Bytes())
		fmt.Println(string(resp), err)
	}
}

func (dpos DPOSConsensus) RecoverFromDB() {
	block, err := dpos.Blockchain.LastBlock()
	// 如果是第一次打开
	if err != nil || block == nil {
		// 将创世块写入数据库
		accounts := conf.EKTConfig.GenesisBlockAccounts
		block = &blockchain.Block{
			Height:       0,
			Nonce:        0,
			Fee:          dpos.Blockchain.Fee,
			TotalFee:     0,
			PreviousHash: nil,
			CurrentHash:  nil,
			BlockBody:    blockchain.NewBlockBody(0),
			Body:         nil,
			Round:        nil,
			Timestamp:    0,
			Locker:       sync.RWMutex{},
			StatTree:     MPTPlus.NewMTP(db.GetDBInst()),
			StatRoot:     nil,
			TxTree:       MPTPlus.NewMTP(db.GetDBInst()),
			TxRoot:       nil,
			EventTree:    MPTPlus.NewMTP(db.GetDBInst()),
			EventRoot:    nil,
			TokenTree:    MPTPlus.NewMTP(db.GetDBInst()),
			TokenRoot:    nil,
		}
		for _, account := range accounts {
			block.InsertAccount(account)
		}
		block.UpdateMPTPlusRoot()
		block.CaculateHash()
		dpos.Blockchain.SaveBlock(block)
	}
	dpos.Blockchain.CurrentHeight = block.Height
	dpos.Blockchain.CurrentBlock = block
}

//获取存活的DPOS节点数量
func AliveDPoSPeerCount(peers p2p.Peers, print bool) int {
	count := 0
	for _, peer := range peers {
		if peer.IsAlive() {
			if print {
				fmt.Printf("Peer %s is alive, address: %s \n", peer.PeerId, peer.Address)
			}
			count++
		}
	}
	return count
}

func (dpos DPOSConsensus) SyncHeight(height int64) bool {
	fmt.Printf("Synchronizing block at height %d \n", height)
	round := &i_consensus.Round{
		Peers:        param.MainChainDPosNode,
		CurrentIndex: -1,
	}
	if dpos.Blockchain.CurrentHeight > 0 {
		round = dpos.Blockchain.CurrentBlock.Round
	}
	peers := param.MainChainDPosNode
	if dpos.Blockchain.CurrentHeight > 0 {
		peers = round.Peers
	}
	for _, peer := range peers {
		block, err := getBlockHeader(peer, height)
		if err != nil || block.Height != height {
			fmt.Println("Geting block header by height failed.", err)
			continue
		}
		votes, err := getVotes(peer, hex.EncodeToString(block.CurrentHash))
		if err != nil {
			fmt.Println("Error peer has no votes.", err)
			continue
		}
		if votes.Validate() {
			if dpos.Blockchain.CurrentBlock.ValidateNextBlock(*block, dpos.Blockchain.BlockInterval) {
				if dpos.RecieveVoteResult(votes) {
					return true
				} else {
					continue
				}
			}
		}
	}
	return false
}

func (dpos DPOSConsensus) VoteFromPeer(vote blockchain.BlockVote) {
	fmt.Println("Recieved vote from peer.")
	if dpos.VoteResults.Broadcasted(vote.BlockHash) {
		fmt.Println("This block has voted, return.")
		return
	}
	dpos.VoteResults.Insert(vote)
	round := &i_consensus.Round{
		Peers:        param.MainChainDPosNode,
		CurrentIndex: -1,
	}
	if dpos.Blockchain.CurrentHeight > 0 {
		round = dpos.Blockchain.CurrentBlock.Round
	}
	fmt.Println("Is current vote number more than half node?")
	if dpos.VoteResults.Number(vote.BlockHash) > len(round.Peers)/2 {
		fmt.Println("Vote number more than half node, sending vote result to other nodes.")
		votes := dpos.VoteResults.GetVoteResults(hex.EncodeToString(vote.BlockHash))
		for _, peer := range round.Peers {
			url := fmt.Sprintf(`http://%s:%d/vote/api/voteResult`, peer.Address, peer.Port)
			resp, err := util.HttpPost(url, votes.Bytes())
			log.GetLogInst().LogDebug(`Resp: %s, err: %v`, string(resp), err)
		}
	} else {
		fmt.Printf("Current vote results: %s", string(dpos.VoteResults.GetVoteResults(hex.EncodeToString(vote.BlockHash)).Bytes()))
		fmt.Printf("Vote number is %d, less than %d, waiting for vote. \n", dpos.VoteResults.Number(vote.BlockHash), len(round.Peers)/2+1)
	}
}

func (dpos DPOSConsensus) RecieveVoteResult(votes blockchain.Votes) bool {
	if !dpos.ValidateVotes(votes) {
		fmt.Println("Votes validate failed. ", votes)
		return false
	}
	status := blockchain.BlockRecorder.GetStatus(hex.EncodeToString(votes[0].BlockHash))
	// 未同步区块body
	if status == -1 {
		// 未同步区块体通过sync同步区块
		return false
	}
	if block := blockchain.BlockRecorder.GetBlock(hex.EncodeToString(votes[0].BlockHash)); block != nil {
		if status == 100 {
			// 已同步区块body，但是未写入区块链中
			fmt.Println("Recieve vote result and get this block, saving block.")
			dpos.SaveVotes(votes)
			dpos.Blockchain.NotifyPool(block)
			dpos.Blockchain.SaveBlock(block)
			blockchain.BlockRecorder.SetStatus(hex.EncodeToString(block.CurrentHash), 200)
		} else if status == 200 {
			// 已经写入区块链中
			fmt.Println("This block is already wrote to blockchain.")
		}
		return true
	} else {
		fmt.Println("Haven't recieve this block,  abort.")
	}
	return false
}

func (dpos DPOSConsensus) ValidateVotes(votes blockchain.Votes) bool {
	if !votes.Validate() {
		return false
	}
	round := &i_consensus.Round{
		Peers:        param.MainChainDPosNode,
		CurrentIndex: -1,
	}
	if dpos.Blockchain.CurrentHeight > 0 {
		round = dpos.Blockchain.CurrentBlock.Round
	}
	if votes.Len() <= len(round.Peers)/2 {
		return false
	}
	return true
}

func (dpos DPOSConsensus) SaveVotes(votes blockchain.Votes) {
	dbKey := []byte(fmt.Sprintf("block_votes:%s", hex.EncodeToString(votes[0].BlockHash)))
	db.GetDBInst().Set(dbKey, votes.Bytes())
}

func (dpos DPOSConsensus) GetVotes(blockHash string) blockchain.Votes {
	dbKey := []byte(fmt.Sprintf("block_votes:%s", blockHash))
	data, err := db.GetDBInst().Get(dbKey)
	if err != nil {
		return nil
	}
	var votes blockchain.Votes
	err = json.Unmarshal(data, &votes)
	if err != nil {
		return nil
	}
	return votes
}

//获取当前的peers
func (dpos DPOSConsensus) GetCurrentDPOSPeers() p2p.Peers {
	return param.MainChainDPosNode
}

func getBlockHeader(peer p2p.Peer, height int64) (*blockchain.Block, error) {
	url := fmt.Sprintf(`http://%s:%d/block/api/blockByHeight?height=%d`, peer.Address, peer.Port, height)
	body, err := util.HttpGet(url)
	if err != nil {
		return nil, err
	}
	var resp x_resp.XRespBody
	err = json.Unmarshal(body, &resp)
	data, err := json.Marshal(resp.Result)
	if err == nil {
		var block blockchain.Block
		err = json.Unmarshal(data, &block)
		return &block, err
	}
	return nil, err
}

func getVotes(peer p2p.Peer, blockHash string) (blockchain.Votes, error) {
	url := fmt.Sprintf(`http://%s:%d/vote/api/getVotes?hash=%s`, peer.Address, peer.Port, blockHash)
	body, err := util.HttpGet(url)
	if err != nil {
		return nil, err
	}
	var resp x_resp.XRespBody
	err = json.Unmarshal(body, &resp)
	if err == nil && resp.Status == 0 {
		var votes blockchain.Votes
		data, err := json.Marshal(resp.Result)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, &votes)
		return votes, err
	}
	return nil, err
}
