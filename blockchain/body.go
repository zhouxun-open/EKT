package blockchain

import (
	"encoding/json"
	"github.com/EducationEKT/EKT/core/userevent"
)

type BlockBody struct {
	Events []string
}

func NewBlockBody() *BlockBody {
	return &BlockBody{
		Events: make([]string, 0),
	}
}

func (body *BlockBody) AddEvent(event userevent.IUserEvent) {
	body.Events = append(body.Events, event.EventId())
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
