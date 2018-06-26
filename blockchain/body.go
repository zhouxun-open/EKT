package blockchain

import (
	"encoding/json"

	"encoding/hex"
	"github.com/EducationEKT/EKT/core/common"
	"github.com/EducationEKT/EKT/ctxlog"
	"github.com/EducationEKT/EKT/event"
	"github.com/EducationEKT/EKT/pool"
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

func (body *BlockBody) Bytes() []byte {
	data, _ := json.Marshal(body)
	return data
}

func (body *BlockBody) AddTxResult(txResult common.TxResult) {
	body.TxResults = append(body.TxResults, txResult)
}

func (body *BlockBody) AddEventResult(eventResult event.EventResult) {
	body.EventResults = append(body.EventResults, eventResult)
}

func (body *BlockBody) Size() int {
	return len(body.EventResults) + len(body.TxResults)
}
func (blockchain BlockChain) NewTransaction(log *ctxlog.ContextLog, tx *common.Transaction) bool {
	from, _ := hex.DecodeString(tx.From)
	log.Log("from", tx.From)
	if account, err := blockchain.GetLastBlock().GetAccount(log, from); err == nil && account != nil {
		log.Log("account", account)
		status := pool.Block
		if account.Nonce+1 == tx.Nonce {
			status = pool.Ready
		}
		log.Log("txStatus", status)
		blockchain.Pool.ParkTx(tx, status)
		log.Log("parked", true)
		return true
	}
	return false
}
