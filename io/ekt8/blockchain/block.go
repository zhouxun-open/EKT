package blockchain

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/EducationEKT/EKT/io/ekt8/MPTPlus"
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"github.com/EducationEKT/EKT/io/ekt8/crypto"
)

var currentBlock *Block = nil

type Block struct {
	Height       int64        `json:"height"`
	Nonce        int64        `json:"nonce"`
	PreviousHash []byte       `json:"previousHash"`
	CurrentHash  []byte       `json:"currentHash"`
	Locker       sync.RWMutex `json:"-"`
	StatTree     *MPTPlus.MTP `json:"-"`
	TxTree       *MPTPlus.MTP `json:"-"`
	EventTree    *MPTPlus.MTP `json:"-"`
}

func (block *Block) String() string {
	return fmt.Sprintf(`{"height": %d, "statRoot": "%s", "txRoot": "%s", "eventRoot": "%s", "nonce": %d, "previousHash": "%s"}`,
		block.Height, block.StatTree.Root, block.TxTree.Root, block.EventTree.Root, block.Nonce, block.PreviousHash)
}

func (block *Block) Bytes() []byte {
	return []byte(block.String())
}

func (block *Block) Hash() []byte {
	return crypto.Sha3_256(block.Hash())
}

func (block *Block) NewNonce() {
	block.Nonce++
}

func (block Block) Validate() error {
	if !bytes.Equal(block.Hash(), block.CurrentHash) {
		return errors.New("Invalid Hash")
	}
	return nil
}

func (block *Block) GetAccount(address []byte) (*common.Account, error) {
	value, err := block.StatTree.GetValue(address)
	if err != nil {
		return nil, err
	}
	var account common.Account
	err = json.Unmarshal(value, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (block *Block) ExistAddress(address []byte) bool {
	return block.StatTree.ContainsKey(address)
}

func (block *Block) NewAccount(address []byte, pubKey []byte) error {
	account := &common.Account{
		Address:    address,
		PublickKey: pubKey,
		Amount:     0,
		Nonce:      0,
	}
	value, err := json.Marshal(account)
	if err != nil {
		return err
	}
	err = block.StatTree.MustInsert(address, value)
	return err
}
