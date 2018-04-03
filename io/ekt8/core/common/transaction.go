package common

import (
	"encoding/json"
	"fmt"
	"github.com/EducationEKT/EKT/io/ekt8/crypto"
	"strings"
)

type Transactions []*Transaction

type Transaction struct {
	From         string `json:"from"`
	To           string `json:"to"`
	TimeStamp    Time   `json:"time"` // UnixTimeStamp
	Amount       int64  `json:"amount"`
	Nonce        int64  `json:"nonce"`
	TokenAddress string `json:"tokenAddress"`
	Sign         string `json:"sign"`
}

type TxResult struct {
	TxId    string `json:"txId"`
	Fee     int64  `json:"fee"`
	Success bool   `json:"success"`
	FailMsg string `json:"failMsg"`
}

func NewTransactionResult(tx *Transaction, fee int64, success bool, failMessage string) *TxResult {
	return &TxResult{
		TxId:    tx.TransactionId(),
		Fee:     fee,
		Success: success,
		FailMsg: failMessage,
	}
}

func (txResult *TxResult) ToBytes() []byte {
	data, _ := json.Marshal(txResult)
	return data
}

func (txResult *TxResult) TxResult() (bool, string) {
	return txResult.Success, txResult.FailMsg
}

func (transactions Transactions) Len() int {
	return len(transactions)
}

func (transactions Transactions) Less(i, j int) bool {
	return strings.Compare(transactions[i].TransactionId(), transactions[j].TransactionId()) > 0
}

func (transactions Transactions) Swap(i, j int) {
	transactions[i], transactions[j] = transactions[j], transactions[i]
}

func (tx *Transaction) Bytes() []byte {
	return []byte(tx.String())
}

func (tx *Transaction) String() string {
	return fmt.Sprintf(`{"from": "%s", "to": "%s", "time": %d, "Amount": %d, "Nonce": %d}`,
		tx.From, tx.To,
		tx.TimeStamp, tx.Amount, tx.Nonce)
}

func (tx *Transaction) TransactionId() (ID string) {
	txData, _ := json.Marshal(tx)
	ID = string(crypto.Sha3_256(txData))
	return
}
