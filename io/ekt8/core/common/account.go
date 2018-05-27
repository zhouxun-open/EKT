package common

import (
	"bytes"
	"encoding/hex"
	"encoding/json"

	"github.com/EducationEKT/EKT/io/ekt8/crypto"
)

const (
	Default_Crypto_Method = "secp256k1"
)

type Account struct {
	hexAddress    string           `json:"address"`
	hexPublickKey string           `json:"publicKey"`
	cryptoMethod  string           `json:"cryptoMethod"`
	amount        int64            `json:"amount"`
	nonce         int64            `json:"nonce"`
	balances      map[string]int64 `json:"balances"`
}

func NewAccount(address, pubKey []byte) *Account {
	return &Account{
		hexAddress:    hex.EncodeToString(address),
		hexPublickKey: hex.EncodeToString(pubKey),
		nonce:         0,
		amount:        0,
	}
}

func (account Account) ToBytes() []byte {
	data, _ := json.Marshal(account)
	return data
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
}

func (account Account) ReduceAmount(amount int64) {
	account.amount -= amount
	account.nonce++
}

func (account Account) AlterPublicKey(newPublicKey []byte) {
	account.hexPublickKey = hex.EncodeToString(newPublicKey)
	account.nonce++
}

func FromPubKeyToAddress(pubKey []byte) []byte {
	hash := crypto.Sha3_256(pubKey)
	address := crypto.Sha3_256(crypto.Sha3_256(append([]byte("EKT"), hash...)))
	return address
}

func ValidatePubKey(pubKey, address []byte) bool {
	return bytes.Equal(FromPubKeyToAddress(pubKey), address)
}
