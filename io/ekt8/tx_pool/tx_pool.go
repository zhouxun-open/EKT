package tx_pool

import (
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
)

const (
	Block = 0
	Ready = 1
)

//等待依赖队列 k:user address v:transactions of user
type BlockQueue map[string]UserTransactions

//就绪队列 k:transaction id v:transaction
type ReadyQueue map[string]*common.Transaction

//wrapper for sort
type UserTransactions []*common.Transaction

type TxPool struct {
	ready ReadyQueue
	block BlockQueue
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
	if reason == Ready {
		txPool.ready[tx.TransactionId] = tx
	} else if reason == Block {
		txs_slice := txPool.block[tx.From]
		txPool.block[tx.From] = append(txs_slice, tx)
	}
}

/*
当交易被区块打包后,将交易移出txPool
如果当前用户有Nonce比当前大一的tx在Block队列，则移动至ready队列
*/
func (txPool TxPool) Notify(tx *common.Transaction) {
	delete(txPool.ready, tx.TransactionId)
	txs := txPool.block[tx.From]
	if txs != nil {
		for i, _tx := range txs {
			if _tx.Nonce == tx.Nonce+1 {
				txs = append(txs[:i], txs[i+1:]...)
				txPool.ready[_tx.TransactionId] = _tx
				break
			}
		}
	}
}

/*当交易被区块打包后,将交易批量移出txPool

 */
func (txPool TxPool) BatchNotify(txs []*common.Transaction) {
	for _, tx := range txs {
		txPool.Notify(tx)
	}
}

/*
返回就绪队列中指定数量的交易
如果size小于等于0，返回全部
*/
func (tx TxPool) Fetch(size int) map[string]*common.Transaction {
	if size <= 0 {
		return tx.ready
	} else if size > len(tx.ready) {
		return tx.ready
	} else {
		returnMap := make(map[string]*common.Transaction)
		count := 0
		for k, v := range tx.ready {
			if count >= size {
				break
			}
			count++
			returnMap[k] = v
			//delete(txPool.ready, k) //delete (k,v) from readyqueue
		}
		return returnMap
	}
}

func (u UserTransactions) Len() int {
	return len(u)
}

func (u UserTransactions) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (u UserTransactions) Less(i, j int) bool {
	return u[i].Nonce < u[j].Nonce
}
