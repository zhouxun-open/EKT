package event

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/EducationEKT/EKT/crypto"
	"github.com/EducationEKT/EKT/db"
)

const (
	NewAccountEvent      = "NewAccount"
	UpdatePublicKeyEvent = "UpdatePublicKey"
)

type EventParam interface {
	EventType() string
	Validate() bool
	Id() string
}

type Event struct {
	EventType  string
	EventParam EventParam
}

type NewAccountParam struct {
	Address string
	PubKey  string
	Nonce   int64
	EventId string
}

type UpdatePublicKeyParam struct {
	Address   string
	NewPubKey string
	Nonce     int64
	EventId   string
}

type EventResult struct {
	EventId string `json:"eventId"`
	Success bool   `json:"success"`
	Reason  string `json:"reason"`
}

func (evtResult EventResult) Bytes() []byte {
	data, _ := json.Marshal(evtResult)
	return data
}

func (event Event) EventId() []byte {
	data, _ := json.Marshal(event)
	return crypto.Sha3_256(data)
}

func GetEvent(eventId []byte) *Event {
	data, err := db.GetDBInst().Get(eventId)
	if err != nil {
		return nil
	}
	return FromBytes(data)
}

func FromBytes(data []byte) *Event {
	var event Event
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil
	}
	return &event
}

func (newAccountParam NewAccountParam) EventType() string {
	return NewAccountEvent
}

func (newAccountParam NewAccountParam) Validate() bool {
	msg := []byte(fmt.Sprintf(`{"address": "%s", "pubKey": "%s", "nonce": %d}`, newAccountParam.Address, newAccountParam.PubKey, newAccountParam.Nonce))
	msg = crypto.Sha3_256(msg)
	value, err := hex.DecodeString(newAccountParam.EventId)
	if err != nil {
		return false
	}
	pubKey2, err := crypto.RecoverPubKey(msg, value)
	pubKey, err := hex.DecodeString(newAccountParam.PubKey)
	if err != nil {
		return false
	}
	if !bytes.Equal(pubKey, pubKey2) {
		return false
	}
	return true
}

func (newAccountParam NewAccountParam) Id() string {
	return newAccountParam.EventId
}

func (updatePublicKeyParam UpdatePublicKeyParam) Validate() bool {
	msg := []byte(fmt.Sprintf(`{"address": "%s", "pubKey": "%s", "nonce": %d}`, updatePublicKeyParam.Address, updatePublicKeyParam.NewPubKey, updatePublicKeyParam.Nonce))
	msg = crypto.Sha3_256(msg)
	value, err := hex.DecodeString(updatePublicKeyParam.EventId)
	if err != nil {
		return false
	}
	pubKey2, err := crypto.RecoverPubKey(msg, value)
	pubKey, err := hex.DecodeString(updatePublicKeyParam.NewPubKey)
	if err != nil {
		return false
	}
	if !bytes.Equal(pubKey, pubKey2) {
		return false
	}
	return true
}

func (updatePublicKeyParam UpdatePublicKeyParam) EventType() string {
	return UpdatePublicKeyEvent
}

func (updatePublicKeyParam UpdatePublicKeyParam) Id() string {
	return updatePublicKeyParam.EventId
}

func (event Event) ValidateEvent() bool {
	if event.EventParam.EventType() != event.EventType {
		return false
	}
	return event.EventParam.Validate()
}
