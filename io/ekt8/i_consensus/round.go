package i_consensus

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/EducationEKT/EKT/io/ekt8/conf"
	"github.com/EducationEKT/EKT/io/ekt8/log"
	"github.com/EducationEKT/EKT/io/ekt8/p2p"
	"github.com/EducationEKT/EKT/io/ekt8/param"
	"github.com/EducationEKT/EKT/io/ekt8/util"
)

type Round struct {
	CurrentIndex int        `json:"currentIndex"` // default -1
	Peers        []p2p.Peer `json:"peers"`
	Random       int        `json:"random"`
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

func (round *Round) IndexPlus(CurrentHash []byte) *Round {
	if round.CurrentIndex == len(round.Peers)-1 {
		Random := util.BytesToInt(CurrentHash[22:])
		round_ := &Round{
			CurrentIndex: 0,
			Peers:        round.Peers,
			Random:       Random,
		}
		sort.Sort(round_)
		return round_
	} else {
		round.CurrentIndex++
	}
	return round
}

func (round *Round) NewRandom(CurrentHash []byte) *Round {
	round1 := &Round{
		Peers:        round.Peers,
		Random:       util.BytesToInt(CurrentHash[22:]),
		CurrentIndex: round.CurrentIndex,
	}
	return round1
}

func (round *Round) NewRound() *Round {
	if round == nil {
		return nil
	}
	return &Round{
		Peers:        round.Peers,
		CurrentIndex: round.CurrentIndex,
		Random:       round.Random,
	}
}

func MyRound(round_ *Round, Hash []byte) *Round {
	if round_ == nil {
		return &Round{
			Peers:        param.MainChainDPosNode,
			CurrentIndex: 0,
			Random:       util.BytesToInt(Hash[22:]),
		}
	} else {
		return round_.MyRound(Hash)
	}
}

func (round *Round) MyRound(CurrentHash []byte) *Round {
	log.Debug("Current Round is %s", round.String())
	_round := round.NewRound()
	if round.CurrentIndex == round.Len()-1 {
		_round = round.NewRandom(CurrentHash)
		sort.Sort(_round)
	}
	_round.CurrentIndex = _round.MyIndex()
	log.Debug("My Round is %s", _round.String())
	return _round
}

func (round *Round) NextRound(CurrentHash []byte) *Round {
	if round.CurrentIndex == len(round.Peers)-1 {
		bytes := CurrentHash[22:]
		Random := util.BytesToInt(bytes)
		round = &Round{
			CurrentIndex: 0,
			Peers:        round.Peers,
			Random:       Random,
		}
	} else {
		round.CurrentIndex = round.MyIndex()
	}
	return round
}

func (round Round) IsMyTurn() bool {
	if round.Peers[(round.CurrentIndex+1)%len(round.Peers)].Equal(conf.EKTConfig.Node) {
		return true
	}
	return false
}

func (round Round) NextPeerRight(peer p2p.Peer, hash []byte) bool {
	if round.CurrentIndex < round.Len()-1 {
		if round.Peers[round.CurrentIndex+1].Equal(peer) {
			return true
		}
		return false
	} else {
		_round := round.NewRandom(hash)
		sort.Sort(_round)
		return _round.Peers[0].Equal(peer)
	}
}

func (round Round) MyIndex() int {
	for i, peer := range round.Peers {
		if peer.Equal(conf.EKTConfig.Node) {
			return i
		}
	}
	return -1
}

func (round Round) Len() int {
	return len(round.Peers)
}

func (round Round) Swap(i, j int) {
	round.Peers[i], round.Peers[j] = round.Peers[j], round.Peers[i]
}

func (round Round) Less(i, j int) bool {
	return round.Random%(i+j)%2 == 1
}

func (round Round) String() string {
	peers, _ := json.Marshal(round.Peers)
	return fmt.Sprintf(`{"currentIndex": %d, "peers": %s, "random": %d}`, round.CurrentIndex, string(peers), round.Random)
}
