package crypto

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"testing"
)

func TestSha3(t *testing.T) {
	fmt.Println(hexutil.Encode(Sha3_256([]byte("123456"))))
	fmt.Println(len(Sha3_256([]byte("aaaa"))))
	if "0xd7190eb194ff9494625514b6d178c87f99c5973e28c398969d2233f2960a573e" != hexutil.Encode(Sha3_256([]byte("123456"))) {
		t.Fail()
	}
}
