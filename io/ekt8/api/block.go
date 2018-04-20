package api

import (
	"errors"

	"fmt"
	"github.com/EducationEKT/EKT/io/ekt8/blockchain_manager"
	"github.com/EducationEKT/xserver/x_err"
	"github.com/EducationEKT/xserver/x_http/x_req"
	"github.com/EducationEKT/xserver/x_http/x_resp"
	"github.com/EducationEKT/xserver/x_http/x_router"
)

func init() {
	x_router.Post("/blocks/api/last", lastBlock)
	x_router.Post("/blocks/api/voteNext", voteNext)
	x_router.Post("/blocks/api/voteResult", voteResult)
	x_router.Get("/blocks/api/blockHeaders", blockHeaders)
	x_router.Get("/block/api/body", body)
	x_router.Get("/block/api/blockByHeight", blockByHeight)
}

func body(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	consensus := blockchain_manager.MainBlockChainConsensus
	if consensus.CurrentBlock().Height == consensus.Blockchain.CurrentBody.Height {
		return x_resp.Success(consensus.Blockchain.CurrentBody), nil
	}
	return nil, x_err.NewXErr(errors.New("can not get body"))
}

func lastBlock(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	block, err := blockchain_manager.GetMainChain().LastBlock()
	return x_resp.Return(block, err)
}

func voteNext(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	return x_resp.Success(make(map[string]interface{})), nil
}

func voteResult(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	return x_resp.Success(make(map[string]interface{})), nil
}

func blockHeaders(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	fromHeight := req.MustGetInt64("fromHeight")
	headers := blockchain_manager.GetMainChain().GetBlockHeaders(fromHeight)
	return x_resp.Success(headers), nil
}

func blockByHeight(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	bc := blockchain_manager.MainBlockChain
	height := req.MustGetInt64("heigth")
	if bc.CurrentHeight < height {
		fmt.Printf("Heigth %d is heigher than current height, current height is %d", height, bc.CurrentHeight)
		return nil, x_err.New(-404, fmt.Sprintf("Heigth %d is heigher than current height, current height is %d", height, bc.CurrentHeight))
	}
	return x_resp.Return(bc.GetBlockByHeight(height))
}

func blockBodyByHeight(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	bc := blockchain_manager.MainBlockChain
	height := req.MustGetInt64("heigth")
	if bc.CurrentHeight < height {
		fmt.Printf("Heigth %d is heigher than current height, current height is %d", height, bc.CurrentHeight)
		return nil, x_err.New(-404, fmt.Sprintf("Heigth %d is heigher than current height, current height is %d", height, bc.CurrentHeight))
	}
	return x_resp.Return(bc.GetBlockBodyByHeight(height))
}
