package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/EducationEKT/EKT/io/ekt8/api"
	"github.com/EducationEKT/EKT/io/ekt8/conf"
	"github.com/EducationEKT/EKT/io/ekt8/db"
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
	err := http.ListenAndServe(":50880", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func InitService() error {
	err := initConfig()
	if err != nil {
		return err
	}
	err = initDB()
	if err != nil {
		return err
	}

	return nil
}

func initConfig() error {
	var confPath string
	if len(os.Args) < 2 {
		confPath = "genesis.conf"
		fmt.Println("No conf file specified, genesis.conf will be default one.")
	} else {
		confPath = os.Args[1]
	}
	err := conf.InitConfig(confPath)
	return err
}

func initDB() error {
	return db.InitEKTDB(conf.EKTConfig.DBPath)
}
