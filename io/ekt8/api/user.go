package api

import (
	"encoding/hex"
	"github.com/EducationEKT/EKT/io/ekt8/blockchain_manager"
	"github.com/EducationEKT/EKT/io/ekt8/ctxlog"
	"github.com/EducationEKT/xserver/x_err"
	"github.com/EducationEKT/xserver/x_http/x_req"
	"github.com/EducationEKT/xserver/x_http/x_resp"
	"github.com/EducationEKT/xserver/x_http/x_router"
)

func init() {
	x_router.Get("/user/api/info", userInfo)
}

func userInfo(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	address := req.MustGetString("address")
	log := ctxlog.NewContextLog("userInfo")
	defer log.Finish()
	hexAddress, err := hex.DecodeString(address)
	if err != nil {
		x_resp.Fail(-1, "error address", nil)
	}
	account, err := blockchain_manager.GetMainChain().GetLastBlock().GetAccount(log, hexAddress)
	if err != nil {
		x_resp.Fail(-1, err.Error(), nil)
	}
	return x_resp.Return(account, nil)
}
