package types

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto/sha3"
)

const AddressLength = 32

const BigAddressLength = 64

type CommonAddress struct {
	Type          int
	AddressLength int
	Address       []byte
}

type NormalAddress [AddressLength]byte

func BytesToAddress(b []byte) NormalAddress {
	var a NormalAddress
	a.SetBytes(b)
	return a
}

func (a *NormalAddress) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

func (a *NormalAddress) SetString(s string) {
	a.SetBytes([]byte(s))
}

func (a NormalAddress) String() string {
	return a.Hex()
}

func (a NormalAddress) Hex() string {
	unchecksummed := hex.EncodeToString(a[:])
	sha := sha3.NewKeccak256()
	sha.Write([]byte(unchecksummed))
	hash := sha.Sum(nil)

	result := []byte(unchecksummed)
	for i := 0; i < len(result); i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if result[i] > '9' && hashByte > 7 {
			result[i] -= 32
		}
	}
	return "0x" + string(result)
}
