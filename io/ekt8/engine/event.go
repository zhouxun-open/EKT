package engine

import (
	"encoding/json"
	"errors"

	"github.com/EducationEKT/EKT/io/ekt8/MPTPlus"
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"github.com/EducationEKT/EKT/io/ekt8/db"
	pkg_event "github.com/EducationEKT/EKT/io/ekt8/event"
)

func GetAccount(address []byte) (*common.Account, error) {
	block, err := mainBlockChain.CurrentBlock()
	if err != nil {
		return nil, err
	}
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

func ExistAddress(address []byte) bool {
	block, err := mainBlockChain.CurrentBlock()
	if err != nil {
		return false
	}
	statRoot := block.StatRoot
	statTree := MPTPlus.MTP_Tree(db.GetDBInst(), statRoot)
	return statTree.ContainsKey(address)
}

func NewAccount(address []byte, pubKey []byte) error {
	block, err := mainBlockChain.CurrentBlock()
	if err != nil {
		return err
	}
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

func NewEvent(event pkg_event.Event) error {
	if !event.ValidateEvent() {
		return errors.New("Invalid event")
	}
	if event.EventType == pkg_event.NewAccountEvent {
		param := event.EventParam.(pkg_event.NewAccountParam)
		if ExistAddress(param.Address) {
			return errors.New("Address Exist")
		} else {
			if err := NewAccount(param.Address, param.PubKey); err == nil {

			}
		}
	}
	return nil
}
