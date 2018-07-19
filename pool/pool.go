package pool

import (
	"encoding/hex"
	"github.com/EducationEKT/EKT/userevent"
	"strings"
	"sync"
)

const (
	Block = 0
	Ready = 1
)

type Pool struct {
	ready sync.Map
	block sync.Map
}

func NewPool() *Pool {
	return &Pool{
		ready: sync.Map{},
		block: sync.Map{},
	}
}

// 根据address获取该地址pending/queue的交易信息
func (pool Pool) GetReadyEvents(address string) userevent.SortedUserEvent {
	events, exist := pool.ready.Load(address)
	if exist {
		return events.(userevent.SortedUserEvent)
	}
	return nil
}

/*
 * 把交易放在 pool 里等待打包
 */
func (pool *Pool) Park(event userevent.IUserEvent, reason int) bool {
	if reason == Ready {
		return pool.parkReady(event)
	}

	return pool.parkBlock(event)
}

func (pool *Pool) parkReady(event userevent.IUserEvent) bool {
	readyEvents, exist := pool.ready.Load(hex.EncodeToString(event.GetFrom()))

	var list userevent.SortedUserEvent

	if exist {
		list = readyEvents.(userevent.SortedUserEvent)
		if len(list) == 0 {
			list = make(userevent.SortedUserEvent, 0)
		} else if list[list.Len()-1].GetNonce()+1 != event.GetNonce() {
			return false
		}
	} else {
		list = make(userevent.SortedUserEvent, 0)
	}
	list = append(list, event)

	pool.ready.Store(hex.EncodeToString(event.GetFrom()), list)

	pool.MergeReadyAndBlock(event.GetFrom())

	return true
}

func (pool *Pool) parkBlock(event userevent.IUserEvent) bool {
	blockEvents, exist := pool.block.Load(hex.EncodeToString(event.GetFrom()))
	if exist {
		list := blockEvents.(userevent.SortedUserEvent)
		if list == nil {
			list = make(userevent.SortedUserEvent, 0)
		}
		list = append(list, event)
		pool.block.Store(hex.EncodeToString(event.GetFrom()), list)
	} else {
		list := make(userevent.SortedUserEvent, 0)
		list = append(list, event)
		pool.block.Store(hex.EncodeToString(event.GetFrom()), list)
	}
	pool.MergeReadyAndBlock(event.GetFrom())
	return true
}

func (pool *Pool) MergeReadyAndBlock(from []byte) {
	address := hex.EncodeToString(from)
	readyEvents, exist := pool.ready.Load(address)
	if !exist {
		return
	}
	blockEvents, exist := pool.block.Load(address)
	if !exist {
		return
	}
	readyList := readyEvents.(userevent.SortedUserEvent)
	blockList := blockEvents.(userevent.SortedUserEvent)
	if len(readyList) > 0 && len(blockList) > 0 {
		lastNonce := readyList[readyList.Len()-1].GetNonce()
		numMerged := 0
		for i, event := range blockList {
			if event.GetNonce() == lastNonce+1 {
				lastNonce++
				readyList = append(readyList, blockList[i])
				numMerged++
			}
		}
		blockList = blockList[numMerged:]
		pool.ready.Store(address, readyList)
		pool.block.Store(address, blockList)
	}
}

/*
当交易被区块打包后,将交易移出pool
*/
func (pool *Pool) Notify(from, eventId string) {
	// notify ready queue
	readyList, exist := pool.ready.Load(from)
	if exist {
		list := readyList.(userevent.SortedUserEvent)
		newList := make(userevent.SortedUserEvent, len(list))
		for _, event := range list {
			if !strings.EqualFold(event.EventId(), eventId) {
				newList = append(newList, event)
			}
		}
		pool.ready.Store(from, newList)
	}

	// notify block
	blockList, exist := pool.block.Load(from)
	if exist {
		list := blockList.(userevent.SortedUserEvent)
		newList := make(userevent.SortedUserEvent, len(list))
		for _, event := range list {
			if !strings.EqualFold(event.EventId(), eventId) {
				newList = append(newList, event)
			}
		}
		pool.block.Store(from, newList)
	}
}

func (pool *Pool) Fetch() userevent.SortedUserEvent {
	var events userevent.SortedUserEvent = nil
	var from string
	pool.ready.Range(func(key, value interface{}) bool {
		list, ok := value.(userevent.SortedUserEvent)
		if !ok {
			return true
		}
		if len(list) > 0 {
			from = key.(string)
			events = list
			return false
		}
		return true
	})
	pool.ready.Delete(from)
	return events
}
