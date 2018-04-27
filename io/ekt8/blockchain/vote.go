package blockchain

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/EducationEKT/EKT/io/ekt8/b_search"
	"github.com/EducationEKT/EKT/io/ekt8/crypto"
	"github.com/EducationEKT/EKT/io/ekt8/p2p"
)

type BlockVote struct {
	BlockchainId []byte   `json:"blockchainId"`
	BlockHash    []byte   `json:"blockHash"`
	BlockHeight  int64    `json:"blockHeight"`
	VoteResult   bool     `json:"voteResult"`
	Peer         p2p.Peer `json:"peer"`
	Signature    []byte   `json:"signature"`
}

type Votes []BlockVote

type VoteResults struct {
	VoteResults Votes
}

func (vote BlockVote) Validate() bool {
	pubKey_, err := crypto.RecoverPubKey(vote.Data(), vote.Signature)
	if err != nil {
		return false
	}
	if !strings.EqualFold(hex.EncodeToString(crypto.Sha3_256(pubKey_)), vote.Peer.PeerId) {
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

func (vote BlockVote) Bytes() []byte {
	data, _ := json.Marshal(vote)
	return data
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
