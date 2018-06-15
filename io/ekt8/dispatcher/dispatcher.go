package dispatcher

import (
	"encoding/hex"
	"errors"

	"github.com/EducationEKT/EKT/io/ekt8/blockchain"
	"github.com/EducationEKT/EKT/io/ekt8/blockchain_manager"
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"github.com/EducationEKT/EKT/io/ekt8/event"
)

var dispatcher DefaultDispatcher

func init() {
	dispatcher = DefaultDispatcher{}
}

type IDispatcher interface {
	NewTransaction(transaction *common.Transaction) error
	NewEvent(event *event.Event)
}

func GetDisPatcher() IDispatcher {
	return dispatcher
}

type DefaultDispatcher struct {
	blockChains map[string]*blockchain.BlockChain
	openFunc    map[string]*blockchain.ChainFunc
}

func (dispatcher DefaultDispatcher) GetBlockChain(chainId []byte) (*blockchain.BlockChain, bool) {
	blockChain, exist := dispatcher.blockChains[hex.EncodeToString(chainId)]
	return blockChain, exist
}

func (dispacher DefaultDispatcher) GetBackBoneBlockChain() *blockchain.BlockChain {
	blockChain := dispacher.blockChains[hex.EncodeToString(blockchain.BackboneChainId)]
	return blockChain
}

func NewTransaction(transaction *common.Transaction) error {
	// 主币的tokenAddress为空
	if transaction.TokenAddress != "" {
		tokenAddress, err := hex.DecodeString(transaction.TokenAddress)
		if err != nil {
			return err
		}
		currentBlock := blockchain_manager.GetMainChain().GetLastBlock()
		var token common.Token
		err = currentBlock.TokenTree.GetInterfaceValue(tokenAddress, &token)
		if err != nil || token.Name == "" || token.Decimals <= 0 || token.Total <= 0 {
			return err
		}
	}
	if !transaction.Validate() {
		return errors.New("error signature")
	}
	if !blockchain_manager.GetMainChain().NewTransaction(transaction) {
		return errors.New("error transaction")
	}
	return nil
}

func (dispatcher DefaultDispatcher) NewTransaction(transaction *common.Transaction) error {
	// 主币的tokenAddress为空
	if transaction.TokenAddress != "" {
		tokenAddress, err := hex.DecodeString(transaction.TokenAddress)
		if err != nil {
			return err
		}
		currentBlock := dispatcher.GetBackBoneBlockChain().GetLastBlock()
		var token common.Token
		err = currentBlock.TokenTree.GetInterfaceValue(tokenAddress, &token)
		if err != nil || token.Name == "" || token.Decimals <= 0 || token.Total <= 0 {
			return err
		}
	}
	if !transaction.Validate() {
		return errors.New("error signature")
	}
	if !blockchain_manager.GetMainChain().NewTransaction(transaction) {
		return errors.New("error transaction")
	}
	return nil
}

func (dispatcher DefaultDispatcher) NewEvent(evt *event.Event) {
	if !evt.ValidateEvent() {
		return
	}
	if evt.EventType == event.NewAccountEvent {
		accountParam := (evt.EventParam).(event.NewAccountParam)
		block := blockchain_manager.MainBlockChain.GetLastBlock()
		address, err := hex.DecodeString(accountParam.Address)
		if err != nil && !block.ExistAddress(address) {
			pubKey, err := hex.DecodeString(accountParam.PubKey)
			if err != nil {
				return
			}
			block.CreateAccount(address, pubKey)
		}
	}
}
