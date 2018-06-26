package blockchain

import (
	"bytes"
	"time"
)

type BlockPolice struct {
	// 记录从其他节点过来的block
	PeerBlocks map[int64][]*Block
}

func NewBlockPolice() BlockPolice {
	return BlockPolice{
		PeerBlocks: make(map[int64][]*Block),
	}
}

// 返回-1表示同一个节点打包了不同的区块
// 返回 0表示已经记录过此区块
// 返回 1表示未记录过此区块，需要进行投票
func (police BlockPolice) BlockFromPeer(block Block, interval time.Duration) int {
	blocks, exist := police.PeerBlocks[block.Height]
	if exist {
		for _, _block := range blocks {
			if !bytes.Equal(block.CurrentHash, _block.CurrentHash) && block.GetRound().Equal(_block.GetRound()) {
				t1, t2 := block.Timestamp, _block.Timestamp
				time := t2 - t1
				if t1 > t2 {
					time = 0 - time
				}
				if time > int64(interval/1e6) {
					continue
				}
				return -1
			}
		}
	}
	blocks = append(blocks, &block)
	police.PeerBlocks[block.Height] = blocks
	return 0
}

// 通过上面的方法获取发现作恶节点之后，可以根据这个方法成成EvilBlock，发送给其他节点
func (police BlockPolice) GetEvilBlock(block Block) *EvilBlock {
	if blocks, exist := police.PeerBlocks[block.Height]; exist {
		blocks = append(blocks, &block)
		police.PeerBlocks[block.Height] = blocks
		for _, _block := range blocks {
			if !bytes.Equal(block.CurrentHash, _block.CurrentHash) && block.GetRound().Equal(_block.GetRound()) {
				return NewEvilBlock(&block, _block)
			}
		}
	}
	return nil
}
