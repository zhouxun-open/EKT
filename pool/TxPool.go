package pool

import (
	"encoding/hex"
	"github.com/EducationEKT/EKT/core/userevent"
)

const (
	TX_POOL_CHAN_SIZE       = 10
	GET_USEREVENT_CHAN_SIZE = 100

	Block = 0
	Ready = 1
)

type MultiFetcher struct {
	Chan chan userevent.SortedUserEvent
	num  int
}

type UserEventGetter struct {
	Chan      chan userevent.SortedUserEvent
	EventType int
	Address   string
}

type EventGetter struct {
	Chan    chan userevent.IUserEvent
	eventId string
}

func NewUserEventGetter(address string, eventType int) UserEventGetter {
	return UserEventGetter{
		Chan:      make(chan userevent.SortedUserEvent),
		EventType: eventType,
		Address:   address,
	}
}

func (pool *TxPool) GetReadyEvents(address string) userevent.SortedUserEvent {
	getter := NewUserEventGetter(address, Ready)
	pool.UserEventGetter <- getter
	return <-getter.Chan
}

func (pool *TxPool) GetBlockEvents(address string) userevent.SortedUserEvent {
	getter := NewUserEventGetter(address, Block)
	pool.UserEventGetter <- getter
	return <-getter.Chan
}

func NewEventGetter(eventId string) EventGetter {
	return EventGetter{
		Chan:    make(chan userevent.IUserEvent),
		eventId: eventId,
	}
}

func NewMultiFatcher(num int) MultiFetcher {
	return MultiFetcher{
		Chan: make(chan userevent.SortedUserEvent),
		num:  num,
	}
}

type TxPool struct {
	all   map[string]userevent.IUserEvent
	ready map[string]userevent.SortedUserEvent
	block map[string]userevent.SortedUserEvent

	SingleReady chan userevent.IUserEvent
	SingleBlock chan userevent.IUserEvent
	MultiReady  chan userevent.SortedUserEvent
	MultiBlock  chan userevent.SortedUserEvent

	SingleFetcher chan chan userevent.IUserEvent
	MultiFetcher  chan MultiFetcher

	EventGetter chan EventGetter
	EventPutter chan userevent.IUserEvent

	UserEventGetter chan UserEventGetter

	Notify chan []string
}

func NewTxPool() *TxPool {
	pool := &TxPool{
		all:   make(map[string]userevent.IUserEvent),
		ready: make(map[string]userevent.SortedUserEvent),
		block: make(map[string]userevent.SortedUserEvent),

		SingleReady: make(chan userevent.IUserEvent, TX_POOL_CHAN_SIZE),
		SingleBlock: make(chan userevent.IUserEvent, TX_POOL_CHAN_SIZE),
		MultiReady:  make(chan userevent.SortedUserEvent, TX_POOL_CHAN_SIZE),
		MultiBlock:  make(chan userevent.SortedUserEvent, TX_POOL_CHAN_SIZE),

		SingleFetcher: make(chan chan userevent.IUserEvent),
		MultiFetcher:  make(chan MultiFetcher),

		EventGetter: make(chan EventGetter, GET_USEREVENT_CHAN_SIZE),
		EventPutter: make(chan userevent.IUserEvent, GET_USEREVENT_CHAN_SIZE),

		UserEventGetter: make(chan UserEventGetter, GET_USEREVENT_CHAN_SIZE),

		Notify: make(chan []string),
	}

	go pool.loop()

	return pool
}

func (pool *TxPool) loop() {
	for {
		select {
		case event := <-pool.SingleReady:
			pool.parkReady(event)
		case event := <-pool.SingleBlock:
			pool.parkBlock(event)
		case events := <-pool.MultiReady:
			pool.parkMultiReady(events)
		case events := <-pool.MultiBlock:
			pool.parkMultiBlock(events)

		case singleer := <-pool.SingleFetcher:
			event := pool.SingleEvent()
			singleer <- event
		case multier := <-pool.MultiFetcher:
			events := pool.MultiEvent(multier.num)
			multier.Chan <- events

		case getter := <-pool.EventGetter:
			eventId := getter.eventId
			event, exist := pool.all[eventId]
			if exist {
				getter.Chan <- event
			} else {
				getter.Chan <- nil
			}
		case event := <-pool.EventPutter:
			pool.all[event.EventId()] = event

		case getter := <-pool.UserEventGetter:
			if getter.EventType == Ready {
				getter.Chan <- pool.ready[getter.Address]
			} else {
				getter.Chan <- pool.block[getter.Address]
			}

		case events := <-pool.Notify:
			pool.notify(events)
		}
	}
}

func (pool *TxPool) notify(events []string) {
	if len(events) == 0 {
		return
	}
	for _, eventId := range events {
		if event, exist := pool.all[eventId]; exist {
			readyList, exist := pool.ready[hex.EncodeToString(event.GetFrom())]
			if exist {
				pool.ready[hex.EncodeToString(event.GetFrom())] = readyList.Delete(eventId)
			}
			blockList, exist := pool.block[hex.EncodeToString(event.GetFrom())]
			if exist {
				pool.block[hex.EncodeToString(event.GetFrom())] = blockList.Delete(eventId)
			}
		}
	}
}

func (pool *TxPool) SingleEvent() userevent.IUserEvent {
	for _, list := range pool.ready {
		if len(list) > 0 {
			event := list[0]
			list = list[1:]
			pool.ready[hex.EncodeToString(event.GetFrom())] = list
			return event
		}
	}
	return nil
}

func (pool *TxPool) MultiEvent(num int) userevent.SortedUserEvent {
	for _, list := range pool.ready {
		if len(list) > 0 {
			address := hex.EncodeToString(list[0].GetFrom())
			var result userevent.SortedUserEvent
			if len(list) > num {
				result = list[:num]
				list = list[num+1:]
				pool.ready[address] = list
			} else {
				result = list
				delete(pool.ready, address)
			}
			return result
		}
	}
	return nil
}

func (pool *TxPool) parkReady(event userevent.IUserEvent) {
	address := hex.EncodeToString(event.GetFrom())
	events := pool.ready[address]
	if len(events) == 0 {
		events = make(userevent.SortedUserEvent, 0)
	}
	pool.all[event.EventId()] = event
	events = events.QuikInsert(event)
	pool.ready[address] = events
	pool.mergeReadyAndBlock(address)
}

func (pool *TxPool) parkBlock(event userevent.IUserEvent) {
	address := hex.EncodeToString(event.GetFrom())
	events := pool.block[address]
	if len(events) == 0 {
		events = make(userevent.SortedUserEvent, 0)
	}
	pool.all[event.EventId()] = event
	events = events.QuikInsert(event)
	pool.block[address] = events
	pool.mergeReadyAndBlock(address)
}

func (pool *TxPool) parkMultiReady(events userevent.SortedUserEvent) {
	if len(events) == 0 {
		return
	}
	address := hex.EncodeToString(events[0].GetFrom())
	readyList, exist := pool.ready[address]
	if !exist {
		readyList = make(userevent.SortedUserEvent, 0)
	}
	for i, _ := range events {
		pool.all[events[i].EventId()] = events[i]
		readyList.QuikInsert(events[i])
	}
	pool.mergeReadyAndBlock(address)
}

func (pool *TxPool) parkMultiBlock(events userevent.SortedUserEvent) {
	if len(events) == 0 {
		return
	}
	address := hex.EncodeToString(events[0].GetFrom())
	blockList, exist := pool.block[address]
	if !exist {
		blockList = make(userevent.SortedUserEvent, 0)
	}
	for i, _ := range events {
		pool.all[events[i].EventId()] = events[i]
		blockList.QuikInsert(events[i])
	}
	pool.mergeReadyAndBlock(address)
}

func (pool *TxPool) mergeReadyAndBlock(address string) {
	readyList, exist := pool.ready[address]
	if !exist {
		return
	}

	blockList, exist := pool.block[address]
	if !exist {
		return
	}

	numMerged := 0
	if len(readyList) > 0 && len(blockList) > 0 {
		lastNonce := readyList[readyList.Len()-1].GetNonce()
		for i, event := range blockList {
			if event.GetNonce() == lastNonce+1 {
				lastNonce++
				numMerged++
				readyList = readyList.QuikInsert(blockList[i])
			} else {
				break
			}
		}
	}

	if numMerged > 0 {
		blockList = blockList[numMerged:]
		pool.ready[address] = readyList
		pool.block[address] = blockList
	}
}
