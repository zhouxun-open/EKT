package tx_pool

import (
	"sort"

	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"github.com/EducationEKT/EKT/io/ekt8/event"
	"strings"
)

const (
	Block = 0
	Ready = 1
)

//等待依赖队列 k:user address v:transactions of user
type BlockTxQueue map[string]UserTransactions

//就绪队列 k:transaction id v:transaction
type ReadyTxQueue map[string]*common.Transaction

type ReadyEventQueue map[string]event.Event

type BlockEventQueue map[string]UserEvents

//wrapper for sort
type UserTransactions []*common.Transaction

type UserEvents []event.Event

type Pool struct {
	txReady    ReadyTxQueue
	txBlock    BlockTxQueue
	eventReady ReadyEventQueue
	eventBlock BlockEventQueue
}

func NewPool() *Pool {
	return &Pool{
		txReady:    make(map[string]*common.Transaction),
		txBlock:    make(map[string]UserTransactions),
		eventBlock: make(map[string]UserEvents),
		eventReady: make(map[string]event.Event),
	}
}

func (pool Pool) ParkEvent(evt event.Event, reason int) {
	if Ready == reason {
		pool.eventReady[evt.EventParam.Id()] = evt
	} else if Block == reason {
		if strings.EqualFold(event.UpdatePublicKeyEvent, evt.EventType) {
			evtParam := evt.EventParam.(event.UpdatePublicKeyParam)
			events := pool.eventBlock[evtParam.Address]
			events = append(events, evt)
			sort.Sort(events)
			pool.eventBlock[evtParam.Address] = events
		}
	}
}

func (pool Pool) NotifyEvent(evt event.Event) {
	delete(pool.eventReady, evt.EventParam.Id())
	if strings.EqualFold(evt.EventType, event.UpdatePublicKeyEvent) {
		param := evt.EventParam.(event.UpdatePublicKeyParam)
		address := param.Address
		nonce := param.Nonce
		txs := pool.txBlock[address]
		if txs != nil && len(txs) > 0 {
			for i, _tx := range txs {
				if _tx.Nonce == nonce+1 {
					txs = append(txs[:i], txs[i+1:]...)
					pool.txReady[_tx.TransactionId()] = _tx
					pool.txBlock[address] = txs
					break
				}
			}
		}
		events := pool.eventBlock[address]
		if events != nil && len(events) > 0 {
			for i, _evt := range events {
				if _evt.EventParam.(event.UpdatePublicKeyParam).Nonce == nonce+1 {
					events = append(events[:i], events[i+1:]...)
					pool.eventReady[_evt.EventParam.Id()] = _evt
					pool.eventBlock[address] = events
					break
				}
			}
		}
	}
}

/*
把交易放在 txPool 里等待打包
*/
func (txPool Pool) ParkTx(tx *common.Transaction, reason int) {
	if reason == Ready {
		txPool.txReady[tx.TransactionId()] = tx
	} else if reason == Block {
		txs_slice := txPool.txBlock[tx.From]
		txs_slice = append(txs_slice, tx)
		sort.Sort(txs_slice)
		txPool.txBlock[tx.From] = txs_slice
	}
}

/*
当交易被区块打包后,将交易移出txPool
如果当前用户有Nonce比当前大一的tx在Block队列，则移动至ready队列
*/
func (pool Pool) Notify(tx *common.Transaction) {
	delete(pool.txReady, tx.TransactionId())

	address := tx.From
	nonce := tx.Nonce

	txs := pool.txBlock[address]
	if txs != nil && len(txs) > 0 {
		for i, _tx := range txs {
			if _tx.Nonce == nonce+1 {
				txs = append(txs[:i], txs[i+1:]...)
				pool.txReady[_tx.TransactionId()] = _tx
				pool.txBlock[address] = txs
				break
			}
		}
	}

	events := pool.eventBlock[address]
	if events != nil && len(events) > 0 {
		for i, _evt := range events {
			if _evt.EventParam.(event.UpdatePublicKeyParam).Nonce == nonce+1 {
				events = append(events[:i], events[i+1:]...)
				pool.eventReady[_evt.EventParam.Id()] = _evt
				pool.eventBlock[address] = events
				break
			}
		}
	}
}

/*当交易被区块打包后,将交易批量移出txPool

 */
func (txPool Pool) BatchNotify(txs []*common.Transaction) {
	for _, tx := range txs {
		txPool.Notify(tx)
	}
}

/*
返回能够打包的指定数量的交易
如果size小于等于0，返回全部
*/
func (txPool Pool) Fetch(size int32) (result []*common.Transaction) {
	result = []*common.Transaction{}
	record := []*common.Transaction{}
	var count int32 = 0
	if size < 0 {
		size = ^(size << 31)
	}
	for {
		for txId, transaction := range txPool.txReady {
			delete(txPool.txReady, txId)
			result = append(result, transaction)
			count++
			record = append(record, transaction)
			if count >= size { //watch
				txPool.BatchNotify(record)
				return
			}
			if len(txPool.txReady) == 0 {
				txPool.BatchNotify(record)
				record = []*common.Transaction{}
			}
			if len(txPool.txReady) == 0 {
				return
			}
		}
	}
	return
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

func (events UserEvents) Len() int {
	return len(events)
}

func (events UserEvents) Swap(i, j int) {
	events[i], events[j] = events[j], events[i]
}

func (events UserEvents) Less(i, j int) bool {
	return events[i].EventParam.(event.UpdatePublicKeyParam).Nonce < events[j].EventParam.(event.UpdatePublicKeyParam).Nonce
}
