package api

import (
	"encoding/json"

	"github.com/EducationEKT/EKT/io/ekt8/blockchain"
	"github.com/EducationEKT/xserver/x_err"
	"github.com/EducationEKT/xserver/x_http/x_req"
	"github.com/EducationEKT/xserver/x_http/x_resp"
	"github.com/EducationEKT/xserver/x_http/x_router"
)

func init() {
	x_router.Post("/vote/api/vote", voteBlock)
	x_router.Post("/blocks/api/voteResult", voteResult)
}

func voteBlock(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	var vote blockchain.BlockVote
	err := json.Unmarshal(req.Body, &vote)
	if err != nil {
		return x_resp.Return(nil, err)
	}
	if !vote.Validate() {
		return x_resp.Return(false, nil)
	}
	return nil, nil
}

func voteResult(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	return x_resp.Success(make(map[string]interface{})), nil
}
