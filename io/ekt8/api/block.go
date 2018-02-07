package api

import (
	"github.com/EducationEKT/xserver/x_err"
	"github.com/EducationEKT/xserver/x_http/x_req"
	"github.com/EducationEKT/xserver/x_http/x_resp"
	"github.com/EducationEKT/xserver/x_http/x_router"
)

func init() {
	x_router.Post("/blocks/api/last", lastBlock)
	x_router.Post("/blocks/api/voteNext", voteNext)
	x_router.Post("/blocks/api/voteResult", voteResult)
}

func lastBlock(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	return x_resp.Success(make(map[string]interface{})), nil
}

func voteNext(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	return x_resp.Success(make(map[string]interface{})), nil
}

func voteResult(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	return x_resp.Success(make(map[string]interface{})), nil
}
