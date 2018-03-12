package tx_pool

import "github.com/EducationEKT/EKT/io/ekt8/core/common"

var txPool TxPool

func init() {
	txPool = TxPool{}
}

type TxPool struct {
	txs map[string]*common.Transaction
}

func GetTxPool() TxPool {
	return txPool
}

func (txPool TxPool) Park(tx *common.Transaction) {
	if txPool.txs[tx.TransactionId] != nil {
		return
	}
	txPool.txs[tx.TransactionId] = tx
}
