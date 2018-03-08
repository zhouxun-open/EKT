package common

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/EducationEKT/EKT/io/ekt8/blockchain_manager"
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
	return fmt.Sprintf(`{"from": "%s", "to": "%s", "time": %d, "amount": %d, "nonce": %d}`,
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
	if !bytes.Equal(signedTxId, txIdBytes) {
		return errors.New("Invalid transaction")
	}
	block, err := blockchain_manager.MainBlockChain.CurrentBlock()
	if err != nil {
		return err
	}
	var account Account
	address, err := hex.DecodeString(tx.From)
	if err != nil {
		return err
	}
	block.StatTree.GetInterfaceValue(address, &account)
	if !crypto.Verify(sign, account.PublicKey(), tx.Bytes()) {
		return errors.New("Invalid Signature")
	}
	return nil
}
