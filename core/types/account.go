package types

import (
	"bytes"
	"encoding/json"

	"github.com/EducationEKT/EKT/crypto"
)

const (
	Default_Crypto_Method = "secp256k1"
)

type Account struct {
	Address  HexBytes         `json:"address"`
	Amount   int64            `json:"amount"`
	Nonce    int64            `json:"nonce"`
	Balances map[string]int64 `json:"balances"`
}

func CreateAccount(address []byte, Amount int64) Account {
	return Account{
		Address: address,
		Amount:  Amount,
		Nonce:   0,
	}
}

func NewAccount(address []byte) *Account {
	return &Account{
		Address: address,
		Nonce:   0,
		Amount:  0,
	}
}

func (account Account) ToBytes() []byte {
	data, _ := json.Marshal(account)
	return data
}

func (account Account) GetNonce() int64 {
	return account.Nonce
}

func (account Account) GetAmount() int64 {
	return account.Amount
}

func (account *Account) AddAmount(amount int64) {
	account.Amount = account.Amount + amount
}

func (account *Account) ReduceAmount(amount int64) {
	account.Amount = account.Amount - amount
	account.Nonce++
}

func FromPubKeyToAddress(pubKey []byte) []byte {
	hash := crypto.Sha3_256(pubKey)
	address := crypto.Sha3_256(crypto.Sha3_256(append([]byte("EKT"), hash...)))
	return address
}

func ValidatePubKey(pubKey, address []byte) bool {
	return bytes.Equal(FromPubKeyToAddress(pubKey), address)
}
