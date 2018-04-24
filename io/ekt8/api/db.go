package api

import (
	"fmt"

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
	if len(v) != 32 {
		fmt.Println("Remote peer want a db value that len(key) is not 32 byte, return fail.")
		return x_resp.Fail(-403, "Invalid Key", string(v)), nil
	}
	resp := &x_resp.XRespContainer{
		HttpCode: 200,
		Body:     v,
	}
	return resp, x_err.NewXErr(err)
}
