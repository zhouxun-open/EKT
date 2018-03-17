package p2p

import (
	"fmt"
	"strings"
	"time"

	"github.com/EducationEKT/EKT/io/ekt8/util"
)

const (
	PingInterval = 1 * time.Second
)

var DPOSPeersKey = []byte("DPOSPeersKey")

func IsDPosPeer(address string) bool {
	for _, peer := range MainChainDPosNode {
		if strings.Contains(address, string(peer.Address)) {
			return true
		}
	}
	return false
}

func BroadcastRequest(path string, body []byte) {
	for _, peer := range MainChainDPosNode {
		url := fmt.Sprintf("http://%s:%d%s", string(peer.Address), peer.Port, path)
		util.HttpPost(url, body)
	}
}
