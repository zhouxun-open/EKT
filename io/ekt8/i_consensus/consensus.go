package i_consensus

const (
	DPOS = 1
	POW  = 2
	POS  = 3
)

type ConsensusType int

type Consensus interface {
	Run()
}
