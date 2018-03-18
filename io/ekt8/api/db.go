package api

import (
	"github.com/EducationEKT/EKT/io/ekt8/db"
	"github.com/EducationEKT/xserver/x_err"
	"github.com/EducationEKT/xserver/x_http/x_req"
	"github.com/EducationEKT/xserver/x_http/x_resp"
	"github.com/EducationEKT/xserver/x_http/x_router"
)

func init() {
	x_router.Post("/db/api/get", GetValue)
}

func GetValue(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	v, err := db.GetDBInst().Get(req.Body)
	return x_resp.Return(v, err)
}
