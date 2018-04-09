package blockchain

import (
	"encoding/json"

	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"github.com/EducationEKT/EKT/io/ekt8/event"
)

type BlockBody struct {
	Height       int64               `json:"height"`
	TxResults    []common.TxResult   `json:"txResults"`
	EventResults []event.EventResult `json:"eventResults"`
}

func NewBlockBody(height int64) *BlockBody {
	return &BlockBody{
		Height:       height,
		TxResults:    make([]common.TxResult, 0),
		EventResults: make([]event.EventResult, 0),
	}
}

func FromBytes(data []byte) (*BlockBody, error) {
	var body BlockBody
	err := json.Unmarshal(data, &body)
	return &body, err
}

func (body *BlockBody) AddTxResult(txResult common.TxResult) {
	body.TxResults = append(body.TxResults, txResult)
}

func (body *BlockBody) AddEventResult(eventResult event.EventResult) {
	body.EventResults = append(body.EventResults, eventResult)
}
