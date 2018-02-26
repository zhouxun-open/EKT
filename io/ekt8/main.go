package main

import (
	_ "github.com/EducationEKT/EKT/io/ekt8/api"
	_ "github.com/EducationEKT/EKT/io/ekt8/db"
	"github.com/EducationEKT/xserver/x_http"
	"net/http"
	"fmt"
)

func init() {
	http.HandleFunc("/", x_http.Service)
}

func main() {
	fmt.Println("main ")
}
