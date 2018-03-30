package api

import (
	"github.com/EducationEKT/xserver/x_http/x_req"
	"github.com/EducationEKT/xserver/x_http/x_resp"
	"github.com/EducationEKT/xserver/x_err"
	"fmt"
	"github.com/EducationEKT/xserver/x_http/x_router"
)

/*
用来与其他区块进行基于共识机制的通信
*/
func init(){
	x_router.Post("/consenus/api/receive",receive)
}

func receive(req *x_req.XReq)(*x_resp.XRespContainer, *x_err.XErr){
	return x_resp.Success("receive"),nil
}
