package common

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/EducationEKT/EKT/io/ekt8/crypto"
)

type Transactions []*Transaction

type Transaction struct {
	TransactionId string
	From          string
	To            string
	TimeStamp     Time // UnixTimeStamp
	Amount        int64
	Nonce         int64
	Sign          string
}

type TxResult struct {
	TxId      string `json:"txId"`
	From      string `json:"from"`
	To        string `json:"to"`
	Amount    int64  `json:"Amount"`
	TimeStamp Time   `json:"timestamp"`
	Nonce     int64  `json:"Nonce"`
	Fee       int64  `json:"fee"`
	Sign      string `json:"sign"`
	Success   bool   `json:"success"`
	FailMsg   string `json:"failMsg"`
}

func NewTransactionResult(tx *Transaction, success bool, failMessage string) *TxResult {
	return &TxResult{
		TxId:      tx.TransactionId,
		From:      tx.From,
		To:        tx.To,
		TimeStamp: tx.TimeStamp,
		Amount:    tx.Amount,
		Nonce:     tx.Nonce,
		Sign:      tx.Sign,
		Fee:       1e6,
		Success:   success,
		FailMsg:   failMessage,
	}
}

func (txResult *TxResult) ToTransaction() *Transaction {
	return &Transaction{
		TransactionId: txResult.TxId,
		From:          txResult.From,
		To:            txResult.To,
		TimeStamp:     txResult.TimeStamp,
		Amount:        txResult.Amount,
		Nonce:         txResult.Nonce,
		Sign:          txResult.Sign,
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
	return strings.Compare(transactions[i].TransactionId, transactions[j].TransactionId) > 0
}

func (transactions Transactions) Swap(i, j int) {
	transactions[i], transactions[j] = transactions[j], transactions[i]
}

func (transactions Transactions) Hash() string {
	sort.Sort(transactions)
	bytes, _ := json.Marshal(transactions)
	fmt.Println(string(bytes))
	return ""
}

func (tx *Transaction) Bytes() []byte {
	return []byte(tx.String())
}

func (tx *Transaction) String() string {
	return fmt.Sprintf(`{"from": "%s", "to": "%s", "time": %d, "Amount": %d, "Nonce": %d}`,
		tx.From, tx.To,
		tx.TimeStamp, tx.Amount, tx.Nonce)
}

func (tx *Transaction) Validate() error {
	sign, err := hex.DecodeString(tx.Sign)
	if err != nil {
		return err
	}
	signedTxId := crypto.Sha3_256(sign)
	txIdBytes, err := hex.DecodeString(tx.TransactionId)
	if err != nil {
		return err
	}
	if !bytes.Equal(signedTxId, txIdBytes) {
		return errors.New("Invalid transaction")
	}
	var account Account
	//HexAddress, err := hex.DecodeString(tx.From)
	if err != nil {
		return err
	}
	if !crypto.Verify(sign, account.PublicKey(), tx.Bytes()) {
		return errors.New("Invalid Signature")
	}
	return nil
}
