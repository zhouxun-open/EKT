package crypto

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	pub, priv := GenerateKeyPair()
	fmt.Printf("pubKey=%s, secKey=%s\n", hex.EncodeToString(pub), hex.EncodeToString(priv))
	data := Sha3_256([]byte("hello world"))
	sign, err := Crypto(data, priv)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if pubKey, err := RecoverPubKey(data, sign); err == nil && bytes.Equal(pubKey, pub) {
		fmt.Printf(`recovered public key: %s`, hex.EncodeToString(pubKey))
	} else {
		t.Fail()
	}
	if !Verify(sign, pub, data) {
		fmt.Println("verify fail")
		t.Fail()
	}
}

func BenchmarkCrypto(b *testing.B) {
	_, priv := GenerateKeyPair()
	for i := 0; i < b.N; i++ {
		data := Sha3_256([]byte("hello world"))
		Crypto(data, priv)
	}
}

func BenchmarkRecoverPubKey(b *testing.B) {
	_, priv := GenerateKeyPair()
	data := Sha3_256([]byte("hello world"))
	sign, err := Crypto(data, priv)
	if err != nil {
		b.Fail()
	}
	for i := 0; i < b.N; i++ {
		RecoverPubKey(data, sign)
	}
}

func BenchmarkVerify(b *testing.B) {
	pub, priv := GenerateKeyPair()
	data := Sha3_256([]byte("hello world"))
	sign, err := Crypto(data, priv)
	if err != nil {
		b.Fail()
	}
	for i := 0; i < b.N; i++ {
		Verify(sign, pub, data)
	}
}
