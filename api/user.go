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
	x_router.Get("/user/api/info", userInfo)
	x_router.Get("/user/api/lastTxNonce", userNonce)
}

func userInfo(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	address := req.MustGetString("address")
	log := ctxlog.NewContextLog("userInfo")
	defer log.Finish()
	hexAddress, err := hex.DecodeString(address)
	if err != nil {
		return x_resp.Return(nil, err)
	}
	account, err := blockchain_manager.GetMainChain().GetLastBlock().GetAccount(log, hexAddress)
	if err != nil {
		return x_resp.Return(nil, err)
	}
	return x_resp.Return(account, nil)
}

func userNonce(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	ctxLog:=ctxlog.NewContextLog("Get user last transaction nonce")
	defer ctxLog.Finish()

	hexAddress := req.MustGetString("address")
	address, err:=hex.DecodeString(hexAddress)
	if err!=nil {
		return x_resp.Return(nil, err)
	}
	ctxLog.Log("address", hexAddress)

	// get user nonce by user stat tree
	account, err:=blockchain_manager.GetMainChain().GetLastBlock().GetAccount(ctxLog, address)
	if err!=nil {
		return x_resp.Return(nil, err)
	}
	nonce:=account.GetNonce()

	txs:=blockchain_manager.GetMainChain().Pool.GetTxs(hexAddress)
	if txs!=nil &&len(txs) > 0 {
		for _, tx:=range txs {
			if tx.Nonce> nonce {
				nonce=tx.Nonce
			}
		}
	}

	return x_resp.Return(nonce, nil)
}