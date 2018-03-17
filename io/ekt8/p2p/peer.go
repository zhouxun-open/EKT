package p2p

import (
	"fmt"
	"strings"
	"time"

	"github.com/EducationEKT/EKT/io/ekt8/db"
	"github.com/EducationEKT/EKT/io/ekt8/util"
)

const (
	PingInterval = 1 * time.Second
)

var DPosPeers Peers
var DPOSPeersKey = []byte("DPOSPeersKey")

func init() {
	DPosPeers = BootNodes
	db.GetDBInst().Set(DPOSPeersKey, DPosPeers.Bytes())
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
