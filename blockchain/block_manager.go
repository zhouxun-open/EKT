package blockchain

import (
	"encoding/hex"
	"sync"
	"time"
)

const (
	// 1xx未处理
	BLOCK_TO_BE_HANDLE = 100

	// 2xx 处理成功待后续处理
	BLOCK_VALID = 201
	BLOCK_VOTED = 202

	// 3xx 错误的区块
	BLOCK_ERROR_START          = 300
	BLOCK_ERROR_PACK_TIME      = 301
	BLOCK_ERROR_BROADCAST_TIME = 302
	BLOCK_ERROR_HASH           = 303
	BLOCK_ERROR_SIGN           = 304
	BLOCK_ERROR_BODY           = 305
	BLOCK_ERROR_END            = 399

	// 400已经写入区块链
	BLOCK_SAVED = 400
)

// 内部操作不加lock，外部在需要加锁的地方加锁，保证操作的原子性
type BlockManager struct {
	Blocks        *sync.Map
	BlockStatus   *sync.Map // 根据区块hash计算，主要是从peer来的区块 100：待处理 	101：已经处理成功，未写入区块 	400：错误的区块头 		200：处理成功，已经写入区块
	HeightManager *sync.Map // 根据block的height进行计算，主要是防止内部多次进行打包 100代表未打包，101代表已打包
	HeightVote    *sync.Map //上次在某个高度的投票时间，防止重复投票
}

func NewBlockManager() *BlockManager {
	return &BlockManager{
		Blocks:        &sync.Map{},
		BlockStatus:   &sync.Map{},
		HeightManager: &sync.Map{},
		HeightVote:    &sync.Map{},
	}
}

// 获取指定区块的状态， -1表示不存在
func (manager *BlockManager) GetBlockStatus(hash []byte) int {
	status, exist := manager.BlockStatus.Load(hex.EncodeToString(hash))
	if !exist {
		return -1
	}
	return status.(int)
}

func (manager *BlockManager) SetBlockStatus(hash []byte, status int) {
	_status := manager.GetBlockStatus(hash)
	if _status < status {
		manager.BlockStatus.Store(hex.EncodeToString(hash), status)
	}
}

func (manager *BlockManager) GetVoteTime(height int64) int64 {
	time, exist := manager.HeightVote.Load(height)
	if exist {
		return time.(int64)
	}
	return -1
}

func (manager *BlockManager) SetVoteTime(height int64, time int64) {
	manager.HeightVote.Store(height, time)
}

// 根据区块高度判断自己是否可以对此高度进行打包
// 一个区块在1个interval内不可以对同一个高度的区块进行打包
func (manager *BlockManager) GetBlockStatusByHeight(height, interval int64) bool {
	t, exist := manager.HeightManager.Load(height)
	if !exist {
		return true
	}
	return t.(int64)+interval/1e6 < time.Now().UnixNano()/1e6
}

func (manager *BlockManager) SetBlockStatusByHeight(height, nanoSecond int64) {
	manager.HeightManager.Store(height, nanoSecond)
}

//将指定区块插入，默认是100
func (manager *BlockManager) Insert(block Block) {
	hash := hex.EncodeToString(block.CurrentHash)
	if _, exist := manager.Blocks.Load(hash); exist {
		return
	} else {
		manager.Blocks.Store(hash, block)
		manager.BlockStatus.Store(hash, BLOCK_TO_BE_HANDLE)
	}
}

func (manager *BlockManager) GetBlock(hash []byte) (Block, bool) {
	block, exist := manager.Blocks.Load(hex.EncodeToString(hash))
	if !exist {
		return Block{}, exist
	}
	return block.(Block), exist
}
