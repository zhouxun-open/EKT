package blockchain

import (
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"github.com/EducationEKT/EKT/io/ekt8/event"
)

type BlockBody struct {
	TxResults   []common.TxResult
	EventResult []event.EventResult
}

func NewBlockBody() *BlockBody {
	return &BlockBody{
		TxResults:   make([]common.TxResult, 0),
		EventResult: make([]event.EventResult, 0),
	}
}

func (body *BlockBody) AddTxResult(txResult common.TxResult) {
	body.TxResults = append(body.TxResults, txResult)
}

func (body *BlockBody) AddEventResult(eventResult event.EventResult) {
	body.EventResult = append(body.EventResult, eventResult)
}
