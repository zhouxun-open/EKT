package p2p

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/EducationEKT/EKT/crypto"
)

func TestPeer(t *testing.T) {
	fmt.Println(hex.EncodeToString(crypto.Sha3_256([]byte("192.168.6.67:1995"))))
	fmt.Println(hex.EncodeToString(crypto.Sha3_256([]byte("192.168.6.68:1995"))))
	fmt.Println(hex.EncodeToString(crypto.Sha3_256([]byte("192.168.6.69:1995"))))
}
