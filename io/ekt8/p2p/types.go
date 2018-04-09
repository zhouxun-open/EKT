package p2p

import (
	"encoding/json"
	"fmt"

	"github.com/EducationEKT/EKT/io/ekt8/util"
)

type Peer struct {
	PeerId         []byte `json:"peerId"`
	Address        []byte `json:"address"`
	Port           int32  `json:"port"`
	AddressVersion int    `json:"addressVersion"`
}

type Peers []Peer

func (peers Peers) Bytes() []byte {
	bts, _ := json.Marshal(peers)
	return bts
}

//func (peer Peer) CurrentHeight() (int64, error) {
//	url := fmt.Sprintf(`http://%s:%d/blocks/api/last`, peer.Address, peer.Port)
//	body, err := util.HttpGet(url)
//	if err != nil {
//		return -1, err
//	}
//	var block blockchain.Block
//	err = json.Unmarshal(body, &block)
//	return block.Height, err
//}
//
//func (peer Peer) CurrentBlock() (*blockchain.Block, error) {
//	url := fmt.Sprintf(`http://%s:%d/blocks/api/last`, peer.Address, peer.Port)
//	body, err := util.HttpGet(url)
//	if err != nil {
//		return nil, err
//	}
//	var block blockchain.Block
//	err = json.Unmarshal(body, &block)
//	return &block, err
//}

func (peer Peer) GetDBValue(key []byte) ([]byte, error) {
	url := fmt.Sprintf(`http://%s:%d/db/api/get`, peer.Address, peer.Port)
	return util.HttpPost(url, key)
}
