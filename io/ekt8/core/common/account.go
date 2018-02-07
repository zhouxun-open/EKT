package common

import (
	"bytes"

	"github.com/EducationEKT/EKT/io/ekt8/crypto"
)

type Account struct {
	Address    Address
	PublickKey []byte
}
