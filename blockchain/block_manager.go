package blockchain

import (
	"encoding/hex"
	"sync"
	"time"
)

const (
	BLOCK_TO_BE_HANDLE = 100
	BODY_SAVED         = 101
	ERROR_BLOCK        = 400
	BLOCK_SAVED        = 200
)

// 内部操作不加lock，外部在需要加锁的地方加锁，保证操作的原子性
type BlockManager struct {
	Blocks        map[string]*Block
	BlockStatus   map[string]int  // 根据区块hash计算，主要是从peer来的区块 100：待处理 	101：已经处理成功，未写入区块 	400：错误的区块头 		200：处理成功，已经写入区块
	HeightManager map[int64]int64 // 根据block的height进行计算，主要是防止内部多次进行打包 100代表未打包，101代表已打包
	HeightVote    map[int64]int64 //上次在某个高度的投票时间，防止重复投票
	locker        sync.RWMutex
}

func NewBlockManager() *BlockManager {
	return &BlockManager{
		Blocks:        make(map[string]*Block),
		BlockStatus:   make(map[string]int),
		HeightManager: make(map[int64]int64),
		HeightVote:    make(map[int64]int64),
		locker:        sync.RWMutex{},
	}
}

// 获取指定区块的状态， -1表示不存在
func (manager *BlockManager) GetBlockStatus(hash []byte) int {
	status, exist := manager.BlockStatus[hex.EncodeToString(hash)]
	if !exist {
		return -1
	}
	return status
}

func (manager *BlockManager) GetVoteTime(height int64) int64 {
	manager.RLock()
	defer manager.RUnlock()
	time, exist := manager.HeightVote[height]
	if exist {
		return time
	}
	return -1
}

func (manager *BlockManager) SetVoteTime(height int64, time int64) {
	manager.Lock()
	defer manager.Unlock()
	manager.HeightVote[height] = time
}

// 根据区块高度判断自己是否可以对此高度进行打包
// 一个区块在2个interval内不可以对同一个高度的区块进行打包
func (manager *BlockManager) GetBlockStatusByHeight(height, interval int64) bool {
	manager.locker.RLock()
	defer manager.locker.RUnlock()
	status, exist := manager.HeightManager[height]
	if !exist {
		return true
	}
	return status < (time.Now().UnixNano()-interval*2)/1e6
}

func (manager *BlockManager) SetBlockStatusByHeight(height, nanoSecond int64) {
	manager.Lock()
	defer manager.Unlock()
	manager.HeightManager[height] = nanoSecond
}

//将指定区块插入，默认是100
func (manager *BlockManager) Insert(block *Block) {
	hash := hex.EncodeToString(block.CurrentHash)
	if _, exist := manager.Blocks[hash]; exist {
		return
	} else {
		manager.Blocks[hash] = block
		manager.BlockStatus[hash] = 100
	}
}

func (manager *BlockManager) RLock() {
	manager.locker.RLock()
}

func (manager *BlockManager) Lock() {
	manager.locker.Lock()
}

func (manager *BlockManager) RUnlock() {
	manager.locker.RUnlock()
}

func (manager *BlockManager) Unlock() {
	manager.locker.Unlock()
}
