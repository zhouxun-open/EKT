package x_parser

import (
	"net/http"
	"github.com/EducationEKT/xserver/x_http/x_req"
)

type XParser struct {
}

func (xParser *XParser) Parse(w http.ResponseWriter, r *http.Request) (xReq *x_req.XReq) {
	xReq = x_req.New(r, w)
	xReq.ParseRequest()
	return
}
