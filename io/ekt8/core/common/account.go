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
	HexAddress    string           `json:"address"`
	HexPublickKey string           `json:"-"`
	CryptoMethod  string           `json:"-"`
	Amount        int64            `json:"amount"`
	Nonce         int64            `json:"nonce"`
	Balances      map[string]int64 `json:"balances"`
}

func CreateAccount(address string, Amount int64) Account {
	return Account{
		HexAddress: address,
		Amount:     Amount,
		Nonce:      0,
	}
}

func NewAccount(address, pubKey []byte) *Account {
	return &Account{
		HexAddress:    hex.EncodeToString(address),
		HexPublickKey: hex.EncodeToString(pubKey),
		Nonce:         0,
		Amount:        0,
	}
}

func (account Account) ToBytes() []byte {
	data, _ := json.Marshal(account)
	return data
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
}

func (account Account) ReduceAmount(amount int64) {
	account.Amount -= amount
	account.Nonce++
}

func (account Account) AlterPublicKey(newPublicKey []byte) {
	account.HexPublickKey = hex.EncodeToString(newPublicKey)
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
