package api

import (
	"github.com/EducationEKT/xserver/x_err"
	"github.com/EducationEKT/xserver/x_http/x_req"
	"github.com/EducationEKT/xserver/x_http/x_resp"
	"github.com/EducationEKT/xserver/x_http/x_router"
)

func init() {
	x_router.Post("/user/api/regist")
}

func NewAccount(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	return x_resp.Fail(-1, "This function is not open now.", nil), nil
}
