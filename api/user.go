package api

import (
	"encoding/hex"
	"github.com/EducationEKT/EKT/blockchain_manager"
	"github.com/EducationEKT/EKT/ctxlog"
	"github.com/EducationEKT/xserver/x_err"
	"github.com/EducationEKT/xserver/x_http/x_req"
	"github.com/EducationEKT/xserver/x_http/x_resp"
	"github.com/EducationEKT/xserver/x_http/x_router"
)

func init() {
	x_router.Get("/account/api/info", userInfo)
	x_router.Get("/account/api/nonce", userNonce)
}

func userInfo(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	address := req.MustGetString("address")
	log := ctxlog.NewContextLog("userInfo")
	defer log.Finish()
	hexAddress, err := hex.DecodeString(address)
	if err != nil {
		return x_resp.Return(nil, err)
	}
	account, err := blockchain_manager.GetMainChain().GetLastBlock().GetAccount(hexAddress)
	if err != nil {
		return x_resp.Return(nil, err)
	}
	return x_resp.Return(account, nil)
}

func userNonce(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	ctxLog := ctxlog.NewContextLog("Get user last transaction nonce")
	defer ctxLog.Finish()

	hexAddress := req.MustGetString("address")
	address, err := hex.DecodeString(hexAddress)
	if err != nil {
		return x_resp.Return(nil, err)
	}
	ctxLog.Log("address", hexAddress)

	// get user nonce by user stat tree
	account, err := blockchain_manager.GetMainChain().GetLastBlock().GetAccount(address)
	if err != nil {
		return x_resp.Return(nil, err)
	}
	nonce := account.GetNonce()

	txs := blockchain_manager.GetMainChain().Pool.GetReadyEvents(hexAddress)
	if len(txs) > 0 {
		nonce = txs[len(txs)-1].GetNonce()
	}

	return x_resp.Return(nonce, nil)
}
