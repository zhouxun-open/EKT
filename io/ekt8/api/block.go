package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/EducationEKT/EKT/io/ekt8/blockchain"
	"github.com/EducationEKT/EKT/io/ekt8/blockchain_manager"
	"github.com/EducationEKT/EKT/io/ekt8/conf"
	"github.com/EducationEKT/EKT/io/ekt8/context_log"
	"github.com/EducationEKT/EKT/io/ekt8/util"
	"github.com/EducationEKT/xserver/x_err"
	"github.com/EducationEKT/xserver/x_http/x_req"
	"github.com/EducationEKT/xserver/x_http/x_resp"
	"github.com/EducationEKT/xserver/x_http/x_router"
)

func init() {
	x_router.Post("/blocks/api/last", lastBlock)
	x_router.Get("/block/api/blockByHeight", blockByHeight)
	x_router.Post("/block/api/newBlock", newBlock)
}

func lastBlock(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	block := blockchain_manager.GetMainChain().GetLastBlock()
	return x_resp.Return(block, nil)
}

func blockByHeight(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	bc := blockchain_manager.MainBlockChain
	height := req.MustGetInt64("height")
	if bc.GetLastHeight() < height {
		fmt.Printf("Heigth %d is heigher than current height, current height is %d \n", height, bc.GetLastHeight())
		return nil, x_err.New(-404, fmt.Sprintf("Heigth %d is heigher than current height, current height is %d \n ", height, bc.GetLastHeight()))
	}
	return x_resp.Return(bc.GetBlockByHeight(height))
}

func newBlock(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	cLog := context_log.NewContextLog("Block from peer")
	defer cLog.Finish()
	var block blockchain.Block
	json.Unmarshal(req.Body, &block)
	cLog.Log("block", block)
	fmt.Printf("Recieved new block : block=%v, blockHash=%s \n", string(block.Bytes()), hex.EncodeToString(block.Hash()))
	lastHeight := blockchain_manager.GetMainChain().GetLastHeight()
	if lastHeight+1 != block.Height {
		cLog.Log("Invalid height", true)
		fmt.Printf("Block height is not right, want %d, get %d, give up voting. \n", lastHeight+1, block.Height)
		return x_resp.Fail(-1, "error invalid height", nil), nil
	}
	IP := strings.Split(req.R.RemoteAddr, ":")[0]
	if !strings.EqualFold(IP, conf.EKTConfig.Node.Address) &&
		strings.EqualFold(block.GetRound().Peers[block.GetRound().CurrentIndex].Address, IP) && block.GetRound().MyIndex() != -1 &&
		(block.GetRound().MyIndex()-block.GetRound().CurrentIndex+len(block.GetRound().Peers))%len(block.GetRound().Peers) < len(block.GetRound().Peers)/2 {
		//当前节点是打包节点广播，而且当前节点满足(currentIndex - miningIndex + len(DPoSNodes)) % len(DPoSNodes) < len(DPoSNodes) / 2
		if _, forward := req.Query["forward"]; !forward {
			for i := 0; i < len(block.GetRound().Peers); i++ {
				if i == block.GetRound().CurrentIndex || i == block.GetRound().MyIndex() {
					continue
				}
				util.HttpPost(fmt.Sprintf(`http://%s:%d/block/api/newBlock?forward=true`, block.GetRound().Peers[i].Address, block.GetRound().Peers[i].Port), req.Body)
			}
			fmt.Println("Forward block to other succeed.")
		}
	}
	blockchain_manager.MainBlockChainConsensus.BlockFromPeer(cLog, block)
	return x_resp.Return("recieved", nil)
}
