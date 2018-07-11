package param

import (
	"github.com/EducationEKT/EKT/p2p"
	"github.com/EducationEKT/EKT/param"
)

var (
	Localnet bool = false
	Testnet  bool = false
	Mainnet  bool = false
)

func GetPeers() []p2p.Peer {
	if Localnet {
		return param.LocalNet
	} else if Testnet {
		return param.TestNet
	} else if Mainnet {
		return param.MainNet
	}
	return param.LocalNet
}
