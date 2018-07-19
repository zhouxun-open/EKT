package common

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/EducationEKT/EKT/crypto"
	"github.com/EducationEKT/EKT/db"
)

type Transactions []*Transaction

type Transaction struct {
	From         HexBytes `json:"from"`
	To           HexBytes `json:"to"`
	TimeStamp    int64    `json:"time"` // UnixTimeStamp
	Amount       int64    `json:"amount"`
	Fee          int64    `json:"fee"`
	Nonce        int64    `json:"nonce"`
	Data         string   `json:"data"`
	TokenAddress string   `json:"tokenAddress"`
	Sign         HexBytes `json:"sign"`
}

type TxResult struct {
	TxId    string `json:"txId"`
	Fee     int64  `json:"fee"`
	Success bool   `json:"success"`
	FailMsg string `json:"failMsg"`
}

func NewTransaction(from, to []byte, timestamp, amount, fee, nonce int64, data, tokenAddress string) *Transaction {
	return &Transaction{
		From:         from,
		To:           to,
		TimeStamp:    timestamp,
		Amount:       amount,
		Fee:          fee,
		Nonce:        nonce,
		Data:         data,
		TokenAddress: tokenAddress,
	}
}

func NewTransactionResult(tx Transaction, fee int64, success bool, failMessage string) *TxResult {
	return &TxResult{
		TxId:    tx.TransactionId(),
		Fee:     fee,
		Success: success,
		FailMsg: failMessage,
	}
}

func (tx Transaction) GetNonce() int64 {
	return tx.Nonce
}

func (tx *Transaction) Signature(priv []byte) error {
	sign, err := crypto.Crypto(tx.Msg(), priv)
	tx.Sign = sign
	return err
}

func (tx Transaction) GetSign() []byte {
	return tx.Sign
}

func (tx Transaction) Msg() []byte {
	return crypto.Sha3_256([]byte(tx.String()))
}

func (tx Transaction) GetFrom() []byte {
	return tx.From
}

func (tx Transaction) GetTo() []byte {
	return tx.To
}

func (tx Transaction) SetFrom(from []byte) {
	tx.From = from
}

func (tx Transaction) EventId() string {
	return tx.TransactionId()
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
	return hex.EncodeToString(crypto.Sha3_256(txData))
}

func (tx *Transaction) String() string {
	return fmt.Sprintf(`{"from": "%s", "to": "%s", "time": %d, "amount": %d, "fee": %d, "nonce": %d, "data": "%s", "tokenAddress": "%s"}`,
		hex.EncodeToString(tx.From), hex.EncodeToString(tx.To), tx.TimeStamp, tx.Amount, tx.Fee, tx.Nonce, tx.Data, tx.TokenAddress)
}

func (tx Transaction) Bytes() []byte {
	data, _ := json.Marshal(tx)
	return data
}
