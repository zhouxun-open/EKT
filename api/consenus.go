package api

import (
	"fmt"
	"github.com/EducationEKT/EKT/p2p"
	"github.com/EducationEKT/EKT/util"
	"github.com/EducationEKT/xserver/x_err"
	"github.com/EducationEKT/xserver/x_http/x_req"
	"github.com/EducationEKT/xserver/x_http/x_resp"
	"github.com/EducationEKT/xserver/x_http/x_router"
)

/*
用来与其他区块进行基于共识机制的通信
*/
func init() {
	x_router.Post("/consenus/api/receive", receive)
}

func receive(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	return x_resp.Success("receive"), nil
}

func broadcast(req *x_req.XReq, peers p2p.Peers) {
	for _, peer := range peers {
		url := fmt.Sprintf(`http://%s:%d/consenus/api/receive`, peer.Address, peer.Port)
		util.HttpPost(url, []byte("block header"))
	}
}
