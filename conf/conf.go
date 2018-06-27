package conf

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"

	"github.com/EducationEKT/EKT/core/common"
	"github.com/EducationEKT/EKT/p2p"
)

type EKTConf struct {
	Version              string           `json:"version"`
	DBPath               string           `json:"dbPath"`
	LogPath              string           `json:"logPath"`
	Debug                bool             `json:"debug"`
	Node                 p2p.Peer         `json:"node"`
	BlockchainManagePwd  string           `json:"blockchainManagePwd"`
	GenesisBlockAccounts []common.Account `json:"genesisBlock"`
	PrivateKey           string           `json:"privateKey"`
	Env                  string           `json:"env"`
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

func (conf EKTConf) GetPrivateKey() []byte {
	privKey, err := hex.DecodeString(conf.PrivateKey)
	if err != nil {
		panic(err)
	}
	return privKey
}
