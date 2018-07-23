package pool

import (
	"encoding/hex"
	"github.com/EducationEKT/EKT/userevent"
	"sort"
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
	address := hex.EncodeToString(event.GetFrom())
	readyEvents, exist := pool.ready.Load(address)

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

	blockEvents, exist := pool.block.Load(address)
	var block userevent.SortedUserEvent
	if exist {
		block = blockEvents.(userevent.SortedUserEvent)
	}

	pool.MergeReadyAndBlock(event.GetFrom(), list, block)

	return true
}

func (pool *Pool) parkBlock(event userevent.IUserEvent) bool {
	address := hex.EncodeToString(event.GetFrom())
	blockEvents, exist := pool.block.Load(address)
	var list userevent.SortedUserEvent
	if exist {
		list = blockEvents.(userevent.SortedUserEvent)
		if list == nil {
			list = make(userevent.SortedUserEvent, 0)
		}
		list = append(list, event)
	} else {
		list = make(userevent.SortedUserEvent, 0)
		list = append(list, event)
	}
	sort.Sort(list)

	readyEvents, exist := pool.ready.Load(address)
	var ready userevent.SortedUserEvent
	if exist {
		ready = readyEvents.(userevent.SortedUserEvent)
	}

	pool.MergeReadyAndBlock(event.GetFrom(), ready, list)

	return true
}

func (pool *Pool) MergeReadyAndBlock(from []byte, ready userevent.SortedUserEvent, block userevent.SortedUserEvent) {
	readyList, blockList := ready, block
	numMerged := 0
	if len(ready) > 0 && len(block) > 0 {
		lastNonce := readyList[readyList.Len()-1].GetNonce()
		for i, event := range blockList {
			if event.GetNonce() == lastNonce+1 {
				lastNonce++
				readyList = append(readyList, blockList[i])
				numMerged++
			}
		}
		blockList = blockList[numMerged:]
	}

	address := hex.EncodeToString(from)

	pool.ready.Store(address, readyList)
	pool.block.Store(address, blockList)
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
			from, ok = key.(string)
			if ok {
				events = list
				return false
			}
		}
		return true
	})
	pool.ready.Delete(from)
	return events
}
