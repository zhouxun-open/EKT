package common

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strings"
)

type Transactions []*Transaction

type Transaction struct {
	From      NormalAddress
	To        NormalAddress
	TimeStamp Time // UnixTimeStamp
	Amount    int64
	Nonce     int64
	R, S, V   *big.Int
	Sign      string
}

func (transactions Transactions) Len() int {
	return len(transactions)
}

func (transactions Transactions) Less(i, j int) bool {
	return strings.Compare(transactions[i].Sign, transactions[j].Sign) > 0
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
	return fmt.Sprintf(`{"from": "%s", "to": "%s": "%s", "time": %d, "amount": %d, "nonce": %d}`,
		hex.EncodeToString(tx.From[:]),
		hex.EncodeToString(tx.To[:]),
		tx.TimeStamp, tx.Amount, tx.Nonce)
}

func (tx *Transaction) ValidateSignature() error {
	// TODO
	return nil
}
