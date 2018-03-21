package tx_pool

import (
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"sort"
)

const (
	Block = iota
	Queue
)

var txPool TxPool

type BlockQueue [] *common.Transaction

func init() {
	txPool = TxPool{}
}

type TxPool struct {
	ready map[string]*common.Transaction
	wait BlockQueue
}

func GetTxPool() TxPool {
	return txPool
}

//func (txPool TxPool) Park(tx *common.Transaction) {
//	if txPool.ready[tx.TransactionId] != nil {
//		return
//	}
//	txPool.ready[tx.TransactionId] = tx
//}

func (txPool TxPool) Park(tx *common.Transaction ,reason int){

}

func (TxPool TxPool) Notify(tx *common.Transaction ){

}

func (txPool TxPool) BatchNotify(txs []*common.Transaction){

}


func (blockQueue BlockQueue) Len() int{
	return len(blockQueue)
}

func (blockQueue BlockQueue) Swap(i, j int){
	blockQueue[i].Nonce, blockQueue[j].Nonce =blockQueue[j].Nonce, blockQueue[i].Nonce
}

func (blockQueue BlockQueue) Less(i, j int) bool {
	return blockQueue[i].Nonce < blockQueue[j].Nonce
}