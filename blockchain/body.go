package blockchain

import (
	"encoding/hex"
	"encoding/json"
	"github.com/EducationEKT/EKT/userevent"
	"sync"
)

type BlockBody struct {
	Events sync.Map
}

func NewBlockBody(height int64) *BlockBody {
	return &BlockBody{
		Events: sync.Map{},
	}
}

func (body *BlockBody) AddEvent(event userevent.IUserEvent) {
	value, exist := body.Events.Load(hex.EncodeToString(event.GetFrom()))
	var list []string
	if !exist {
		list = make([]string, 1)
	} else {
		list = value.([]string)
		if list == nil {
			list = make([]string, 1)
		}
	}
	list = append(list, event.EventId())
	body.Events.Store(hex.EncodeToString(event.GetFrom()), list)
}

func FromBytes(data []byte) (*BlockBody, error) {
	var body BlockBody
	err := json.Unmarshal(data, &body)
	return &body, err
}

func (body *BlockBody) Bytes() []byte {
	data, _ := json.Marshal(body)
	return data
}
