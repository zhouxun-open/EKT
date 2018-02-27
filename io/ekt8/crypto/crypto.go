package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

var theCurve = new(secp256k1.BitCurve)

func init() {
	theCurve.P, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)
	theCurve.N, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)
	theCurve.B, _ = new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000007", 16)
	theCurve.Gx, _ = new(big.Int).SetString("79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798", 16)
	theCurve.Gy, _ = new(big.Int).SetString("483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8", 16)
	theCurve.BitSize = 256
}

func Crypto(data []byte, secKey []byte) ([]byte, error) {
	return secp256k1.Sign(data, secKey)
}

func Verify(sign, pubKey, msg []byte) bool {
	pub2, err:=secp256k1.RecoverPubkey(msg, sign)
	if err!=nil || !bytes.Equal(pubKey, pub2){
		return false
	}
	return true
}

func GenerateKeyPair() (pubkey, privkey []byte) {
	key, err := ecdsa.GenerateKey(S256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	pubkey = elliptic.Marshal(S256(), key.X, key.Y)
	return pubkey, math.PaddedBigBytes(key.D, 32)
}

func S256() *secp256k1.BitCurve {
	return theCurve
}
