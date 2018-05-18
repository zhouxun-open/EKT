package x_parser

import (
	"fmt"
	"net/http"
	"github.com/EducationEKT/xserver/x_http/x_req"
)

type XParser struct {
}

func (xParser *XParser) Parse(w http.ResponseWriter, r *http.Request) (xReq *x_req.XReq) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic occurred.", r)
		}
	}()
	xReq = x_req.New(r, w)
	xReq.ParseRequest()
	return
}
