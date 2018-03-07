package conf

import (
	"encoding/json"
	"io/ioutil"
)

type EKTConf struct {
	DBPath  string `json:"dbPath"`
	LogPath string `json:"logPath"`
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
