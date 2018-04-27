package conf

import (
	"encoding/json"
	"io/ioutil"

	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"github.com/EducationEKT/EKT/io/ekt8/p2p"
)

type EKTConf struct {
	Version              string           `json:"version"`
	DBPath               string           `json:"dbPath"`
	LogPath              string           `json:"logPath"`
	Node                 p2p.Peer         `json:"node"`
	BlockchainManagePwd  string           `json:"blockchainManagePwd"`
	GenesisBlockAccounts []common.Account `json:"genesisBlock"`
	PrivateKey           []byte           `json:"privateKey"`
}

var EKTConfig EKTConf

func InitConfig(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &EKTConfig)
	return err
}
