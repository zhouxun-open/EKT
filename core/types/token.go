package types

import (
	"encoding/json"
	"github.com/EducationEKT/EKT/crypto"
)

type Token struct {
	Name     string `json:"name"`
	Total    int64  `json:"total"`
	Decimals int64  `json:"decimals"`
}

func (token Token) Address() []byte {
	v, err := json.Marshal(token)
	if err != nil {
		return nil
	}
	return crypto.Sha3_256(v)
}
