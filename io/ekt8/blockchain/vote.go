package blockchain

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/EducationEKT/EKT/io/ekt8/b_search"
	"github.com/EducationEKT/EKT/io/ekt8/crypto"
	"github.com/EducationEKT/EKT/io/ekt8/p2p"
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
	Locker      sync.RWMutex
	broadcast   map[string]bool
	voteResults map[string]Votes
}

func NewVoteResults() VoteResults {
	return VoteResults{
		broadcast:   make(map[string]bool),
		voteResults: make(map[string]Votes),
		Locker:      sync.RWMutex{},
	}
}

func (vote VoteResults) GetVoteResults(hash string) Votes {
	vote.Locker.RLock()
	votes := vote.voteResults[hash]
	vote.Locker.RUnlock()
	return votes
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
	vote.Locker.Lock()
	defer vote.Locker.Unlock()
	votes, exist := vote.voteResults[hex.EncodeToString(voteResult.BlockHash)]
	if exist && len(votes) > 0 {
		for _, _vote := range votes {
			if strings.EqualFold(_vote.Value(), voteResult.Value()) {
				return
			}
		}
		votes = append(votes, voteResult)
		vote.voteResults[hex.EncodeToString(voteResult.BlockHash)] = votes
	} else {
		votes = make([]BlockVote, 0)
		votes = append(votes, voteResult)
		vote.voteResults[hex.EncodeToString(voteResult.BlockHash)] = votes
	}
}

func (vote VoteResults) Number(blockHash []byte) int {
	vote.Locker.RLock()
	defer vote.Locker.RUnlock()
	votes, exist := vote.voteResults[hex.EncodeToString(blockHash)]
	if !exist {
		return 0
	}
	return len(votes)
}

func (vote VoteResults) Broadcasted(blockHash []byte) bool {
	vote.Locker.RLock()
	vote.Locker.RUnlock()
	return vote.broadcast[hex.EncodeToString(blockHash)]
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
		fmt.Println("Votes.Validate: length of votes is 0, return false.")
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

func (vote Votes) Index(index int) b_search.Interface {
	if index > vote.Len() || index < 0 {
		panic("Index out of bound.")
	}
	return vote[index]
}
