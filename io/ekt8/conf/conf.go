package conf

import (
	"encoding/json"
	"io/ioutil"

	"github.com/EducationEKT/EKT/io/ekt8/core/common"
)

type EKTConf struct {
	DBPath               string           `json:"dbPath"`
	LogPath              string           `json:"logPath"`
	SelfNodeId           string           `json:"nodeId"`
	BlockchainManagePwd  string           `json:"blockchainManagePwd"`
	GenesisBlockAccounts []common.Account `json:"genesisBlock"`
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
