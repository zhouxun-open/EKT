package engine

import (
	"github.com/EducationEKT/EKT/io/ekt8/blockchain"
)

var mainBlockChain blockchain.BlockChain

func init() {
	mainBlockChain = blockchain.BlockChain{[]byte{(byte)(1 & 0xFF)}}
}
