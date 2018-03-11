package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/EducationEKT/xserver/x_err"
	"github.com/EducationEKT/xserver/x_http"
	"github.com/EducationEKT/xserver/x_http/x_req"
	"github.com/EducationEKT/xserver/x_http/x_resp"
	"github.com/EducationEKT/xserver/x_http/x_router"
	"net/http"
	"os"
)

var addr = flag.String("addr", ":8098", "http service address")

func main() {
	x_router.All("/Hello", Hello, World)
	x_router.Get("/test", Hello, World)
	http.HandleFunc("/", x_http.Service)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func Hello(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	req.Context.Sticker["userId"] = 123
	return nil, nil
}

func World(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	req.Context.Sticker["Hello"] = "World"
	body, _ := json.Marshal(req.Context.Sticker)
	return &x_resp.XRespContainer{
		HttpCode: 200,
		Body:     body,
	}, nil
}
