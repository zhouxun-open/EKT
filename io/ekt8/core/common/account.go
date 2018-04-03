package common

import (
	"encoding/hex"
	"fmt"
)

type Account struct {
	hexAddress    string `json:"address"`
	hexPublickKey string `json:"publicKey"`
	amount        int64  `json:"amount"`
	nonce         int64  `json:"nonce"`
}

func NewAccount(address, pubKey []byte) *Account {
	return &Account{
		hexAddress:    hex.EncodeToString(address),
		hexPublickKey: hex.EncodeToString(pubKey),
		nonce:         0,
		amount:        0,
	}
}

func (account Account) ToString() string {
	return fmt.Sprintf(`{"HexAddress": "%s", "publicKey": "%s", "Amount": %s, "Nonce": %d}`,
		account.hexAddress, account.hexPublickKey, account.amount, account.nonce)
}

func (account Account) ToBytes() []byte {
	return []byte(account.ToString())
}

func (account Account) GetNonce() int64 {
	return account.nonce
}

func (account Account) Address() []byte {
	address, _ := hex.DecodeString(account.hexAddress)
	return address
}

func (account Account) PublicKey() []byte {
	publicKey, _ := hex.DecodeString(account.hexPublickKey)
	return publicKey
}

func (account Account) GetAmount() int64 {
	return account.amount
}

func (account Account) AddAmount(amount int64) {
	account.amount += amount
	//account.Nonce++
}

func (account Account) ReduceAmount(amount int64) {
	account.amount -= amount
	account.nonce++
}

func (account Account) AlterPublicKey(newPublicKey []byte) {
	account.hexPublickKey = hex.EncodeToString(newPublicKey)
	account.nonce++
}
