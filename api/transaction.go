package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/EducationEKT/EKT/conf"
	"github.com/EducationEKT/EKT/core/common"
	"github.com/EducationEKT/EKT/crypto"
	"github.com/EducationEKT/EKT/ctxlog"
	"github.com/EducationEKT/EKT/db"
	"github.com/EducationEKT/EKT/dispatcher"
	"github.com/EducationEKT/EKT/param"
	"github.com/EducationEKT/EKT/util"

	"github.com/EducationEKT/xserver/x_err"
	"github.com/EducationEKT/xserver/x_http/x_req"
	"github.com/EducationEKT/xserver/x_http/x_resp"
	"github.com/EducationEKT/xserver/x_http/x_router"
)

func init() {
	x_router.Post("/transaction/api/newTransaction", broadcastTx, newTransaction)
}

func newTransaction(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	log := ctxlog.NewContextLog("NewTransaction")
	defer log.Finish()
	log.Log("body", req.Body)
	var tx common.Transaction
	err := json.Unmarshal(req.Body, &tx)
	log.Log("tx", tx)
	if err != nil {
		return nil, x_err.New(-1, err.Error())
	}
	if tx.Amount <= 0 {
		return nil, x_err.New(-100, "error amount")
	}
	err = dispatcher.NewTransaction(log, &tx)
	if err == nil {
		txId := crypto.Sha3_256(tx.Bytes())
		db.GetDBInst().Set(txId, tx.Bytes())
	}
	return x_resp.Return(nil, err)
}

func broadcastTx(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	IP := strings.Split(req.R.RemoteAddr, ":")[0]
	broadcasted := false
	for _, peer := range param.MainChainDPosNode {
		if peer.Address == IP {
			broadcasted = true
			break
		}
	}
	if !broadcasted {
		for _, peer := range param.MainChainDPosNode {
			if !peer.Equal(conf.EKTConfig.Node) {
				url := fmt.Sprint(`http://%s:%d/transaction/api/newTransaction`, peer.Address, peer.Port)
				util.HttpPost(url, req.Body)
			}
		}
	}
	return nil, nil
}
