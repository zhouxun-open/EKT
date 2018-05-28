package blockchain

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"errors"

	"github.com/EducationEKT/EKT/io/ekt8/conf"
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"github.com/EducationEKT/EKT/io/ekt8/crypto"
	"github.com/EducationEKT/EKT/io/ekt8/db"
	"github.com/EducationEKT/EKT/io/ekt8/i_consensus"
	"github.com/EducationEKT/EKT/io/ekt8/log"
	"github.com/EducationEKT/EKT/io/ekt8/param"
	"github.com/EducationEKT/EKT/io/ekt8/pool"
	"github.com/EducationEKT/EKT/io/ekt8/util"
)

var BackboneChainId []byte
var BackboneChainDifficulty []byte
var EKTTokenId []byte

func init() {
	BackboneChainId, _ = hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000001")
	BackboneChainDifficulty = []byte("F")
	EKTTokenId, _ = hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000001")
}

const (
	BackboneConsensus     = i_consensus.DPOS
	BackboneBlockInterval = 3 * time.Second
	BackboneChainFee      = 210000
)

const (
	CurrentBlockKey = "CurrentBlockKey"
	InitStatus      = 0
	StartPackStatus = 100
)

type BlockChain struct {
	ChainId       []byte
	Consensus     i_consensus.ConsensusType
	CurrentBlock  *Block
	Locker        sync.RWMutex
	Status        int
	Fee           int64
	Difficulty    []byte
	Pool          *pool.Pool
	CurrentHeight int64
	Validator     *BlockValidator
	BlockInterval time.Duration
	Police        BlockPolice
	BlockManager  *BlockManager
	PackLock      sync.RWMutex
}

func NewBlockChain(chainId []byte, consensusType i_consensus.ConsensusType, fee int64, difficulty []byte, interval time.Duration) *BlockChain {
	return &BlockChain{
		ChainId:       chainId,
		Consensus:     consensusType,
		CurrentBlock:  nil,
		Locker:        sync.RWMutex{},
		Status:        InitStatus, // 100 正在计算MTProot, 150停止计算root,开始计算block Hash
		Fee:           fee,
		Difficulty:    difficulty,
		Pool:          pool.NewPool(),
		CurrentHeight: 0,
		Validator:     nil,
		BlockInterval: interval,
		Police:        NewBlockPolice(),
		BlockManager:  NewBlockManager(),
		PackLock:      sync.RWMutex{},
	}
}

func (blockchain *BlockChain) PackSignal(height int64) {
	blockchain.PackLock.Lock()
	defer blockchain.PackLock.Unlock()
	if blockchain.Status != StartPackStatus {
		blockchain.Status = StartPackStatus
		block := blockchain.WaitAndPack()
		hash := hex.EncodeToString(block.CurrentHash)
		blockchain.BlockManager.Lock()
		blockchain.BlockManager.Blocks[hash] = block
		blockchain.BlockManager.BlockStatus[hash] = BODY_SAVED
		blockchain.BlockManager.HeightManager[block.Height] = block.Timestamp
		blockchain.BlockManager.Unlock()
		if err := block.Sign(); err != nil {
			fmt.Println("Sign block failed.", err)
		} else {
			if err := blockchain.broadcastBlock(block); err != nil {
				fmt.Println("Broadcast block failed, reason: ", err)
			}
		}
		blockchain.Status = InitStatus
	}
}

func (blockchain *BlockChain) PackHeightValidate(height int64) bool {
	if blockchain.CurrentHeight+1 != height {
		return false
	}
	blockchain.BlockManager.RLock()
	defer blockchain.BlockManager.RUnlock()
	if !blockchain.BlockManager.GetBlockStatusByHeight(height, int64(blockchain.BlockInterval)) {
		return false
	}
	return true
}

func (blockchain *BlockChain) GetBlockByHeight(height int64) (*Block, error) {
	if height > blockchain.CurrentHeight {
		return nil, errors.New("Invalid height")
	}
	key := blockchain.GetBlockByHeightKey(height)
	data, err := db.GetDBInst().Get(key)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("Too heigher.")
	}
	block, err := FromBytes2Block(data)
	if block.Height != height {
		return nil, errors.New("Can not get block from db.")
	}
	return block, err
}

func (blockchain *BlockChain) GetBlockByHeightKey(height int64) []byte {
	return []byte(fmt.Sprint(`GetBlockByHeight: _%s_%d`, hex.EncodeToString(blockchain.ChainId), height))
}

func (blockchain *BlockChain) broadcastBlock(block *Block) error {
	fmt.Println("Broadcasting block to the other peers.")
	data := block.Bytes()
	for _, peer := range block.Round.Peers {
		url := fmt.Sprintf(`http://%s:%d/block/api/newBlock`, peer.Address, peer.Port)
		go util.HttpPost(url, data)
	}
	return nil
}

func (blockchain *BlockChain) SaveBlock(block *Block) {
	fmt.Println("Saving block to database.")
	db.GetDBInst().Set(block.Hash(), block.Data())
	data, _ := json.Marshal(block)
	db.GetDBInst().Set(blockchain.GetBlockByHeightKey(block.Height), data)
	db.GetDBInst().Set(blockchain.CurrentBlockKey(), data)
	blockchain.CurrentBlock = block
	blockchain.CurrentHeight = block.Height
	fmt.Println("Save block to database succeed.")
}

func (blockchain *BlockChain) LastBlock() (*Block, error) {
	var err error = nil
	var block *Block
	if currentBlock == nil {
		key := blockchain.CurrentBlockKey()
		data, err := db.GetDBInst().Get(key)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, &block)
		if err != nil {
			return nil, err
		}
		currentBlock = block
		return block, err
	}
	return currentBlock, err
}

func (blockchain *BlockChain) CurrentBlockKey() []byte {
	buffer := bytes.Buffer{}
	buffer.WriteString(CurrentBlockKey)
	buffer.Write(blockchain.ChainId)
	return buffer.Bytes()
}

func (blockchain *BlockChain) WaitAndPack() *Block {
	// 打包10500个交易大概需要0.95秒
	eventTimeout := time.After(blockchain.BlockInterval / 3)
	if blockchain.BlockInterval > 3*time.Second {
		eventTimeout = time.After(blockchain.BlockInterval - 2*time.Second)
	}
	round := &i_consensus.Round{
		Peers:        param.MainChainDPosNode,
		CurrentIndex: 0,
	}
	if blockchain.CurrentBlock.Height != 0 {
		round = blockchain.CurrentBlock.Round.MyRound(blockchain.CurrentBlock.CurrentHash)
	}
	block := NewBlock(blockchain.CurrentBlock, round)
	fmt.Println("Packing transaction and other events.")
	for {
		flag := false
		select {
		case <-eventTimeout:
			flag = true
			break
		default:
			// 因为要进行以太坊ERC20的映射和冷钱包，因此一期不支持地址的申请和加密算法的替换，只能打包转账交易 和 token发行
			tx := blockchain.Pool.FetchTx()
			if tx != nil {
				txResult := block.NewTransaction(tx, block.Fee)
				blockchain.Pool.Notify(tx.TransactionId())
				block.BlockBody.AddTxResult(*txResult)
			}
		}
		if flag {
			break
		}
	}
	blockchain.Pack(block)
	return block
}

// 当区块写入区块时，notify交易池，一些nonce比较大的交易可以进行打包
func (blockchain *BlockChain) NotifyPool(block *Block) {
	if block.BlockBody == nil {
		return
	}
	// Notify transaction
	if block.BlockBody.TxResults != nil && len(block.BlockBody.TxResults) > 0 {
		for _, txResult := range block.BlockBody.TxResults {
			blockchain.Pool.Notify(txResult.TxId)
		}
	}
	// Notify event
	if block.BlockBody.EventResults != nil && len(block.BlockBody.EventResults) > 0 {
		for _, eventResult := range block.BlockBody.EventResults {
			blockchain.Pool.NotifyEvent(eventResult.EventId)
		}
	}
}

// consensus 模块调用这个函数，获得一个block对象之后发送给其他节点，其他节点同意之后调用上面的NewBlock方法
func (blockchain *BlockChain) Pack(block *Block) {
	block.Locker.Lock()
	defer block.Locker.Unlock()
	bodyData, _ := json.Marshal(block.BlockBody)
	block.Body = crypto.Sha3_256(bodyData)
	db.GetDBInst().Set(block.Body, bodyData)
	start := time.Now().Nanosecond()
	fmt.Println("Caculating block hash.")
	block.UpdateMPTPlusRoot()
	for ; !bytes.HasPrefix(block.CaculateHash(), blockchain.Difficulty); block.NewNonce() {
	}
	end := time.Now().Nanosecond()
	fmt.Printf("Caculated block hash, cost %d ms. \n", (end-start+1e9)%1e9/1e6)
}

func (blockchain *BlockChain) BlockFromPeer(block Block) {
	fmt.Printf("Validating block from peer, block info: %s, block.Hash=%s \n", string(block.Bytes()), hex.EncodeToString(block.Hash()))
	if err := block.Validate(); err != nil {
		fmt.Printf("Block validate failed, %s. \n", err.Error())
		return
	}
	status := blockchain.Police.BlockFromPeer(block)
	//收到了当前节点的其他区块
	if status == -1 {
		evilBlock := blockchain.Police.GetEvilBlock(block)
		for _, peer := range block.Round.Peers {
			defer func() {
				if r := recover(); r != nil {
					log.GetLogInst().LogCrit("Sending evil block fail, recovered.", r)
				}
			}()
			url := fmt.Sprintf(`http://%s:%d/block/api/evilBlock`, peer.Address, peer.Port)
			util.HttpPost(url, evilBlock.Bytes())
		}
	}
	// 1500是毫秒和纳秒的单位乘以2/3计算得来的
	if time.Now().UnixNano()/1e6-block.Timestamp > int64(blockchain.BlockInterval/1500) {
		fmt.Printf("time.Now=%d, block.Time=%d, block.Interval=%d \n", time.Now().UnixNano()/1e6, block.Timestamp, int64(blockchain.BlockInterval/1500))
		fmt.Println("Block timestamp is more than 2/3 block interval, abort vote.")
		return
	}
	if !blockchain.CurrentBlock.ValidateNextBlock(block, blockchain.BlockInterval) {
		fmt.Println("This block from peer can not recover by last block, abort.")
		return
	}
	BlockRecorder.Blocks[hex.EncodeToString(block.Hash())] = &block
	// 签名
	vote := &BlockVote{
		BlockchainId: blockchain.ChainId,
		BlockHash:    block.Hash(),
		BlockHeight:  block.Height,
		VoteResult:   true,
		Peer:         conf.EKTConfig.Node,
	}
	err := vote.Sign(conf.EKTConfig.PrivateKey)
	if err != nil {
		log.GetLogInst().LogCrit("Sign vote failed, recorded. %v", err)
		fmt.Println("Sign vote failed, recorded.")
		return
	}
	fmt.Println("Sending vote result to other peers.")
	for i, peer := range block.Round.Peers {
		if (i-block.Round.CurrentIndex+len(block.Round.Peers))%len(block.Round.Peers) <= len(block.Round.Peers)/2 {
			url := fmt.Sprintf(`http://%s:%d/vote/api/vote`, peer.Address, peer.Port)
			util.HttpPost(url, vote.Bytes())
		}
	}
}

func (blockchain BlockChain) NewTransaction(tx *common.Transaction) bool {
	from, _ := hex.DecodeString(tx.From)
	if account, err := blockchain.CurrentBlock.GetAccount(from); err == nil && account != nil {
		if account.Nonce+1 == tx.Nonce {
			blockchain.Pool.ParkTx(tx, pool.Ready)
			return true
		} else if account.Nonce+1 < tx.Nonce {
			blockchain.Pool.ParkTx(tx, pool.Block)
			return true
		}
	}
	return false
}

func (blockchain BlockChain) VoteFromPeer(vote BlockVote) {
	fmt.Println("Recieved vote from peer.")
	if VoteResultManager.Broadcasted(vote.BlockHash) {
		fmt.Println("This block has voted, return.")
		return
	}
	VoteResultManager.Insert(vote)
	round := &i_consensus.Round{
		Peers:        param.MainChainDPosNode,
		CurrentIndex: -1,
	}
	if blockchain.CurrentHeight > 0 {
		round = blockchain.CurrentBlock.Round
	}
	fmt.Println("Is current vote number more than half node?")
	if VoteResultManager.Number(vote.BlockHash) > len(round.Peers)/2 {
		fmt.Println("Vote number more than half node, sending vote result to other nodes.")
		VoteResultManager.Locker.RLock()
		defer VoteResultManager.Locker.RUnlock()
		votes := VoteResultManager.VoteResults[hex.EncodeToString(vote.BlockHash)]
		for _, peer := range round.Peers {
			url := fmt.Sprintf(`http://%s:%d/vote/api/voteResult`, peer.Address, peer.Port)
			util.HttpPost(url, votes.Bytes())
		}
	} else {
		fmt.Printf("Current vote results: %s", string(VoteResultManager.VoteResults[hex.EncodeToString(vote.BlockHash)].Bytes()))
		fmt.Printf("Vote number is %d, less than %d, waiting for vote. \n", VoteResultManager.Number(vote.BlockHash), len(round.Peers)/2+1)
	}
}
