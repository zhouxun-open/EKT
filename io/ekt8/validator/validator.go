package validator

import (
	"bytes"
	"encoding/hex"

	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"github.com/EducationEKT/EKT/io/ekt8/crypto"
	"github.com/EducationEKT/EKT/io/ekt8/stat"
)

func Validate(address, msg, sign []byte) bool {
	account, err := stat.GetAccount(address)
	if err != nil {
		return false
	}
	pubkey, err := crypto.RecoverPubKey(msg, sign)
	if err != nil {
		return false
	}
	return bytes.Equal(pubkey, account.PublicKey())
}

func ValidateTx(tx *common.Transaction) bool {
	address, err := hex.DecodeString(tx.From)
	if err != nil {
		return false
	}
	account, err := stat.GetAccount(address)
	if err != nil || account == nil {
		return false
	}
	sign, err := hex.DecodeString(tx.Sign)
	if err != nil {
		return false
	}
	pubkey, err := crypto.RecoverPubKey(tx.Bytes(), sign)
	if err != nil {
		return false
	}
	return bytes.Equal(pubkey, account.PublicKey())
}
