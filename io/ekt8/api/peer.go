package api

import (
	"github.com/EducationEKT/EKT/io/ekt8/blockchain_manager"
	"github.com/EducationEKT/xserver/x_err"
	"github.com/EducationEKT/xserver/x_http/x_req"
	"github.com/EducationEKT/xserver/x_http/x_resp"
	"github.com/EducationEKT/xserver/x_http/x_router"
)

func init() {
	x_router.All("/peer/api/ping", ping)
	x_router.Post("/peer/api/peers", dposPeers)
}

func dposPeers(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	peers := blockchain_manager.MainBlockChain.CurrentBlock.Round.Peers
	return x_resp.Success(peers), x_err.NewXErr(nil)
}

func ping(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	resp := &x_resp.XRespContainer{
		HttpCode: 200,
		Body:     []byte("pong"),
	}
	return resp, nil
}
