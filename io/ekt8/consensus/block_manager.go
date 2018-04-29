package consensus

import "github.com/EducationEKT/EKT/io/ekt8/blockchain"

type BlockManager struct {
	Blockchain blockchain.BlockChain
}

// TODO delte
//func (m *BlockManager) BlockFromPeer(block *blockchain.Block) {
//	//从peer过来的block先对hash进行校验
//	err := block.Validate(nil)
//	if err != nil {
//		return
//	}
//	//恢复blockBody
//	err = block.Recover()
//	if err != nil {
//		return
//	}
//	m.Blockchain.Locker.RLock()
//	height := m.Blockchain.CurrentHeight
//	m.Blockchain.Locker.Unlock()
//	//对block的高度进行校验
//	//如果当前区块链高度离此区块高度相差不为1或者已经打包这个节点，直接返回，不对其他值进行校验
//	if height+1 != block.Height {
//		return
//	}
//}

func (m *BlockManager) VoteFromPeer(vote blockchain.VoteResults) {

}
