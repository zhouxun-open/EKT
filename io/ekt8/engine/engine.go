package engine

import (
	"errors"

	"github.com/EducationEKT/EKT/io/ekt8/blockchain"
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
)

var mainBlockChain blockchain.BlockChain

func init() {
	mainBlockChain = blockchain.BlockChain{[]byte{(byte)(1 & 0xFF)}}
}

type Engine struct {
	blockChain  *blockchain.BlockChain
	Pack        chan bool
	Transaction chan *common.Transaction
	Status      int // 100 正在计算MTProot, 150停止计算root,开始计算blockHash
}

func (engine *Engine) NewTransaction(transaction *common.Transaction) error {
	if engine.Status == 100 {
		block, err := engine.blockChain.CurrentBlock()
		if err != nil {
			return err
		}
		Transaction(block, transaction)
	}
	return errors.New("Wait next block")
}
