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

var TestPubKey = []byte("0419e18e6ba7feb6e76fbd933030a20dfe59a24e04ec51133ef5b6825d0457524f1044f00826cc01bbfa981f50e6cbe16446a0760b3d872d716a51e409d8a5deb7")
var TestSecKey = []byte("ebd71b84d374f881e5280b708c502b2e69de22578d2b1c0fd4dcccd29ad7f4f3")
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
	pub2, err := secp256k1.RecoverPubkey(msg, sign)
	if err != nil || !bytes.Equal(pubKey, pub2) {
		return false
	}
	return true
}

func RecoverPubKey(msg, sign []byte) ([]byte, error) {
	return secp256k1.RecoverPubkey(msg, sign)
}

func PubKey(priv []byte) ([]byte, error) {
	data := Sha3_256([]byte("helloworld"))
	sign, err := Crypto(data, priv)
	if err != nil {
		return nil, err
	}
	return RecoverPubKey(data, sign)
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
