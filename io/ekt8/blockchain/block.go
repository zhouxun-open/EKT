package blockchain

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/EducationEKT/EKT/io/ekt8/MPTPlus"
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"github.com/EducationEKT/EKT/io/ekt8/crypto"
	"github.com/EducationEKT/EKT/io/ekt8/db"
)

type Block struct {
	Height       int64  `json:"height"`
	StatRoot     []byte `json:"statRoot"`
	TxRoot       []byte `json:"txRoot"`
	EventRoot    []byte `json:"eventRoot"`
	Nonce        int64  `json:"nonce"`
	PreviousHash []byte `json:"previousHash"`
	CurrentHash  []byte `json:"currentHash"`
}

func (block *Block) String() string {
	return fmt.Sprintf(`{"height": %d, "statRoot": "%s", "txRoot": "%s", "eventRoot": "%s", "nonce": %d, "previousHash": "%s"}`,
		block.Height, block.StatRoot, block.TxRoot, block.Nonce, block.PreviousHash)
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
	statRoot := block.StatRoot
	statTree := MPTPlus.MTP_Tree(db.GetDBInst(), statRoot)
	value, err := statTree.GetValue(address)
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
	statRoot := block.StatRoot
	statTree := MPTPlus.MTP_Tree(db.GetDBInst(), statRoot)
	return statTree.ContainsKey(address)
}

func (block *Block) NewAccount(address []byte, pubKey []byte) error {
	statRoot := block.StatRoot
	trie := MPTPlus.MTP_Tree(db.GetDBInst(), statRoot)
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
	err = trie.MustInsert(address, value)
	block.StatRoot = trie.Root
	return err
}
