package stat

import (
	"github.com/EducationEKT/EKT/io/ekt8/MPTPlus"
	"github.com/EducationEKT/EKT/io/ekt8/blockchain"
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"github.com/EducationEKT/EKT/io/ekt8/db"
)

var block *blockchain.Block

func NewBlock(block_ *blockchain.Block) {
	block = block_
}

func GetAccount(address []byte) (*common.Account, error) {
	if block == nil {
		return nil, WaitingSyncErr
	}
	block.StatTree = MPTPlus.MTP_Tree(db.GetDBInst(), block.StatRoot)
	account, err := block.GetAccount(address)
	if err != nil {
		return nil, NoAccountErr
	}
	return account, nil
}
