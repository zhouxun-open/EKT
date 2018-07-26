package userevent

import (
	"bytes"
	"github.com/EducationEKT/EKT/core/common"
	"github.com/EducationEKT/EKT/crypto"
	"strings"
)

type IUserEvent interface {
	GetNonce() int64
	Msg() []byte
	GetSign() []byte
	GetFrom() []byte
	EventId() string
}

type SortedUserEvent []IUserEvent

func Validate(userEvent IUserEvent) bool {
	pubKey, err := crypto.RecoverPubKey(userEvent.Msg(), userEvent.GetSign())
	if err != nil {
		return false
	}
	return bytes.EqualFold(common.FromPubKeyToAddress(pubKey), userEvent.GetFrom())
}

func (events SortedUserEvent) Index(eventId string) int {
	for index, event := range events {
		if strings.EqualFold(event.EventId(), eventId) {
			return index
		}
	}
	return -1
}

func (events SortedUserEvent) QuikInsert(event IUserEvent) SortedUserEvent {
	low, high := 0, len(events)
	for high > low+1 {
		m := (low + high) / 2
		if events[m].GetNonce() > event.GetNonce() {
			high = m
		} else {
			low = m
		}
	}

	left := events[:high]
	right := events[high:]
	list := append(left, event)
	return append(list, right...)
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
