package common

import (
	"encoding/hex"
	"fmt"
)

type Account struct {
	address    string `json:"address"`
	publickKey string `json:"publicKey"`
	amount     int64  `json:"amount"`
	nonce      int64  `json:"nonce"`
}

func (account Account) ToString() string {
	return fmt.Sprintf(`{"address": "%s", "publicKey": "%s", "amount": %s, "nonce": %d}`,
		account.address, account.publickKey, account.amount, account.nonce)
}

func (account Account) ToBytes() []byte {
	return []byte(account.ToString())
}

func (account Account) Nonce() int64 {
	return account.nonce
}

func (account Account) Address() []byte {
	address, _ := hex.DecodeString(account.address)
	return address
}

func (account Account) PublicKey() []byte {
	publicKey, _ := hex.DecodeString(account.publickKey)
	return publicKey
}

func (account Account) Amount() int64 {
	return account.amount
}

func (account Account) AddAmount(amount int64) {
	account.amount += amount
	//account.nonce++
}

func (account Account) ReduceAmount(amount int64) {
	account.amount -= amount
	account.nonce++
}

func (account Account) AlterPublicKey(newPublicKey []byte) {
	account.publickKey = hex.EncodeToString(newPublicKey)
	account.nonce++
}
