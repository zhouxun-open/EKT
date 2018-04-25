package blockchain

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/EducationEKT/EKT/io/ekt8/b_search"
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

func (vote BlockVote) Validate(pubKey []byte) bool {
	pubKey_, err := crypto.RecoverPubKey(vote.Data(), vote.Signature)
	if err != nil || !bytes.Equal(pubKey_, pubKey) {
		return false
	}
	return true
}

func (vote BlockVote) Data() []byte {
	str := fmt.Sprintf(`{"blockHash": "%s", "blockHeight": %d, "voteResult": %v, "peer": %s}`,
		hex.EncodeToString(vote.BlockHash), vote.BlockHeight, vote.VoteResult, vote.Peer.String())
	return crypto.Sha3_256([]byte(str))
}

func (vote BlockVote) Value() string {
	return string(vote.Data())
}

func (vote BlockVote) Sign(PrivKey []byte) error {
	signature, err := crypto.Crypto(vote.Data(), PrivKey)
	if err != nil {
		return err
	} else {
		vote.Signature = signature
	}
	return nil
}

func (vote VoteResults) Insert(voteResult BlockVote) {
	if -1 == b_search.Search(voteResult, vote) {
		vote.VoteResults = append(vote.VoteResults, voteResult)
		sort.Sort(vote)
	}
}

func (vote VoteResults) Len() int {
	return len(vote.VoteResults)
}

func (vote VoteResults) Swap(i, j int) {
	vote.VoteResults[i], vote.VoteResults[j] = vote.VoteResults[j], vote.VoteResults[i]
}

func (vote VoteResults) Less(i, j int) bool {
	return vote.VoteResults[i].Peer.String() < vote.VoteResults[j].Peer.String()
}

func (vote VoteResults) Index(index int) b_search.Interface {
	if index > vote.Len() || index < 0 {
		panic("Index out of bound.")
	}
	return vote.VoteResults[index]
}
