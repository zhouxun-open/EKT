package crypto

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/randentropy"
)

func TestGenerateKeyPair(t *testing.T) {
	pub, priv := GenerateKeyPair()
	fmt.Printf("pubKey=%s, secKey=%s\n", hex.EncodeToString(pub), hex.EncodeToString(priv))
	data := randentropy.GetEntropyCSPRNG(32)
	sign, err := Crypto(data, priv)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if !Verify(sign, pub, data) {
		fmt.Println("verify fail")
		t.Fail()
	}
}
