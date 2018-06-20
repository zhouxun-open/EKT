package blockchain

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"errors"

	"github.com/EducationEKT/EKT/io/ekt8/context_log"
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"github.com/EducationEKT/EKT/io/ekt8/crypto"
	"github.com/EducationEKT/EKT/io/ekt8/db"
	"github.com/EducationEKT/EKT/io/ekt8/i_consensus"
	"github.com/EducationEKT/EKT/io/ekt8/log"
	"github.com/EducationEKT/EKT/io/ekt8/param"
	"github.com/EducationEKT/EKT/io/ekt8/pool"
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
	currentLocker sync.RWMutex
	currentBlock  *Block
	currentHeight int64
	Locker        sync.RWMutex
	Status        int
	Fee           int64
	Difficulty    []byte
	Pool          *pool.Pool
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
		currentBlock:  nil,
		Locker:        sync.RWMutex{},
		currentLocker: sync.RWMutex{},
		Status:        InitStatus, // 100 正在计算MTProot, 150停止计算root,开始计算block Hash
		Fee:           fee,
		Difficulty:    difficulty,
		Pool:          pool.NewPool(),
		currentHeight: 0,
		Validator:     nil,
		BlockInterval: interval,
		Police:        NewBlockPolice(),
		BlockManager:  NewBlockManager(),
		PackLock:      sync.RWMutex{},
	}
}

func (blockchain *BlockChain) GetLastBlock() *Block {
	blockchain.currentLocker.RLock()
	defer blockchain.currentLocker.RUnlock()
	return blockchain.currentBlock
}

func (blockchain *BlockChain) SetLastBlock(block *Block) {
	blockchain.currentLocker.Lock()
	defer blockchain.currentLocker.Unlock()
	blockchain.currentBlock = block
}

func (blockchain *BlockChain) GetLastHeight() int64 {
	blockchain.currentLocker.RLock()
	defer blockchain.currentLocker.RUnlock()
	return blockchain.currentHeight
}

func (blockchain *BlockChain) SetLastHeight(height int64) {
	blockchain.currentLocker.RLock()
	defer blockchain.currentLocker.RUnlock()
	blockchain.currentHeight = height
}

func (blockchain *BlockChain) PackSignal(height int64) *Block {
	blockchain.PackLock.Lock()
	defer blockchain.PackLock.Unlock()
	if blockchain.Status != StartPackStatus {
		if !blockchain.BlockManager.GetBlockStatusByHeight(height, int64(blockchain.BlockInterval)) {
			return nil
		}
		blockchain.Status = StartPackStatus
		blockchain.BlockManager.SetBlockStatusByHeight(height, time.Now().UnixNano())
		defer func() {
			if r := recover(); r != nil {
				log.GetLogInst().LogCrit("Panic while pack. %v", r)
			}
			blockchain.Status = InitStatus
		}()
		log.GetLogInst().LogInfo("Start pack block at height %d .\n", blockchain.GetLastHeight()+1)
		log.GetLogInst().LogDebug("Start pack block at height %d .\n", blockchain.GetLastHeight()+1)
		block := blockchain.WaitAndPack()
		log.GetLogInst().LogInfo("Packed a block at height %d, block info: %s .\n", blockchain.GetLastHeight()+1, string(block.Bytes()))
		log.GetLogInst().LogDebug("Packed a block at height %d, block info: %s .\n", blockchain.GetLastHeight()+1, string(block.Bytes()))
		return block
	}
	return nil
}

func (blockchain *BlockChain) PackHeightValidate(height int64) bool {
	if blockchain.GetLastHeight()+1 != height {
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
	if height > blockchain.GetLastHeight() {
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

func (blockchain *BlockChain) SaveBlock(block *Block) {
	blockchain.Locker.Lock()
	defer blockchain.Locker.Unlock()
	if blockchain.GetLastHeight()+1 == block.Height {
		fmt.Println("Saving block to database.")
		db.GetDBInst().Set(block.Hash(), block.Data())
		data, _ := json.Marshal(block)
		db.GetDBInst().Set(blockchain.GetBlockByHeightKey(block.Height), data)
		db.GetDBInst().Set(blockchain.CurrentBlockKey(), data)
		blockchain.SetLastBlock(block)
		blockchain.SetLastHeight(block.Height)
		fmt.Println("Save block to database succeed.")
	}
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

func (blockchain *BlockChain) PackTime() time.Duration {
	d := blockchain.BlockInterval / 3
	if blockchain.BlockInterval > 3*time.Second {
		d = blockchain.BlockInterval - 2*time.Second
	}
	return d
}

func (blockchain *BlockChain) WaitAndPack() *Block {
	// 打包10500个交易大概需要0.95秒
	eventTimeout := time.After(blockchain.PackTime())
	round := &i_consensus.Round{
		Peers:        param.MainChainDPosNode,
		CurrentIndex: 0,
	}
	if blockchain.GetLastHeight() != 0 {
		round = blockchain.GetLastBlock().GetRound().MyRound(blockchain.GetLastBlock().CurrentHash)
	}
	log.GetLogInst().LogDebug("")
	block := NewBlock(blockchain.GetLastBlock(), round)
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
	bodyData := block.BlockBody.Bytes()
	block.Body = crypto.Sha3_256(bodyData)
	db.GetDBInst().Set(block.Body, bodyData)
	block.UpdateMPTPlusRoot()
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

func (blockchain *BlockChain) BlockFromPeer(cLog *context_log.ContextLog, block Block) bool {
	fmt.Printf("Validating block from peer, block info: %s, block.Hash=%s \n", string(block.Bytes()), hex.EncodeToString(block.Hash()))
	if err := block.Validate(); err != nil {
		cLog.Log("InvalidBlock", true)
		fmt.Printf("Block validate failed, %s. \n", err.Error())
		return false
	}
	// 1500是毫秒和纳秒的单位乘以2/3计算得来的
	if time.Now().UnixNano()/1e6-block.Timestamp > int64(blockchain.BlockInterval/1500) {
		fmt.Printf("time.Now=%d, block.Time=%d, block.Interval=%d \n", time.Now().UnixNano()/1e6, block.Timestamp, int64(blockchain.BlockInterval/1500))
		fmt.Println("Block timestamp is more than 2/3 block interval, abort vote.")
		return false
	}
	if !blockchain.GetLastBlock().ValidateNextBlock(block, blockchain.BlockInterval) {
		fmt.Println("This block from peer can not recover by last block, abort.")
		return false
	}
	return true
}

func (blockchain BlockChain) NewTransaction(tx *common.Transaction) bool {
	from, _ := hex.DecodeString(tx.From)
	if account, err := blockchain.GetLastBlock().GetAccount(from); err == nil && account != nil {
		status := pool.Block
		if account.Nonce+1 == tx.Nonce {
			status = pool.Ready
		}
		blockchain.Pool.ParkTx(tx, status)
		return true
	}
	return false
}
