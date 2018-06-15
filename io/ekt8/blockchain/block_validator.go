package blockchain

import (
	"bytes"
	"fmt"
	"sync"
	"time"
)

type BlockValidator struct {
	Blocks map[int64][]*Block
	Block  *Block
	/*
		status = 0   	初始化，状态未知，只有status状态为0才可以修改status的值为成功
		status = 200	表示校验成功
		status=500 		表示校验失败
	*/
	Status   int
	Locker   sync.RWMutex
	Block_ch chan *Block
}

// 此方法用来校验是否收到相同高度不同的block
func (validator BlockValidator) NewBlock(block *Block) bool {
	validator.Locker.Lock()
	defer validator.Locker.Unlock()
	blocks, exist := validator.Blocks[block.Height]
	//如果当前节点没有收到过其他节点的block，返回校验成功
	if !exist || 0 == len(blocks) {
		return true
	}

	for _, _block := range blocks {
		// 如果收到了两个不同的block，会先判断是否是同一个节点打包的
		// 如果是同一个节点打包了不同的block，则放弃当前区块，进行下一个区块的打包，并对此行为进行记录
		// 如果不是同一个节点打包的不同的block，则判断哪个节点拥有权限，对对应的区块进行投票，半数以上同意的区块可以写入区块链中
		// 这一部分会放在上一层blockchain进行判断，因为需要读取上一个区块的信息，所以此模块根据其他模块去判断
		if !bytes.Equal(_block.Hash(), block.Hash()) {
			// 相同节点的打包
			// note： 是否是某一个节点的校验在network层会进行校验，如果签名校验失败会进行拦截
			if _block.GetRound().Peers[_block.GetRound().CurrentIndex].Equal(block.GetRound().Peers[block.GetRound().CurrentIndex]) {

			} else {
				continue
			}
		}
	}
	return true
}

func NewBlockchainValidator(block *Block) *BlockValidator {
	return &BlockValidator{
		Block:    block,
		Status:   0,
		Locker:   sync.RWMutex{},
		Block_ch: make(chan *Block),
	}
}

func (blockValidator BlockValidator) Run() {
	timeout := time.After(1 * time.Second)
	for {
		select {
		case <-timeout:
			return
		case block := <-blockValidator.Block_ch:
			if !bytes.Equal(block.Hash(), blockValidator.Block.Hash()) {
				blockValidator.Fail()
			}
			fmt.Println(block.Height)
		}
	}
}

func (blockValidator BlockValidator) Fail() {
	blockValidator.Locker.Lock()
	defer blockValidator.Locker.Unlock()
	blockValidator.Status = 500
}
