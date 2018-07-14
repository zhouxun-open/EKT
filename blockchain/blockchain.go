package blockchain

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"errors"

	"github.com/EducationEKT/EKT/crypto"
	"github.com/EducationEKT/EKT/ctxlog"
	"github.com/EducationEKT/EKT/db"
	"github.com/EducationEKT/EKT/i_consensus"
	"github.com/EducationEKT/EKT/log"
	"github.com/EducationEKT/EKT/pool"
)

var BackboneChainId int64 = 1

const (
	BackboneConsensus     = i_consensus.DPOS
	BackboneBlockInterval = 10 * time.Second
	BackboneChainFee      = 510000
)

const (
	InitStatus      = 0
	StartPackStatus = 100
)

type BlockChain struct {
	ChainId       int64
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

func NewBlockChain(chainId int64, consensusType i_consensus.ConsensusType, fee int64, difficulty []byte, interval time.Duration) *BlockChain {
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
	blockchain.currentHeight = block.Height
}

func (blockchain *BlockChain) GetLastHeight() int64 {
	blockchain.currentLocker.RLock()
	defer blockchain.currentLocker.RUnlock()
	return blockchain.currentHeight
}

func (blockchain *BlockChain) PackSignal(height int64) *Block {
	blockchain.PackLock.Lock()
	defer blockchain.PackLock.Unlock()
	if blockchain.Status != StartPackStatus {
		defer func() {
			if r := recover(); r != nil {
				log.Crit("Panic while pack. %v", r)
			}
			blockchain.Status = InitStatus
		}()
		if !blockchain.PackHeightValidate(height) {
			log.Info("This height is packed within an interval, return nil.")
			return nil
		}
		log.Info("Start pack block at height %d .\n", blockchain.GetLastHeight()+1)
		log.Debug("Start pack block at height %d .\n", blockchain.GetLastHeight()+1)
		block := blockchain.WaitAndPack()
		log.Info("Packed a block at height %d, block info: %s .\n", blockchain.GetLastHeight()+1, string(block.Bytes()))
		log.Debug("Packed a block at height %d, block info: %s .\n", blockchain.GetLastHeight()+1, string(block.Bytes()))
		return block
	}
	return nil
}

func (blockchain *BlockChain) PackHeightValidate(height int64) bool {
	if blockchain.GetLastHeight()+1 != height {
		return false
	}
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
	return []byte(fmt.Sprint(`GetBlockByHeight: _%d_%d`, blockchain.ChainId, height))
}

func (blockchain *BlockChain) SaveBlock(block *Block) {
	blockchain.Locker.Lock()
	defer blockchain.Locker.Unlock()
	if blockchain.GetLastHeight()+1 == block.Height {
		log.Info("Saving block to database.")
		db.GetDBInst().Set(block.Hash(), block.Data())
		data, _ := json.Marshal(block)
		db.GetDBInst().Set(blockchain.GetBlockByHeightKey(block.Height), data)
		db.GetDBInst().Set(blockchain.CurrentBlockKey(), data)
		blockchain.SetLastBlock(block)
		log.Info("Saved block to database.")
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
	return []byte(fmt.Sprintf("CurrentBlockKey_%d", blockchain.ChainId))
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
	block := NewBlock(blockchain.GetLastBlock())
	if block.Fee <= 0 {
		block.Fee = blockchain.Fee
	}
	log.Info("Packing transaction and other events.")
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
				go blockchain.Pool.Notify(tx.TransactionId())
				log := ctxlog.NewContextLog("BlockFromTxPool")
				defer log.Finish()
				log.Log("tx", tx)
				log.Log("block.StatRoot_p", block.StatTree.Root)
				fee := tx.Fee
				if fee < block.Fee {
					fee = block.Fee
				}
				txResult := block.NewTransaction(log, tx, fee)
				log.Log("fee", fee)
				log.Log("txResult", txResult)
				log.Log("block.StatRoot_a", block.StatTree.Root)
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

func (blockchain *BlockChain) BlockFromPeer(ctxlog *ctxlog.ContextLog, block Block) bool {
	log.Info("Validating block from peer, block info: %s, block.Hash=%s \n", string(block.Bytes()), hex.EncodeToString(block.Hash()))
	if err := block.Validate(ctxlog); err != nil {
		ctxlog.Log("InvalidReason", err.Error())
		return false
	}

	// 1500是毫秒和纳秒的单位乘以2/3计算得来的
	if time.Now().UnixNano()/1e6-block.Timestamp > int64(blockchain.BlockInterval/1500) {
		ctxlog.Log("Invalid timestamp", true)
		log.Info("time.Now=%d, block.Time=%d, block.Interval=%d \n", time.Now().UnixNano()/1e6, block.Timestamp, int64(blockchain.BlockInterval/1500))
		log.Info("Block timestamp is more than 2/3 block interval, abort vote.")
		return false
	}

	if !blockchain.GetLastBlock().ValidateNextBlock(block, blockchain.BlockInterval) {
		log.Info("This block from peer can not recover by last block, abort.")
		return false
	}
	return true
}
