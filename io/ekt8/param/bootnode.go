package param

import (
	"github.com/EducationEKT/EKT/io/ekt8/conf"
	"github.com/EducationEKT/EKT/io/ekt8/p2p"
)

var mapping = make(map[string][]p2p.Peer)
var MainChainDPosNode []p2p.Peer

func InitBootNodes() {
	mapping["mainnet"] = MainNet
	mapping["testnet"] = TestNet
	mapping["localnet"] = LocalNet
	MainChainDPosNode = mapping[conf.EKTConfig.Env]
}
