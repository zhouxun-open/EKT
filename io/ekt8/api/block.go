package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/EducationEKT/EKT/io/ekt8/blockchain"
	"github.com/EducationEKT/EKT/io/ekt8/blockchain_manager"
	"github.com/EducationEKT/EKT/io/ekt8/conf"
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
	block := blockchain_manager.GetMainChain().CurrentBlock
	return x_resp.Return(block, nil)
}

func blockByHeight(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	bc := blockchain_manager.MainBlockChain
	height := req.MustGetInt64("height")
	if bc.CurrentHeight < height {
		fmt.Printf("Heigth %d is heigher than current height, current height is %d \n", height, bc.CurrentHeight)
		return nil, x_err.New(-404, fmt.Sprintf("Heigth %d is heigher than current height, current height is %d \n ", height, bc.CurrentHeight))
	}
	return x_resp.Return(bc.GetBlockByHeight(height))
}

func newBlock(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	var block blockchain.Block
	json.Unmarshal(req.Body, &block)
	fmt.Printf("Recieved new block : block=%v, blockHash=%s \n", string(block.Bytes()), hex.EncodeToString(block.Hash()))
	lastBlock := blockchain_manager.GetMainChain().CurrentBlock
	if lastBlock.Height+1 != block.Height {
		fmt.Printf("Block height is not right, want %d, get %d, give up voting. \n", lastBlock.Height+1, block.Height)
		return x_resp.Fail(-1, "error invalid height", nil), nil
	}
	IP := strings.Split(req.R.RemoteAddr, ":")[0]
	if !strings.EqualFold(IP, conf.EKTConfig.Node.Address) && strings.EqualFold(block.Round.Peers[block.Round.CurrentIndex].Address, IP) && block.Round.MyIndex() != -1 && (block.Round.MyIndex()-block.Round.CurrentIndex+len(block.Round.Peers))%len(block.Round.Peers) < len(block.Round.Peers)/2 {
		//当前节点是打包节点广播，而且当前节点满足(currentIndex - miningIndex + len(DPoSNodes)) % len(DPoSNodes) < len(DPoSNodes) / 2
		if _, forward := req.Query["forward"]; !forward {
			for i := 0; i < len(block.Round.Peers); i++ {
				if i == block.Round.CurrentIndex || i == block.Round.MyIndex() {
					continue
				}
				util.HttpPost(fmt.Sprintf(`http://%s:%d/block/api/newBlock?forward=true`, block.Round.Peers[i].Address, block.Round.Peers[i].Port), req.Body)
			}
			fmt.Println("Forward block to other succeed.")
		}
	}
	go func() {
		blockchain_manager.MainBlockChainConsensus.BlockFromPeer(block)
		//blockchain_manager.MainBlockChainConsensus.Block <- block
	}()
	return x_resp.Return("recieved", nil)
}
