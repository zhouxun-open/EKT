package p2p

import "time"

const (
	PingInterval = 1 * time.Second
)

type Peer struct {
	PeerId         []byte
	Address        []byte
	AddressVersion int
	Port           int32
}
