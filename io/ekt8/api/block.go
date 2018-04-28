package api

import (
	"errors"

	"encoding/json"
	"fmt"
	"github.com/EducationEKT/EKT/io/ekt8/blockchain"
	"github.com/EducationEKT/EKT/io/ekt8/blockchain_manager"
	"github.com/EducationEKT/EKT/io/ekt8/util"
	"github.com/EducationEKT/xserver/x_err"
	"github.com/EducationEKT/xserver/x_http/x_req"
	"github.com/EducationEKT/xserver/x_http/x_resp"
	"github.com/EducationEKT/xserver/x_http/x_router"
	"strings"
)

func init() {
	x_router.Post("/blocks/api/last", lastBlock)
	x_router.Get("/blocks/api/blockHeaders", blockHeaders)
	x_router.Get("/block/api/body", body)
	x_router.Get("/block/api/blockByHeight", blockByHeight)
	x_router.Post("/block/api/newBlock", newBlock)
}

func body(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	consensus := blockchain_manager.MainBlockChainConsensus
	if consensus.CurrentBlock().Height == consensus.Blockchain.CurrentBody.Height {
		return x_resp.Success(consensus.Blockchain.CurrentBody), nil
	}
	return nil, x_err.NewXErr(errors.New("can not get body"))
}

func lastBlock(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	block := blockchain_manager.GetMainChain().CurrentBlock
	return x_resp.Return(block, nil)
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

func newBlock(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	var block blockchain.Block
	err := json.Unmarshal(req.Body, &block)
	if err != nil {
		fmt.Println("Block header unmarshal failed, return fail to broadcast peer.")
		return x_resp.Fail(-1, "error block header", nil), nil
	}
	lastBlock, err := blockchain_manager.GetMainChain().LastBlock()
	if lastBlock.Height+1 != block.Height {
		fmt.Printf("Block height is not right, want %d, get %d, give up voting. \n", lastBlock.Height+1, block.Height)
		return x_resp.Fail(-1, "error invalid height", nil), nil
	}
	IP := strings.Split(req.R.RemoteAddr, ":")[0]
	if strings.EqualFold(block.Round.Peers[block.Round.CurrentIndex].Address, IP) && block.Round.MyIndex() != -1 && (block.Round.MyIndex()-block.Round.CurrentIndex+len(block.Round.Peers))%len(block.Round.Peers) < len(block.Round.Peers)/2 {
		//当前节点是打包节点广播，而且当前节点满足(currentIndex - miningIndex + len(DPoSNodes)) % len(DPoSNodes) < len(DPoSNodes) / 2
		for i := 0; i < len(block.Round.Peers); i++ {
			if i == block.Round.CurrentIndex || i == block.Round.MyIndex() {
				continue
			}
			util.HttpPost(fmt.Sprintf(`http://%s:%d/block/api/newBlock`, block.Round.Peers[i].Address, block.Round.Peers[i].Port), req.Body)
		}
	}
	if !strings.EqualFold(IP, lastBlock.Round.Peers[(lastBlock.Round.CurrentIndex+1)%len(lastBlock.Round.Peers)].Address) {
		return x_resp.Return("error invalid address", errors.New("error invalid address"))
	}
	url := fmt.Sprintf("http://%s:%d/db/api/get", block.Round.Peers[block.Round.CurrentIndex].Address, block.Round.Peers[block.Round.CurrentIndex].Port)
	body, err := util.HttpPost(url, block.Body)
	if err != nil {
		fmt.Println("Get block body from peer failed, return fail to broadcast peer.")
		return x_resp.Return(err.Error(), err)
	}
	var blockBody blockchain.BlockBody
	err = json.Unmarshal(body, &blockBody)
	if err != nil {
		fmt.Printf("Block body unmarshal failed, get %s .\n", string(body))
		return x_resp.Return(err.Error(), err)
	}
	return x_resp.Return(blockchain_manager.GetMainChain().ValidateBlock(&block, &blockBody), nil)
}
