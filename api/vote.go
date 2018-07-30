package api

import (
	"encoding/json"

	"github.com/EducationEKT/EKT/blockchain"
	"github.com/EducationEKT/EKT/blockchain_manager"
	"github.com/EducationEKT/EKT/log"

	"github.com/EducationEKT/xserver/x_err"
	"github.com/EducationEKT/xserver/x_http/x_req"
	"github.com/EducationEKT/xserver/x_http/x_resp"
	"github.com/EducationEKT/xserver/x_http/x_router"
)

func init() {
	x_router.Post("/vote/api/vote", voteBlock)
	x_router.Post("/vote/api/voteResult", voteResult)
	x_router.Get("/vote/api/getVotes", getVotes)
}

func voteBlock(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	var vote blockchain.BlockVote
	err := json.Unmarshal(req.Body, &vote)
	if err != nil {
		log.Info("Invalid vote, abort.")
		return x_resp.Return(nil, err)
	}
	log.Info("Recieved a vote: %s.\n", string(vote.Bytes()))
	if !vote.Validate() {
		log.Info("Invalid vote, abort.")
		return x_resp.Return(false, nil)
	}
	go blockchain_manager.GetMainChainConsensus().VoteFromPeer(vote)
	return nil, nil
}

func voteResult(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	var votes blockchain.Votes
	err := json.Unmarshal(req.Body, &votes)
	if err != nil {
		log.Info("Invalid vote, unmarshal fail, abort.")
		return x_resp.Return(nil, err)
	}
	go blockchain_manager.GetMainChainConsensus().RecieveVoteResult(votes)
	return x_resp.Success(make(map[string]interface{})), nil
}

func getVotes(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	blockHash := req.MustGetString("hash")
	votes := blockchain_manager.GetMainChainConsensus().GetVotes(blockHash)
	return x_resp.Return(votes, nil)
}
