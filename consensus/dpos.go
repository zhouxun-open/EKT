package consensus

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"xserver/x_http/x_resp"

	"sync"
	"time"

	"bytes"
	"errors"
	"github.com/EducationEKT/EKT/MPTPlus"
	"github.com/EducationEKT/EKT/blockchain"
	"github.com/EducationEKT/EKT/conf"
	"github.com/EducationEKT/EKT/core/userevent"
	"github.com/EducationEKT/EKT/crypto"
	"github.com/EducationEKT/EKT/ctxlog"
	"github.com/EducationEKT/EKT/db"
	"github.com/EducationEKT/EKT/log"
	"github.com/EducationEKT/EKT/p2p"
	"github.com/EducationEKT/EKT/param"
	"github.com/EducationEKT/EKT/pool"
	"github.com/EducationEKT/EKT/round"
	"github.com/EducationEKT/EKT/util"
	"os"
	"runtime"
	"strings"
)

type DPOSConsensus struct {
	Blockchain   *blockchain.BlockChain
	BlockManager *blockchain.BlockManager
	Block        chan blockchain.Block
	Vote         chan blockchain.BlockVote
	VoteResults  blockchain.VoteResults
	Locker       sync.RWMutex
}

func NewDPoSConsensus(Blockchain *blockchain.BlockChain) *DPOSConsensus {
	return &DPOSConsensus{
		Blockchain:   Blockchain,
		BlockManager: blockchain.NewBlockManager(),
		Block:        make(chan blockchain.Block),
		Vote:         make(chan blockchain.BlockVote),
		VoteResults:  blockchain.NewVoteResults(),
		Locker:       sync.RWMutex{},
	}
}

// 校验从其他委托人节点过来的区块数据
func (dpos DPOSConsensus) BlockFromPeer(ctxlog *ctxlog.ContextLog, block blockchain.Block) {
	dpos.Locker.Lock()
	defer dpos.Locker.Unlock()

	dpos.BlockManager.Insert(block)

	status := dpos.BlockManager.GetBlockStatus(block.CurrentHash)
	if status == blockchain.BLOCK_SAVED ||
		(status > blockchain.BLOCK_ERROR_START && status < blockchain.BLOCK_ERROR_END) ||
		status == blockchain.BLOCK_VOTED {
		//如果区块已经写入链中 or 是一个有问题的区块 or 已经投票成功 直接返回
		return
	}

	if status == blockchain.BLOCK_VALID {
		ctxlog.Log("SendVote", true)
		dpos.SendVote(block)
		dpos.BlockManager.SetBlockStatus(block.CurrentHash, blockchain.BLOCK_VOTED)
		return
	}

	// 判断此区块是否是一个interval之前打包的，如果是则放弃vote
	// unit： ms    单位：ms
	blockLatencyTime := int(time.Now().UnixNano()/1e6 - block.Timestamp) // 从节点打包到当前节点的延迟，单位ms
	blockInterval := int(dpos.Blockchain.BlockInterval / 1e6)            // 当前链的打包间隔，单位nanoSecond,计算为ms
	if blockLatencyTime > blockInterval {
		log.Info("Recieved a block packed before an interval, return.")
		ctxlog.Log("More than an interval", true)
		dpos.BlockManager.SetBlockStatus(block.CurrentHash, blockchain.BLOCK_ERROR_BROADCAST_TIME)
		return
	}

	// 校验打包节点在打包时是否有打包权限
	log.Info("Validating is the right node.")
	if result := dpos.PeerTurn(ctxlog, block.Timestamp, dpos.Blockchain.GetLastBlock().Timestamp, block.GetRound().Peers[block.GetRound().CurrentIndex]); !result {
		ctxlog.Log("rightNode", result)
		dpos.BlockManager.SetBlockStatus(block.CurrentHash, blockchain.BLOCK_ERROR_PACK_TIME)
		return
	} else {
		ctxlog.Log("rightNode", result)
	}

	// validate hash
	if !block.ValidateHash() {
		dpos.BlockManager.SetBlockStatus(block.CurrentHash, blockchain.BLOCK_ERROR_HASH)
		return
	}

	// validate sign
	if !dpos.ValidateSign(block) {
		dpos.BlockManager.SetBlockStatus(block.CurrentHash, blockchain.BLOCK_ERROR_SIGN)
		return
	}

	_block := block
	if !dpos.syncBlockBody(&_block) {
		return
	}

	events, err := dpos.getBlockEvents(&block)
	if err != nil {
		fmt.Println("Can't get events from local and peer, return false.", err)
		return
	}

	// 对区块进行validate和recover，如果区块数据没问题，则发送投票给其他节点
	if dpos.Blockchain.ValidateNextBlock(ctxlog, _block, events) {
		ctxlog.Log("SendVote", true)
		dpos.BlockManager.SetBlockStatus(block.CurrentHash, blockchain.BLOCK_VOTED)
		dpos.SendVote(block)
	} else {
		dpos.BlockManager.SetBlockStatus(block.CurrentHash, blockchain.BLOCK_ERROR_BODY)
	}
}

func (dpos DPOSConsensus) getBlockEvents(block *blockchain.Block) (list []userevent.IUserEvent, err error) {
	dpos.syncBlockBody(block)
	eventIds := block.BlockBody.Events
	if len(eventIds) > 0 {
		for _, eventId := range eventIds {
			id, err := hex.DecodeString(eventId)
			if err != nil {
				return nil, err
			}
			eventGetter := pool.NewEventGetter(eventId)
			dpos.Blockchain.Pool.EventGetter <- eventGetter
			event := <-eventGetter.Chan
			if event == nil {
				data, err := block.Round.Peers[block.Round.CurrentIndex].GetDBValue(id)
				if !bytes.EqualFold(crypto.Sha3_256(data), id) {
					return nil, errors.New("Invalid Response")
				}
				if err != nil {
					return nil, err
				}
				var tx userevent.Transaction
				err = json.Unmarshal(data, &tx)
				if err != nil {
					return nil, err
				}
				db.GetDBInst().Set(id, data)
				event = tx
				dpos.Blockchain.Pool.EventPutter <- event
			}
			if event != nil {
				list = append(list, event)
			} else {
				return nil, err
			}
		}
	}
	return list, nil
}

func (dpos DPOSConsensus) syncUserEvent(eventId string) bool {
	// TODO
	return true
}

func (dpos DPOSConsensus) syncBlockBody(block *blockchain.Block) bool {
	// 从打包节点获取body
	body, err := block.GetRound().Peers[block.GetRound().CurrentIndex].GetDBValue(block.Body)
	if err != nil {
		log.Info("Can not get body from mining node, return false.")
		return false
	}
	block.BlockBody, err = blockchain.FromBytes(body)
	if err != nil {
		log.Info("Get an error body, return false.")
		return false
	}
	return true
}

func (dpos DPOSConsensus) ValidateSign(block blockchain.Block) bool {
	if pubkey, err := crypto.RecoverPubKey(crypto.Sha3_256(block.CurrentHash), block.Signature); err != nil {
		return false
	} else {
		if !strings.EqualFold(hex.EncodeToString(crypto.Sha3_256(pubkey)), block.GetRound().Peers[block.GetRound().CurrentIndex].PeerId) {
			return false
		}
	}
	return true
}

// 校验从其他委托人节点来的区块成功之后发送投票
func (dpos DPOSConsensus) SendVote(block blockchain.Block) {
	// 同一个节点在一个出块interval内对一个高度只会投票一次，所以先校验是否进行过投票
	log.Info("Validating send vote interval.")
	// 获取上次投票时间 lastVoteTime < 0 表明当前区块没有投票过
	lastVoteTime := dpos.Blockchain.BlockManager.GetVoteTime(block.Height)
	if lastVoteTime > 0 {
		// 距离投票的毫秒数
		intervalInFact := int(time.Now().UnixNano()/1e6 - lastVoteTime)
		// 规则指定的毫秒数
		intervalInRule := int(dpos.Blockchain.BlockInterval / 1e6)

		// 说明在一个intervalInRule内进行过投票
		if intervalInFact < intervalInRule {
			log.Info("This height has voted in paste interval, return.")
			return
		}
	}

	// 记录此次投票的时间
	dpos.Blockchain.BlockManager.SetVoteTime(block.Height, time.Now().UnixNano()/1e6)

	// 生成vote对象
	vote := &blockchain.BlockVote{
		BlockchainId: dpos.Blockchain.ChainId,
		BlockHash:    block.Hash(),
		BlockHeight:  block.Height,
		VoteResult:   true,
		Peer:         conf.EKTConfig.Node,
	}

	// 签名
	err := vote.Sign(conf.EKTConfig.GetPrivateKey())
	if err != nil {
		log.Crit("Sign vote failed, recorded. %v", err)
		return
	}

	// 向其他节点发送签名后的vote信息
	log.Info("Signed this vote, sending vote result to other peers.")
	for i, peer := range block.GetRound().Peers {
		// 为了节省节点间带宽，只会向当前round内，距离打包节点近的n/2个节点
		if (i-block.GetRound().CurrentIndex+len(block.GetRound().Peers))%len(block.GetRound().Peers) <= len(block.GetRound().Peers)/2 {
			url := fmt.Sprintf(`http://%s:%d/vote/api/vote`, peer.Address, peer.Port)
			go util.HttpPost(url, vote.Bytes())
		}
	}
	log.Info("Send vote to other peer succeed.")
}

// for循环+recover保证DPoS线程的安全性
func (dpos *DPOSConsensus) StableRun() {
	for {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Crit(`Consensus occured an unknown error, recovered. %v`, r)
				}
			}()
			dpos.Run()
		}()
	}
}

// 委托人节点检测其他节点未按时出块的情况下， 当前节点进行打包的逻辑
func (dpos DPOSConsensus) DelegateRun() {
	log.Info("DPoS started.")

	//要求有半数以上节点存活才可以进行打包区块
	moreThanHalf := false
	for !moreThanHalf {
		if AliveDPoSPeerCount(param.MainChainDPosNode, false) <= len(param.MainChainDPosNode)/2 {
			log.Info("Alive node is less than half, waiting for other DPoS node restart.")
			time.Sleep(3 * time.Second)
		} else {
			moreThanHalf = true
		}
	}

	// 每1/4个interval检测一次是否有漏块，如果发生漏块且当前节点可以出块，则进入打包流程
	interval := dpos.Blockchain.BlockInterval / 4
	for {
		// 判断是否是当前节点打包区块
		log.Info(`Timer tick: is my turn?`)

		ctxlog := ctxlog.NewContextLog("DPoS is my turn ?")
		if dpos.IsMyTurn(ctxlog) {
			log.Info("This is my turn, current height is %d.", dpos.Blockchain.GetLastHeight())
			dpos.Pack(ctxlog)
		} else {
			log.Info("No, sleeping %d nano second.", int(interval))
		}
		ctxlog.Finish()

		time.Sleep(interval)
	}
}

// 判断peer在指定时间是否有打包区块的权力
func (dpos DPOSConsensus) PeerTurn(ctxlog *ctxlog.ContextLog, packTime, lastBlockTime int64, peer p2p.Peer) bool {
	log.Info("Validating peer has the right to pack block.")

	// 如果当前高度为0，则区块中不包含round，否则从block中取round
	_round := &round.Round{
		Peers:        param.MainChainDPosNode,
		CurrentIndex: -1,
	}
	if dpos.Blockchain.GetLastHeight() > 0 {
		_round = dpos.Blockchain.GetLastBlock().GetRound()
	}

	// 如果当前高度为0，则需要第一个节点进行打包
	if dpos.Blockchain.GetLastHeight() == 0 {
		if _round.Peers[0].Equal(peer) {
			return true
		} else {
			return false
		}
	}

	// 对一些变量进行记录
	ctxlog.Log("lastBlock", dpos.Blockchain.GetLastBlock())
	ctxlog.Log("CurrentNode", peer)
	ctxlog.Log("packTime", packTime)
	ctxlog.Log("lastBlockTime", lastBlockTime)

	intervalInFact, interval := int(packTime-lastBlockTime), int(dpos.Blockchain.BlockInterval/1e6)

	// n表示距离上次打包的间隔
	n := int(intervalInFact) / int(interval)
	remainder := int(intervalInFact) % int(interval)
	if remainder > int(interval)/2 {
		n++
	}

	// 如果距离上次打包在一个interval之内，返回false
	if n == 0 {
		return false
	}

	// 超过n个interval则需要第n+1个节点进行打包
	n++

	// 如果超过了当前round，则重新计算当前round
	nextRound := _round.Clone()
	if _round.CurrentIndex+n >= _round.Len() {
		nextRound = _round.Shuffle(round.GetRandomByHash(dpos.Blockchain.GetLastBlock().CurrentHash))
	}

	// 判断peer是否拥有打包权限
	nextRound.CurrentIndex = (_round.CurrentIndex + n) % _round.Len()
	if _round.Peers[_round.CurrentIndex].Equal(peer) {
		return true
	} else {
		return false
	}
}

// 用于委托人线程判断当前节点是否有打包权限
func (dpos DPOSConsensus) IsMyTurn(ctxlog *ctxlog.ContextLog) bool {
	now := time.Now().UnixNano() / 1e6
	lastPackTime := dpos.Blockchain.GetLastBlock().Timestamp
	result := dpos.PeerTurn(ctxlog, now, lastPackTime, conf.EKTConfig.Node)
	ctxlog.Log("result", result)

	return result
}

func (dpos *DPOSConsensus) Run() {
	// 从数据库中恢复当前节点已同步的区块
	log.Info("Recover data from local database.")
	dpos.RecoverFromDB()
	log.Info("Local data recovered. Current height is %d.\n", dpos.Blockchain.GetLastHeight())

	log.Info("Synchronizing blockchain ...")
	interval, failCount := 50*time.Millisecond, 0
	// 同步区块
	for height := dpos.Blockchain.GetLastHeight() + 1; ; {
		log.Info("Synchronizing block at height %d.", height)
		if dpos.SyncHeight(height) {
			log.Info("Synchronizing block at height %d successed. \n", height)
			height++
			// 同步成功之后，failCount变成0
			failCount = 0
		} else {
			log.Info("Synchronizing block at height %d failed. \n", height)
			failCount++
			// 如果区块同步失败，会重试三次，三次之后判断当前节点是否是DPoS节点，选择不同的同步策略
			if failCount >= 3 {
				log.Info("Fail count more than 3 times.")
				// 如果当前节点是DPoS节点，则不再根据区块高度同步区块，而是通过投票结果来同步区块
				for _, peer := range param.MainChainDPosNode {
					if peer.Equal(conf.EKTConfig.Node) {
						log.Info("This peer is DPoS node, start DPoS thread.")
						// 开启Delegate线程并让此线程sleep
						dpos.startDelegateThread()
						ch := make(chan bool)
						<-ch
					}
				}
				log.Info("Synchronize interval change to blockchain interval")
				interval = dpos.Blockchain.BlockInterval
			}
		}
		time.Sleep(interval)
	}
}

// 开启delegate线程
func (dpos *DPOSConsensus) startDelegateThread() {
	// 稳定启动dpos.DelegateRun()
	go func() {
		for {
			func() {
				defer func() {
					if r := recover(); r != nil {
						var buf [4096]byte
						runtime.Stack(buf[:], false)
						log.Crit("Panic occured at delegate thread, %s", string(buf[:]))
					}
				}()
				dpos.DelegateRun()
			}()
		}

	}()

	// 稳定启动dpos.dposSync()
	go func() {
		for {
			func() {
				defer func() {
					if r := recover(); r != nil {
						var buf [4096]byte
						runtime.Stack(buf[:], false)
						log.Crit("Panic occured at delegate sync thread, %s", string(buf[:]))
					}
				}()
				dpos.dposSync()
			}()
		}

	}()
}

// dposSync同步主要是监控在一定interval如果height没有被委托人间投票改变，则通过height进行同步
func (dpos *DPOSConsensus) dposSync() {
	lastHeight := dpos.Blockchain.GetLastHeight()
	for {
		height := dpos.Blockchain.GetLastHeight()
		log.Debug("Last interval lastHeight is %d, lastHeight is %d now.", lastHeight, height)
		if height == lastHeight {
			log.Debug("Height has not change for an interval, synchronizing block.")
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Crit("Panic occured, %v", r)
					}
				}()
				if dpos.SyncHeight(lastHeight + 1) {
					log.Debug("Synchronized block at lastHeight %d.", lastHeight+1)
					lastHeight = dpos.Blockchain.GetLastHeight()
				} else {
					log.Debug("Synchronize block at lastHeight %d failed.", lastHeight+1)
				}
			}()
		}

		lastHeight = dpos.Blockchain.GetLastHeight()

		time.Sleep(dpos.Blockchain.BlockInterval)
	}
}

// 共识向blockchain发送signal进行下一个区块的打包
func (dpos DPOSConsensus) Pack(ctxlog *ctxlog.ContextLog) {
	// 对下一个区块进行打包
	lastBlock := dpos.Blockchain.GetLastBlock()
	if !dpos.BlockManager.GetBlockStatusByHeight(lastBlock.Height+1, int64(dpos.Blockchain.BlockInterval)) {
		log.Debug("This height is packed within an interval, return nil.")
		return
	}
	block := dpos.Blockchain.PackSignal(ctxlog, lastBlock.Height+1)

	// 如果block不为空，说明打包成功，签名后转发给其他节点
	if block != nil {
		// dpos将当前round的信息更新到区块上
		if lastBlock.Round == nil {
			block.Round = &round.Round{
				Peers:        param.MainChainDPosNode,
				CurrentIndex: 0,
			}
		} else {
			duration := block.Timestamp - lastBlock.Timestamp
			intervalMs := int64(dpos.Blockchain.BlockInterval) / 1e6
			// 判断是否需要进入下一个round
			if int64(lastBlock.Round.MyIndex()-lastBlock.Round.CurrentIndex)*intervalMs < duration {
				block.Round = lastBlock.Round.Shuffle(round.GetRandomByHash(lastBlock.CurrentHash))
			} else {
				block.Round = lastBlock.Round
			}
			block.Round.CurrentIndex = block.Round.MyIndex()
		}

		// 计算hash
		block.CaculateHash()

		// 增加打包信息
		dpos.BlockManager.Insert(*block)
		dpos.BlockManager.SetBlockStatus(block.CurrentHash, blockchain.BLOCK_VALID)
		dpos.BlockManager.SetBlockStatusByHeight(block.Height, block.Timestamp)

		// 签名
		if err := block.Sign(conf.EKTConfig.PrivateKey); err != nil {
			log.Crit("Sign block failed. %v", err)
		} else {
			// 广播
			dpos.broadcastBlock(block)
			ctxlog.Log("block", block)
		}
	}
}

// 广播区块
func (dpos DPOSConsensus) broadcastBlock(block *blockchain.Block) {
	log.Info("Broadcasting block to the other peers.")
	data := block.Bytes()
	for _, peer := range block.GetRound().Peers {
		url := fmt.Sprintf(`http://%s:%d/block/api/newBlock`, peer.Address, peer.Port)
		go util.HttpPost(url, data)
	}
}

// 从db中recover数据
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
			BlockBody:    blockchain.NewBlockBody(),
			Body:         nil,
			Timestamp:    0,
			Locker:       sync.RWMutex{},
			StatTree:     MPTPlus.NewMTP(db.GetDBInst()),
			StatRoot:     nil,
			TxTree:       MPTPlus.NewMTP(db.GetDBInst()),
			TxRoot:       nil,
			TokenTree:    MPTPlus.NewMTP(db.GetDBInst()),
			TokenRoot:    nil,
		}

		for _, account := range accounts {
			block.CreateGenesisAccount(account)
		}

		block.UpdateMPTPlusRoot()

		block.CaculateHash()

		dpos.Blockchain.SaveBlock(*block)
	}
	dpos.Blockchain.SetLastBlock(*block)
}

// 获取存活的DPOS节点数量
func AliveDPoSPeerCount(peers p2p.Peers, print bool) int {
	count := 0
	for _, peer := range peers {
		if peer.IsAlive() {
			if print {
				log.Info("Peer %s is alive, address: %s \n", peer.PeerId, peer.Address)
			}
			count++
		}
	}
	return count
}

// 根据height同步区块
func (dpos DPOSConsensus) SyncHeight(height int64) bool {
	log.Info("Synchronizing block at height %d \n", height)
	if dpos.Blockchain.GetLastHeight() >= height {
		return true
	}
	round := &round.Round{
		Peers:        param.MainChainDPosNode,
		CurrentIndex: -1,
	}
	if dpos.Blockchain.GetLastHeight() > 0 {
		round = dpos.Blockchain.GetLastBlock().GetRound()
	}
	peers := param.MainChainDPosNode
	if dpos.Blockchain.GetLastHeight() > 0 {
		peers = round.Peers
	}
	for _, peer := range peers {
		block, err := getBlockHeader(peer, height)
		if err != nil || block.Height != height {
			log.Info("Geting block header by height failed. %v", err)
			continue
		}
		votes, err := getVotes(peer, hex.EncodeToString(block.CurrentHash))
		if err != nil {
			log.Info("Error peer has no votes. %v", err)
			continue
		}
		if votes.Validate() {
			events, err := dpos.getBlockEvents(block)
			if err != nil {
				continue
			}
			if dpos.Blockchain.GetLastBlock().ValidateNextBlock(*block, events) {
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

// 从其他委托人节点发过来的区块的投票进行记录
func (dpos DPOSConsensus) VoteFromPeer(vote blockchain.BlockVote) {
	log.Info("Recieved vote from peer.")
	if dpos.VoteResults.Broadcasted(vote.BlockHash) {
		log.Info("This block has voted, return.")
		return
	}
	dpos.VoteResults.Insert(vote)

	round := &round.Round{
		Peers:        param.MainChainDPosNode,
		CurrentIndex: -1,
	}
	if dpos.Blockchain.GetLastHeight() > 0 {
		round = dpos.Blockchain.GetLastBlock().GetRound()
	}

	log.Info("Is current vote number more than half node?")
	if dpos.VoteResults.Number(vote.BlockHash) > len(round.Peers)/2 {
		log.Info("Vote number more than half node, sending vote result to other nodes.")
		votes := dpos.VoteResults.GetVoteResults(hex.EncodeToString(vote.BlockHash))
		for _, peer := range round.Peers {
			url := fmt.Sprintf(`http://%s:%d/vote/api/voteResult`, peer.Address, peer.Port)
			go util.HttpPost(url, votes.Bytes())
		}
	} else {
		log.Info("Current vote results: %s", string(dpos.VoteResults.GetVoteResults(hex.EncodeToString(vote.BlockHash)).Bytes()))
		log.Info("Vote number is %d, less than %d, waiting for vote. \n", dpos.VoteResults.Number(vote.BlockHash), len(round.Peers)/2+1)
	}
}

// 收到从其他节点发送过来的voteResult，校验之后可以写入到区块链中
func (dpos DPOSConsensus) RecieveVoteResult(votes blockchain.Votes) bool {
	if !dpos.ValidateVotes(votes) {
		log.Info("Votes validate failed. %v", votes)
		return false
	}

	status := dpos.BlockManager.GetBlockStatus(votes[0].BlockHash)

	// 已经写入到链中
	if status == blockchain.BLOCK_SAVED {
		return true
	}

	// 未同步区块body
	if status > blockchain.BLOCK_ERROR_START && status < blockchain.BLOCK_ERROR_END {
		// 未同步区块体通过sync同步区块
		fmt.Printf("Invalid block and votes, block.hash = %s \n", hex.EncodeToString(votes[0].BlockHash))
		log.Crit("Invalid block and votes, block.hash = %s", hex.EncodeToString(votes[0].BlockHash))
		os.Exit(-1)
	}

	// 区块已经校验但未写入链中
	if status == blockchain.BLOCK_VALID || status == blockchain.BLOCK_VOTED {
		dpos.SaveVotes(votes)
		block, exist := dpos.BlockManager.GetBlock(votes[0].BlockHash)
		if !exist {
			return false
		}
		dpos.Blockchain.SaveBlock(block)
		body, err := block.Round.Peers[block.Round.CurrentIndex].GetDBValue(block.Body)
		if err != nil {

		}
		block.BlockBody, err = blockchain.FromBytes(body)
		dpos.Blockchain.NotifyPool(block)
		contextLog := ctxlog.NewContextLog("pack from vote result")
		if dpos.IsMyTurn(contextLog) {
			dpos.Pack(contextLog)
		}
		return true
	}

	return false
}

// 校验voteResults
func (dpos DPOSConsensus) ValidateVotes(votes blockchain.Votes) bool {
	if !votes.Validate() {
		return false
	}

	round := &round.Round{
		Peers:        param.MainChainDPosNode,
		CurrentIndex: -1,
	}
	if dpos.Blockchain.GetLastHeight() > 0 {
		round = dpos.Blockchain.GetLastBlock().GetRound()
	}

	if votes.Len() <= len(round.Peers)/2 {
		return false
	}
	return true
}

// 保存voteResults，用于同步区块时的校验
func (dpos DPOSConsensus) SaveVotes(votes blockchain.Votes) {
	dbKey := []byte(fmt.Sprintf("block_votes:%s", hex.EncodeToString(votes[0].BlockHash)))
	db.GetDBInst().Set(dbKey, votes.Bytes())
}

// 根据区块hash获取votes
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

// 根据height获取blockHeader
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

// 根据hash向委托人节点获取votes
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
