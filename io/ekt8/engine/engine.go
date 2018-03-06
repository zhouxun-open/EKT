package engine

import (
	"encoding/json"

	"github.com/EducationEKT/EKT/io/ekt8/MPTPlus"
	"github.com/EducationEKT/EKT/io/ekt8/blockchain"
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"github.com/EducationEKT/EKT/io/ekt8/db"
)

var mainBlockChain blockchain.BlockChain

func init() {
	mainBlockChain = blockchain.BlockChain{[]byte{(byte)(1 & 0xFF)}}
}
