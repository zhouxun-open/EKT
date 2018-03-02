package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	_ "github.com/EducationEKT/EKT/io/ekt8/api"
	_ "github.com/EducationEKT/EKT/io/ekt8/db"
	"github.com/EducationEKT/xserver/x_http"
)

var addr = flag.String("addr", ":8098", "http service address")

func init() {
	http.HandleFunc("/", x_http.Service)
}

func main() {
	fmt.Println("main ")
	err := http.ListenAndServe(":50880", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	time.Sleep(time.Hour * 24 * 365)
}
