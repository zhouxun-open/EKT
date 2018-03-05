package blockchain

import "fmt"

type Block struct {
	Height       int64
	StatRoot     []byte
	TxRoot       []byte
	Nonce        int64
	PreviousHash []byte
	Hash         []byte
}

func (block *Block) String() string {
	return fmt.Sprintf(`{"height": %d, "statRoot": "%s", "txRoot": "%s", "nonce": %d, "previousHash": "%s"}`,
		block.Height, block.StatRoot, block.TxRoot)
}
