package consensus

import "github.com/EducationEKT/EKT/io/ekt8/blockchain"

type BlockManager struct {
	Blockchain blockchain.BlockChain
}

func (m *BlockManager) BlockFromPeer(block *blockchain.Block) {
	m.Blockchain.Locker.RLock()
	defer m.Blockchain.Locker.Unlock()
	if m.Blockchain.CurrentHeight+1 < block.Height || m.Blockchain.CurrentHeight >= block.Height {
		//如果当前区块链高度离此区块高度相差不为1或者已经打包这个节点，直接返回，不对其他值进行校验
		return
	}
	block.Validate()
	block.Recover()
}
