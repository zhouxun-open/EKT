package p2p

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"github.com/EducationEKT/EKT/io/ekt8/util"
)

const (
	PingInterval = 1 * time.Second
)

var DPosPeers []Peer

type Peer struct {
	PeerId         []byte
	Address        []byte
	Port           int32
	AddressVersion int
}

func IsDPosPeer(address string) bool {
	for _, peer := range DPosPeers {
		if strings.Contains(address, string(peer.Address)) {
			return true
		}
	}
	return false
}

func BroadcastRequest(path string, body []byte) {
	for _, peer := range DPosPeers {
		url := fmt.Sprintf("http://%s:%d%s", string(peer.Address), peer.Port, path)
		util.HttpPost(url, body)
	}
}
