package blockchain

import (
	"bytes"
	"fmt"
	"sync"
	"time"
)

type BlockValidator struct {
	Block *Block
	/*
		status = 0   	初始化，状态未知，只有status状态为0才可以修改status的值为成功
		status = 200	表示校验成功
		status=500 		表示校验失败
	*/
	Status   int
	Locker   sync.RWMutex
	Block_ch chan *Block
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
	timeout := time.After(2 * time.Second)
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
