package tx_pool

import (
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
)

const (
	Block = iota
	Queue
)

var txPool TxPool

//等待依赖队列 k:user address v:transactions of user
type BlockQueue map[string]UserTransactions

//就绪队列 k:transaction id v:transaction
type ReadyQueue map[string]*common.Transaction

//wrapper for sort
type UserTransactions []*common.Transaction

func init() {
	txPool = TxPool{}
}

type TxPool struct {
	ready ReadyQueue
	block BlockQueue
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
/*
把交易放在 txPool 里等待打包
 */
func (txPool TxPool) Park(tx *common.Transaction, reason int) {

}
/*
当交易被区块打包后,将交易移出txPool
*/
func (TxPool TxPool) Notify(tx *common.Transaction) {

}
/*当交易被区块打包后,将交易批量移出txPool

 */
func (txPool TxPool) BatchNotify(txs []*common.Transaction) {

}

/*
返回就绪队列中指定数量的交易
如果size小于等于0，返回全部
*/
func (tx TxPool)Fetch(size int)[]*common.Transaction{
	var t []*common.Transaction
	return t
}

func (u UserTransactions) Len() int {
	return len(u)
}

func (u UserTransactions) Swap(i, j int) {
	u[i].Nonce, u[j].Nonce = u[j].Nonce, u[i].Nonce
}

func (u UserTransactions) Less(i, j int) bool {
	return u[i].Nonce < u[j].Nonce
}
