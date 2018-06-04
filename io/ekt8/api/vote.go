package api

import (
	"encoding/json"
	"fmt"

	"github.com/EducationEKT/EKT/io/ekt8/blockchain"
	"github.com/EducationEKT/EKT/io/ekt8/blockchain_manager"
	"github.com/EducationEKT/xserver/x_err"
	"github.com/EducationEKT/xserver/x_http/x_req"
	"github.com/EducationEKT/xserver/x_http/x_resp"
	"github.com/EducationEKT/xserver/x_http/x_router"
)

func init() {
	x_router.Post("/vote/api/vote", voteBlock)
	x_router.Post("/vote/api/voteResult", voteResult)
}

func voteBlock(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	var vote blockchain.BlockVote
	err := json.Unmarshal(req.Body, &vote)
	if err != nil {
		fmt.Println("Invalid vote, abort.")
		return x_resp.Return(nil, err)
	}
	fmt.Printf("Recieved a vote: %s.\n", string(vote.Bytes()))
	if !vote.Validate() {
		fmt.Println("Invalid vote, abort.")
		return x_resp.Return(false, nil)
	}
	blockchain_manager.GetMainChain().VoteFromPeer(vote)
	return nil, nil
}

func voteResult(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	var votes blockchain.Votes
	err := json.Unmarshal(req.Body, &votes)
	if err != nil {
		fmt.Println("Invalid vote, unmarshal fail, abort.")
		return x_resp.Return(nil, err)
	}
	blockchain_manager.GetMainChainConsensus().RecieveVoteResult(votes)
	return x_resp.Success(make(map[string]interface{})), nil
}
