package main

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"os"

	_ "github.com/EducationEKT/EKT/io/ekt8/api"
	"github.com/EducationEKT/EKT/io/ekt8/blockchain_manager"
	"github.com/EducationEKT/EKT/io/ekt8/conf"
	"github.com/EducationEKT/EKT/io/ekt8/crypto"
	"github.com/EducationEKT/EKT/io/ekt8/db"
	"github.com/EducationEKT/EKT/io/ekt8/log"
	"github.com/EducationEKT/EKT/io/ekt8/param"
	"github.com/EducationEKT/xserver/x_http"
)

func init() {
	err := InitService()
	if err != nil {
		fmt.Printf("Init service failed, %v \n", err)
		os.Exit(-1)
	}
	http.HandleFunc("/", x_http.Service)
}

func main() {
	fmt.Printf("server listen on :%d \n", conf.EKTConfig.Node.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", conf.EKTConfig.Node.Port), nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func InitService() error {
	err := initConfig()
	if err != nil {
		return err
	}
	fmt.Printf("Current EKT version is %s. \n", conf.EKTConfig.Version)
	err = initDB()
	if err != nil {
		return err
	}
	err = initPeerId()
	if err != nil {
		return err
	}
	err = initLog()
	if err != nil {
		return err
	}
	param.InitBootNodes()
	blockchain_manager.Init()

	return nil
}

func initPeerId() error {
	peerInfoKey := []byte("peerIdInfo")
	v, err := db.GetDBInst().Get(peerInfoKey)
	if err != nil || nil == v || 0 == len(v) {
		pub, priv := crypto.GenerateKeyPair()
		conf.EKTConfig.PrivateKey = priv
		conf.EKTConfig.Node.PeerId = hex.EncodeToString(crypto.Sha3_256(pub))
		fmt.Printf("Current peerId is: %s . \n", conf.EKTConfig.Node.PeerId)
		return db.GetDBInst().Set(peerInfoKey, priv)
	} else {
		conf.EKTConfig.PrivateKey = v
		data := crypto.Sha3_256(v)
		cryptoData, err := crypto.Crypto(data, v)
		if err != nil {
			fmt.Println(err)
			return err
		}
		pub, err := crypto.RecoverPubKey(data, cryptoData)
		conf.EKTConfig.Node.PeerId = hex.EncodeToString(crypto.Sha3_256(pub))
		fmt.Printf("Current peerId is %s. \n", conf.EKTConfig.Node.PeerId)
	}

	return nil
}

func initConfig() error {
	var confPath string
	if len(os.Args) < 2 {
		confPath = "genesis.json"
		fmt.Println("No conf file specified, genesis.json will be default one.")
	} else {
		confPath = os.Args[1]
	}
	err := conf.InitConfig(confPath)
	return err
}

func initDB() error {
	return db.InitEKTDB(conf.EKTConfig.DBPath)
}

func initLog() error {
	return log.InitLog()
}
