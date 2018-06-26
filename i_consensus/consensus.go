package i_consensus

const (
	DPOS = 1
	POW  = 2
	POS  = 3
)

type ConsensusType int

type Consensus interface {
	//接口中没有声明，但是在Consensus的所有的实现struct中都必须拥有blockchain结构体
	Run()
}
