package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/EducationEKT/EKT/conf"
	"github.com/EducationEKT/EKT/crypto"
	"github.com/EducationEKT/EKT/db"
	"github.com/EducationEKT/EKT/dispatcher"
	"github.com/EducationEKT/EKT/param"
	"github.com/EducationEKT/EKT/util"

	"encoding/hex"
	"github.com/EducationEKT/EKT/blockchain_manager"
	"github.com/EducationEKT/EKT/core/userevent"
	"github.com/EducationEKT/xserver/x_err"
	"github.com/EducationEKT/xserver/x_http/x_req"
	"github.com/EducationEKT/xserver/x_http/x_resp"
	"github.com/EducationEKT/xserver/x_http/x_router"
)

func init() {
	x_router.Post("/transaction/api/newTransaction", broadcastTx, newTransaction)
	x_router.Get("/transaction/api/queueTxs", queueTxs)
	x_router.Get("/transaction/api/blockTxs", blockTxs)
	x_router.Get("/transaction/api/status", txStatus)
}

func queueTxs(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	address := req.MustGetString("address")
	return x_resp.Return(blockchain_manager.GetMainChain().Pool.GetReadyEvents(address), nil)
}

func blockTxs(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	address := req.MustGetString("address")
	return x_resp.Return(blockchain_manager.GetMainChain().Pool.GetBlockEvents(address), nil)
}

func txStatus(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	// get transaction by transactionId
	transactionId := req.MustGetString("txId")
	txId, err := hex.DecodeString(transactionId)
	if err != nil {
		return x_resp.Return(nil, err)
	}
	if tx := userevent.GetTransaction(txId); tx == nil {
		synchronizeTransaction(txId)
	}
	tx := userevent.GetTransaction(txId)
	if tx == nil {
		return x_resp.Return("error transaction not found", nil)
	}

	// get account by address
	account, err := blockchain_manager.GetMainChain().GetLastBlock().GetAccount(tx.GetFrom())
	if err != nil {
		return x_resp.Return(nil, err)
	}

	// transaction has been processed
	if account.Nonce >= tx.Nonce {
		// 200 = processed
		return x_resp.Return(200, nil)
	}
	// 100 = pending
	return x_resp.Return(100, nil)
}

func newTransaction(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	var tx userevent.Transaction
	err := json.Unmarshal(req.Body, &tx)
	if err != nil {
		return nil, x_err.New(-1, err.Error())
	}
	if tx.Amount <= 0 {
		return nil, x_err.New(-100, "error amount")
	}
	err = dispatcher.NewTransaction(&tx)
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
				url := fmt.Sprintf(`http://%s:%d/transaction/api/newTransaction`, peer.Address, peer.Port)
				util.HttpPost(url, req.Body)
			}
		}
	}
	return nil, nil
}

func synchronizeTransaction(txId []byte) {
	for _, peer := range param.MainChainDPosNode {
		if value, err := peer.GetDBValue(txId); err != nil {
			db.GetDBInst().Set(txId, value)
		}
	}
}
