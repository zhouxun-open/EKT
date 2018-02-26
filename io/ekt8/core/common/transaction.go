package common

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type Transactions []*Transaction

type Transaction struct {
	From      NormalAddress
	To        NormalAddress
	TimeStamp Time // UnixTimeStamp
	Amount    int64
	Type      CoinType
	Nonce     int64
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
