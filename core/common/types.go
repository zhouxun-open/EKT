package common

import (
	"encoding/hex"
	"fmt"
)

type HexBytes []byte

func (bytes HexBytes) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, hex.EncodeToString(bytes))), nil
}

type Object interface{}

type CoinType int

type Time int64
