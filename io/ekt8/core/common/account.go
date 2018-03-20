package common

import (
	"encoding/hex"
	"fmt"
)

type Account struct {
	HexAddress    string `json:"HexAddress"`
	HexPublickKey string `json:"publicKey"`
	Amount        int64  `json:"Amount"`
	Nonce         int64  `json:"Nonce"`
}

func (account Account) ToString() string {
	return fmt.Sprintf(`{"HexAddress": "%s", "publicKey": "%s", "Amount": %s, "Nonce": %d}`,
		account.HexAddress, account.HexPublickKey, account.Amount, account.Nonce)
}

func (account Account) ToBytes() []byte {
	return []byte(account.ToString())
}

func (account Account) GetNonce() int64 {
	return account.Nonce
}

func (account Account) Address() []byte {
	address, _ := hex.DecodeString(account.HexAddress)
	return address
}

func (account Account) PublicKey() []byte {
	publicKey, _ := hex.DecodeString(account.HexPublickKey)
	return publicKey
}

func (account Account) GetAmount() int64 {
	return account.Amount
}

func (account Account) AddAmount(amount int64) {
	account.Amount += amount
	//account.Nonce++
}

func (account Account) ReduceAmount(amount int64) {
	account.Amount -= amount
	account.Nonce++
}

func (account Account) AlterPublicKey(newPublicKey []byte) {
	account.HexPublickKey = hex.EncodeToString(newPublicKey)
	account.Nonce++
}
