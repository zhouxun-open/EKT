package userevent

import (
	"bytes"
	"github.com/EducationEKT/EKT/core/types"
	"github.com/EducationEKT/EKT/crypto"
	"sort"
	"strings"
)

const (
	TYPE_USEREVENT_TRANSACTION  = "transaction"
	TYPE_USEREVENT_PUBLIC_TOKEN = "issuetoken"
)

type IUserEvent interface {
	GetNonce() int64
	Msg() []byte
	GetSign() []byte
	GetFrom() []byte
	Type() string
	EventId() string
}

type SortedUserEvent []IUserEvent

func Validate(userEvent IUserEvent) bool {
	pubKey, err := crypto.RecoverPubKey(userEvent.Msg(), userEvent.GetSign())
	if err != nil {
		return false
	}
	return bytes.EqualFold(types.FromPubKeyToAddress(pubKey), userEvent.GetFrom())
}

func (events SortedUserEvent) Assert() bool {
	for i := 0; i < len(events)-1; i++ {
		if events[i].GetNonce()+1 != events[i+1].GetNonce() {
			return false
		}
	}
	return true
}

func (events SortedUserEvent) Delete(eventId string) SortedUserEvent {
	if len(events) == 0 {
		return events
	}
	index := events.Index(eventId)
	if index == -1 {
		return events
	}
	return append(events[:index], events[index+1:]...)
}

func (events SortedUserEvent) Index(eventId string) int {
	for index, event := range events {
		if strings.EqualFold(event.EventId(), eventId) {
			return index
		}
	}
	return -1
}

func (events SortedUserEvent) QuickInsert(event IUserEvent) SortedUserEvent {
	if len(events) == 0 {
		return append(events, event)
	}
	if event.GetNonce() < events[0].GetNonce() {
		list := make(SortedUserEvent, 0)
		list = append(list, event)
		list = append(list, events...)
		return list
	}
	if event.GetNonce() > events[len(events)-1].GetNonce() {
		return append(events, event)
	}
	for i := 0; i < len(events)-1; i++ {
		if events[i].GetNonce() < event.GetNonce() && event.GetNonce() < events[i+1].GetNonce() {
			list := make(SortedUserEvent, 0)
			list = append(list, events[:i+1]...)
			list = append(list, event)
			list = append(list, events[i+1:]...)
			return list
		}
	}
	return events.quickSort(event)
}

func (events SortedUserEvent) quickSort(event IUserEvent) SortedUserEvent {
	events = append(events, event)
	sort.Sort(events)
	return events
}

func (events SortedUserEvent) Len() int {
	return len(events)
}

func (events SortedUserEvent) Less(i, j int) bool {
	return events[i].GetNonce() < events[j].GetNonce()
}

func (events SortedUserEvent) Swap(i, j int) {
	events[i], events[j] = events[j], events[i]
}
