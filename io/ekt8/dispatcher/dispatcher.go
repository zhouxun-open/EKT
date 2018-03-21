package dispatcher

import (
	"encoding/hex"

	"github.com/EducationEKT/EKT/io/ekt8/blockchain"
	"github.com/EducationEKT/EKT/io/ekt8/blockchain_manager"
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"github.com/EducationEKT/EKT/io/ekt8/event"
	"github.com/EducationEKT/EKT/io/ekt8/tx_pool"
)

var dispatcher DefaultDispatcher

func init() {
	dispatcher = DefaultDispatcher{}
}

func NewTransaction(transaction common.Transaction) {
	blockchain_manager.MainBlockChainConsensus.NewTransaction(transaction)
}

type IDispatcher interface {
	NewTransaction(transaction *common.Transaction)
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

func (dispatcher DefaultDispatcher) NewTransaction(transaction *common.Transaction) {
	// TODO move validate to blockchain
	if err := transaction.Validate(); err != nil {
		return
	}
	blockChain := dispatcher.GetBackBoneBlockChain()
	if blockChain.GetStatus() == 100 {
		if block := blockChain.CurrentBlock; block != nil {
			address, _ := hex.DecodeString(transaction.From)
			account, _ := block.GetAccount(address)
			if transaction.Nonce <= account.GetNonce() {
				return
			} else if transaction.Nonce-account.GetNonce() > 1 {
				tx_pool.GetTxPool().Park(transaction, tx_pool.Block)
			} else {
				toAddress, _ := hex.DecodeString(transaction.To)
				if !block.ExistAddress(toAddress) {
					return
				}
				block.NewTransaction(transaction)
			}
		}
	}
}

func (dispatcher DefaultDispatcher) NewEvent(evt *event.Event) {
	if !evt.ValidateEvent() {
		return
	}
	if evt.EventType == event.NewAccountEvent {
		accountParam := (evt.EventParam).(event.NewAccountParam)
		block := blockchain_manager.MainBlockChain.CurrentBlock
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
