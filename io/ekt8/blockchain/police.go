package blockchain

import "bytes"

type BlockPolice struct {
	// 记录从其他节点过来的block
	PeerBlocks map[int64][]*Block
}

// 返回-1表示同一个节点打包了不同的区块
// 返回 0表示已经记录过此区块
// 返回 1表示未记录过此区块，需要进行投票
func (police BlockPolice) BlockFromPeer(block *Block) int {
	if blocks, exist := police.PeerBlocks[block.Height]; exist {
		for _, _block := range blocks {
			if !bytes.Equal(block.CurrentHash, _block.CurrentHash) && block.Round.Equal(_block.Round) {
				return -1
			}
		}
	} else {
		blocks = append(blocks, block)
		return 1
	}
	return 0
}
