package dispatcher

import (
	"encoding/hex"

	"github.com/EducationEKT/EKT/io/ekt8/blockchain"
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"github.com/EducationEKT/EKT/io/ekt8/event"
)

type IDispatcher interface {
	NewTransaction(transaction *common.Transaction)
	NewEvent(event *event.Event)
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

}

func (dispatcher DefaultDispatcher) NewEvent(event *event.Event) {
}
