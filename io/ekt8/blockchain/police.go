package blockchain

type BlockPolice struct {
	// 记录从其他节点过来的block
	BlockFromPeer map[int]*Block
	//
	Vote map[string]bool
}
