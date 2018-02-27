package crypto

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/randentropy"
)

func TestGenerateKeyPair(t *testing.T) {
	pub, priv:=GenerateKeyPair()
	data:=randentropy.GetEntropyCSPRNG(32)
	sign, err:=Crypto(data, priv)
	if err!= nil {
		fmt.Println(err)
		t.Fail()
	}
	if !Verify(sign, pub, data) {
		fmt.Println("verify fail")
		t.Fail()
	}
}