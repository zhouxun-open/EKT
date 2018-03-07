package blockchain

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/EducationEKT/EKT/io/ekt8/db"
)

const (
	CurrentBlock = "CurrentBlock"
)

type BlockChain struct {
	ChainId []byte
}

func (blockChain *BlockChain) NewBlock(block Block) error {
	if err := block.Validate(); err != nil {
		return err
	}
	lastBlock, err := blockChain.CurrentBlock()
	if err != nil {
		return err
	}
	if lastBlock.Height > block.Height {
		return errors.New("heigth exist")
	}
	err = db.GetDBInst().Set(block.CurrentHash, block.Hash())
	if err != nil {
		return err
	}
	value, err := json.Marshal(block)
	if err != nil {
		return err
	}
	return db.GetDBInst().Set(blockChain.CurrentBlockKey(), value)
}

func (blockChain BlockChain) CurrentBlock() (*Block, error) {
	blockValue, err := db.GetDBInst().Get(blockChain.CurrentBlockKey())
	if err != nil {
		return nil, err
	}
	var block Block
	err = json.Unmarshal(blockValue, &block)
	return &block, err
}

func (blockChain BlockChain) CurrentBlockKey() []byte {
	buffer := bytes.Buffer{}
	buffer.WriteString(CurrentBlock)
	buffer.Write(blockChain.ChainId)
	return buffer.Bytes()
}
