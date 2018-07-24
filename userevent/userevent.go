package userevent

import (
	"bytes"
	"github.com/EducationEKT/EKT/core/common"
	"github.com/EducationEKT/EKT/crypto"
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

func (events SortedUserEvent) Len() int {
	return len(events)
}

func (events SortedUserEvent) Less(i, j int) bool {
	return events[i].GetNonce() < events[j].GetNonce()
}

func (events SortedUserEvent) Swap(i, j int) {
	events[i], events[j] = events[j], events[i]
}
