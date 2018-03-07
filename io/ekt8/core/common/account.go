package common

import (
	"bytes"

	"github.com/EducationEKT/EKT/io/ekt8/crypto"
)

type Account struct {
	Address    []byte `json:"address"`
	PublickKey []byte `json:"publicKey"`
	Amount     int64  `json:"amount"`
	Nonce      int64  `json:"nonce"`
}
