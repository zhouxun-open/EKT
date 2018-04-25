package blockchain

import (
	"encoding/hex"
	"fmt"

	"github.com/EducationEKT/EKT/io/ekt8/crypto"
	"github.com/EducationEKT/EKT/io/ekt8/p2p"
)

type BlockVote struct {
	BlockHash   []byte   `json:"blockHash"`
	BlockHeight int      `json:"blockHeight"`
	VoteResult  bool     `json:"voteResult"`
	Peer        p2p.Peer `json:"peer"`
	Signature   []byte   `json:"signature"`
}

type Votes []BlockVote

type VoteResults struct {
	VoteResults Votes
}

func (vote BlockVote) Validate() bool {
	return true
}

func (vote BlockVote) Sign(PrivKey []byte) error {
	str := fmt.Sprintf(`{"blockHash": "%s", "blockHeight": %d, "voteResult": %v, "peer": %s}`,
		hex.EncodeToString(vote.BlockHash), vote.BlockHeight, vote.VoteResult, vote.Peer.String())
	data := crypto.Sha3_256([]byte(str))
	signature, err := crypto.Crypto(data, PrivKey)
	if err != nil {
		return err
	} else {
		vote.Signature = signature
	}
	return nil
}
