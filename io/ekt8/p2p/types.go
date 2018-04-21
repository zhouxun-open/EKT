package p2p

import (
	"encoding/json"
	"fmt"

	"bytes"
	"github.com/EducationEKT/EKT/io/ekt8/util"
	"strings"
)

type Peer struct {
	PeerId         []byte `json:"peerId"`
	Address        string `json:"address"`
	Port           int32  `json:"port"`
	AddressVersion int    `json:"addressVersion"`
}

type Peers []Peer

func (peers Peers) Bytes() []byte {
	bts, _ := json.Marshal(peers)
	return bts
}

func (peer Peer) IsAlive() bool {
	body, err := util.HttpGet(fmt.Sprintf(`http://%s:%d/peer/api/ping`, peer.Address, peer.Port))
	if err != nil || !bytes.Equal(body, []byte("pong")) {
		return false
	}
	return true
}

func (peer Peer) Equal(peer_ Peer) bool {
	if strings.EqualFold(peer.Address, peer_.Address) && peer.Port == peer_.Port && peer.AddressVersion == peer_.AddressVersion {
		return true
	}
	return false
}

func (peer Peer) GetDBValue(key []byte) ([]byte, error) {
	url := fmt.Sprintf(`http://%s:%d/db/api/get`, peer.Address, peer.Port)
	return util.HttpPost(url, key)
}
