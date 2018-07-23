package blockchain

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"errors"

	"github.com/EducationEKT/EKT/core/common"
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

func (chain *BlockChain) GetLastBlock() *Block {
	chain.currentLocker.RLock()
	defer chain.currentLocker.RUnlock()
	return chain.currentBlock
}

func (chain *BlockChain) SetLastBlock(block *Block) {
	chain.currentLocker.Lock()
	defer chain.currentLocker.Unlock()
	chain.currentBlock = block
	chain.currentHeight = block.Height
}

func (chain *BlockChain) GetLastHeight() int64 {
	chain.currentLocker.RLock()
	defer chain.currentLocker.RUnlock()
	return chain.currentHeight
}

func (chain *BlockChain) PackSignal(height int64) *Block {
	chain.PackLock.Lock()
	defer chain.PackLock.Unlock()
	if chain.Status != StartPackStatus {
		defer func() {
			if r := recover(); r != nil {
				log.Crit("Panic while pack. %v", r)
			}
			chain.Status = InitStatus
		}()
		if !chain.PackHeightValidate(height) {
			log.Info("This height is packed within an interval, return nil.")
			return nil
		}
		log.Info("Start pack block at height %d .\n", chain.GetLastHeight()+1)
		log.Debug("Start pack block at height %d .\n", chain.GetLastHeight()+1)
		block := chain.WaitAndPack()
		log.Info("Packed a block at height %d, block info: %s .\n", chain.GetLastHeight()+1, string(block.Bytes()))
		log.Debug("Packed a block at height %d, block info: %s .\n", chain.GetLastHeight()+1, string(block.Bytes()))
		return block
	}
	return nil
}

func (chain *BlockChain) PackHeightValidate(height int64) bool {
	if chain.GetLastHeight()+1 != height {
		return false
	}
	if !chain.BlockManager.GetBlockStatusByHeight(height, int64(chain.BlockInterval)) {
		return false
	}
	return true
}

func (chain *BlockChain) GetBlockByHeight(height int64) (*Block, error) {
	if height > chain.GetLastHeight() {
		return nil, errors.New("Invalid height")
	}
	key := chain.GetBlockByHeightKey(height)
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

func (chain *BlockChain) GetBlockByHeightKey(height int64) []byte {
	return []byte(fmt.Sprint(`GetBlockByHeight: _%d_%d`, chain.ChainId, height))
}

func (chain *BlockChain) SaveBlock(block *Block) {
	chain.Locker.Lock()
	defer chain.Locker.Unlock()
	if chain.GetLastHeight()+1 == block.Height {
		log.Info("Saving block to database.")
		db.GetDBInst().Set(block.Hash(), block.Data())
		data, _ := json.Marshal(block)
		db.GetDBInst().Set(chain.GetBlockByHeightKey(block.Height), data)
		db.GetDBInst().Set(chain.CurrentBlockKey(), data)
		chain.SetLastBlock(block)
		log.Info("Saved block to database.")
	}
}

func (chain *BlockChain) LastBlock() (*Block, error) {
	var err error = nil
	var block *Block
	if currentBlock == nil {
		key := chain.CurrentBlockKey()
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

func (chain *BlockChain) CurrentBlockKey() []byte {
	return []byte(fmt.Sprintf("CurrentBlockKey_%d", chain.ChainId))
}

func (chain *BlockChain) PackTime() time.Duration {
	d := chain.BlockInterval / 3
	if chain.BlockInterval > 3*time.Second {
		d = chain.BlockInterval - 2*time.Second
	}
	return d
}

func (chain *BlockChain) WaitAndPack() *Block {
	// 打包10500个交易大概需要0.95秒
	eventTimeout := time.After(chain.PackTime())
	block := NewBlock(chain.GetLastBlock())
	if block.Fee <= 0 {
		block.Fee = chain.Fee
	}
	log.Info("Packing transaction and other events.")
	start := time.Now().UnixNano()
	started := false
	numTx := 0
	for {
		flag := false
		select {
		case <-eventTimeout:
			flag = true
			break
		default:
			event := chain.Pool.Fetch()
			if event != nil {
				if !started {
					started = true
					start = time.Now().UnixNano()
				}
				tx, ok := event.(common.Transaction)
				if ok {
					numTx++
					block.NewTransaction(tx, tx.Fee)
					block.BlockBody.AddEvent(event)
				}
			}
		}
		if flag {
			break
		}
	}
	end := time.Now().UnixNano()
	fmt.Printf("Total tx: %d, startTime: %d, endTime: %d, Total time: %d ns. \n", numTx, start/1e6, end/1e6, end-start)
	bodyData := block.BlockBody.Bytes()
	block.Body = crypto.Sha3_256(bodyData)
	db.GetDBInst().Set(block.Body, bodyData)
	block.UpdateMPTPlusRoot()
	return block
}

// 当区块写入区块时，notify交易池，一些nonce比较大的交易可以进行打包
func (chain *BlockChain) NotifyPool(block *Block) {
	if block.BlockBody == nil {
		return
	}
	block.BlockBody.Events.Range(func(key, value interface{}) bool {
		address, ok1 := key.(string)
		list, ok2 := value.([]string)
		if ok1 && ok2 && len(list) > 0 {
			for _, eventId := range list {
				chain.Pool.Notify(address, eventId)
			}
		}
		return true
	})
}

func (chain *BlockChain) NewTransaction(tx common.Transaction) bool {
	block := chain.GetLastBlock()
	account, err := block.GetAccount(tx.GetFrom())
	if err != nil {
		return false
	}
	if account.GetNonce() >= tx.GetNonce() {
		return false
	}
	result := true
	if account.GetNonce()+1 == tx.GetNonce() {
		result = chain.Pool.Park(tx, pool.Ready)
	} else {
		result = chain.Pool.Park(tx, pool.Block)
	}
	return result
}

func (chain *BlockChain) BlockFromPeer(ctxlog *ctxlog.ContextLog, block Block) bool {
	log.Info("Validating block from peer, block info: %s, block.Hash=%s \n", string(block.Bytes()), hex.EncodeToString(block.Hash()))
	if err := block.Validate(ctxlog); err != nil {
		ctxlog.Log("InvalidReason", err.Error())
		return false
	}

	// 1500是毫秒和纳秒的单位乘以2/3计算得来的
	if time.Now().UnixNano()/1e6-block.Timestamp > int64(chain.BlockInterval/1500) {
		ctxlog.Log("Invalid timestamp", true)
		log.Info("time.Now=%d, block.Time=%d, block.Interval=%d \n", time.Now().UnixNano()/1e6, block.Timestamp, int64(chain.BlockInterval/1500))
		log.Info("Block timestamp is more than 2/3 block interval, abort vote.")
		return false
	}

	if !chain.GetLastBlock().ValidateNextBlock(block, chain.BlockInterval) {
		log.Info("This block from peer can not recover by last block, abort.")
		return false
	}
	return true
}
