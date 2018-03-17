package p2p

import "encoding/json"

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
