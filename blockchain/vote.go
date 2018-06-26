package blockchain

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/EducationEKT/EKT/crypto"
	"github.com/EducationEKT/EKT/log"
	"github.com/EducationEKT/EKT/p2p"
)

var VoteResultManager VoteResults

func init() {
	VoteResultManager = NewVoteResults()
}

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
	broadcast   *sync.Map
	voteResults *sync.Map
}

func NewVoteResults() VoteResults {
	return VoteResults{
		broadcast:   &sync.Map{},
		voteResults: &sync.Map{},
	}
}

func (vote VoteResults) GetVoteResults(hash string) Votes {
	obj, exist := vote.voteResults.Load(hash)
	if exist {
		return obj.(Votes)
	}
	return nil
}

func (vote VoteResults) SetVoteResults(hash string, votes Votes) {
	vote.voteResults.Store(hash, votes)
}

func (vote BlockVote) Validate() bool {
	pubKey_, err := crypto.RecoverPubKey(crypto.Sha3_256(vote.Data()), vote.Signature)
	if err != nil {
		fmt.Printf("BlockVote.Validate: recover public key failed, return false.")
		return false
	}
	if !strings.EqualFold(hex.EncodeToString(crypto.Sha3_256(pubKey_)), vote.Peer.PeerId) {
		fmt.Printf("Recovered public key: %s", hex.EncodeToString(pubKey_))
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

func (vote *BlockVote) Sign(PrivKey []byte) error {
	signature, err := crypto.Crypto(crypto.Sha3_256(vote.Data()), PrivKey)
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
	votes := vote.GetVoteResults(hex.EncodeToString(voteResult.BlockHash))
	if len(votes) > 0 {
		for _, _vote := range votes {
			if strings.EqualFold(_vote.Value(), voteResult.Value()) {
				return
			}
		}
		votes = append(votes, voteResult)
		sort.Sort(votes)
	} else {
		votes = make([]BlockVote, 0)
		votes = append(votes, voteResult)
	}
	vote.SetVoteResults(hex.EncodeToString(voteResult.BlockHash), votes)
}

func (vote VoteResults) Number(blockHash []byte) int {
	votes := vote.GetVoteResults(hex.EncodeToString(blockHash))
	return len(votes)
}

func (vote VoteResults) Broadcasted(blockHash []byte) bool {
	_, exist := vote.broadcast.Load(hex.EncodeToString(blockHash))
	return exist
}

func (vote Votes) Len() int {
	return len(vote)
}

func (vote Votes) Swap(i, j int) {
	vote[i], vote[j] = vote[j], vote[i]
}

func (vote Votes) Less(i, j int) bool {
	return vote[i].Peer.String() < vote[j].Peer.String()
}

func (vote Votes) Bytes() []byte {
	data, _ := json.Marshal(vote)
	return data
}

func (votes Votes) Validate() bool {
	if len(votes) == 0 {
		log.Debug("Votes.Validate: length of votes is 0, return false.")
		return false
	}
	for i, vote := range votes {
		if !vote.Validate() || !vote.VoteResult {
			return false
		}
		for j, _vote := range votes {
			if i != j {
				if bytes.Equal(vote.Data(), _vote.Data()) {
					return false
				}
			}
		}
	}
	return true
}
