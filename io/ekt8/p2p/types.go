package p2p

import (
	"encoding/json"
	"fmt"

	"github.com/EducationEKT/EKT/io/ekt8/blockchain"
	"github.com/EducationEKT/EKT/io/ekt8/util"
)

type Peer struct {
	PeerId         []byte
	Address        []byte
	Port           int32
	AddressVersion int
}

type Peers []Peer

func (peers Peers) Bytes() []byte {
	bts, _ := json.Marshal(peers)
	return bts
}

func (peer Peer) CurrentHeight() (int64, error) {
	url := fmt.Sprintf(`http://%s:%d/blocks/api/last`, peer.Address, peer.Port)
	body, err := util.HttpGet(url)
	if err != nil {
		return -1, err
	}
	var block blockchain.Block
	err = json.Unmarshal(body, &block)
	return block.Height, err
}
