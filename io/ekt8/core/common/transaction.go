package common

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/EducationEKT/EKT/io/ekt8/crypto"
	"github.com/EducationEKT/EKT/io/ekt8/db"
)

type Transactions []*Transaction

type Transaction struct {
	From         string `json:"from"`
	To           string `json:"to"`
	TimeStamp    int64  `json:"time"` // UnixTimeStamp
	Amount       int64  `json:"Amount"`
	Nonce        int64  `json:"Nonce"`
	Data         string `json:"data"`
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

func GetTransaction(txId []byte) *Transaction {
	txData, err := db.GetDBInst().Get(txId)
	if err != nil {
		return nil
	}
	return FromBytes(txData)
}

func FromBytes(data []byte) *Transaction {
	var tx Transaction
	err := json.Unmarshal(data, &tx)
	if err != nil {
		return nil
	}
	return &tx
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

func (tx *Transaction) TransactionId() string {
	txData, _ := json.Marshal(tx)
	return string(crypto.Sha3_256(txData))
}

func (tx *Transaction) String() string {
	return fmt.Sprintf(`{"from": "%s", "to": "%s", "time": %d, "amount": %d, "nonce": %d, "data": "%s", "tokenAddress": "%s"}`,
		tx.From, tx.To, tx.TimeStamp, tx.Amount, tx.Nonce, tx.Data, tx.TokenAddress)
}

func (tx Transaction) Bytes() []byte {
	data, _ := json.Marshal(tx)
	return data
}

func (tx *Transaction) Validate() bool {
	sign, err := hex.DecodeString(tx.Sign)
	if err != nil {
		return false
	}
	data := crypto.Sha3_256([]byte(tx.String()))
	if pubKey, err := crypto.RecoverPubKey(data, sign); err != nil {
		return false
	} else {
		address, err := hex.DecodeString(tx.From)
		if err != nil || !ValidatePubKey(pubKey, address) {
			return false
		}
	}
	return true
}
