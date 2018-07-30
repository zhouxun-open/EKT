package round

import (
	"encoding/json"
	"github.com/EducationEKT/EKT/conf"
	"github.com/EducationEKT/EKT/p2p"
	"github.com/EducationEKT/EKT/util"
)

type Round struct {
	CurrentIndex int        `json:"currentIndex"` // default -1
	Peers        []p2p.Peer `json:"peers"`
}

func (round1 *Round) Equal(round2 *Round) bool {
	if round1.CurrentIndex != round2.CurrentIndex || len(round1.Peers) != len(round2.Peers) {
		return false
	}
	for i, peer := range round1.Peers {
		if !peer.Equal(round2.Peers[i]) {
			return false
		}
	}
	return true
}

func (round *Round) Clone() *Round {
	if round == nil {
		return nil
	}
	return &Round{
		Peers:        round.Peers,
		CurrentIndex: round.CurrentIndex,
	}
}

func (round *Round) Shuffle(random int) *Round {
	newRound := round.Clone()
	for high := newRound.Len() - 1; high > 0; high-- {
		low := random % high
		if random%(low+high)%2 == 1 {
			newRound.Swap(high, low)
		}
	}
	return newRound
}

func (round *Round) Swap(i, j int) {
	round.Peers[i], round.Peers[j] = round.Peers[j], round.Peers[i]
}

func (round *Round) Len() int {
	return len(round.Peers)
}

func GetRandomByHash(hash []byte) int {
	random := util.BytesToInt(hash[22:])
	if random < 0 {
		random = -random
	}
	return random
}

func (round Round) MyIndex() int {
	for i, peer := range round.Peers {
		if peer.Equal(conf.EKTConfig.Node) {
			return i
		}
	}
	return -1
}

func (round Round) String() string {
	bytes, _ := json.Marshal(round)
	return string(bytes)
}
