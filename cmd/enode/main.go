package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"net/http"
	"os"

	"bytes"

	_ "github.com/EducationEKT/EKT/api"
	"github.com/EducationEKT/EKT/blockchain_manager"
	"github.com/EducationEKT/EKT/conf"
	"github.com/EducationEKT/EKT/crypto"
	"github.com/EducationEKT/EKT/db"
	"github.com/EducationEKT/EKT/log"
	"github.com/EducationEKT/EKT/param"
	"github.com/EducationEKT/xserver/x_http"
)

const (
	version = "0.1"
)

func init() {
	var (
		help bool
		ver  bool
		cfg  string
	)
	flag.BoolVar(&help, "h", false, "this help")
	flag.BoolVar(&ver, "v", false, "show version and exit")
	flag.StringVar(&cfg, "c", "genesis.json", "set genesis.json file and start")
	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	if ver {
		fmt.Println(version)
		os.Exit(0)
	}

	err := InitService(cfg)
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

func InitService(confPath string) error {
	err := initConfig(confPath)
	if err != nil {
		return err
	}
	err = initLog()
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
	param.InitBootNodes()
	blockchain_manager.Init()

	return nil
}

func initPeerId() error {
	if !bytes.Equal(conf.EKTConfig.PrivateKey, []byte("")) {
		fmt.Printf("Current peerId is: %s . \n", conf.EKTConfig.Node.PeerId)
		return nil
	}
	peerInfoKey := []byte("peerIdInfo")
	v, err := db.GetDBInst().Get(peerInfoKey)
	if err != nil || nil == v || 0 == len(v) {
		pub, priv := crypto.GenerateKeyPair()
		conf.EKTConfig.PrivateKey = priv
		conf.EKTConfig.Node.PeerId = hex.EncodeToString(crypto.Sha3_256(pub))
		err = db.GetDBInst().Set(peerInfoKey, priv)
		if err != nil {
			fmt.Println(err)
			return err
		}
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
	}

	fmt.Println("Peer private key is: ", hex.EncodeToString(conf.EKTConfig.PrivateKey))
	fmt.Printf("Current peerId is %s . \n", conf.EKTConfig.Node.PeerId)

	return nil
}

func initConfig(confPath string) error {
	return conf.InitConfig(confPath)
}

func initDB() error {
	return db.InitEKTDB(conf.EKTConfig.DBPath)
}

func initLog() error {
	log.InitLog()
	return nil
}
