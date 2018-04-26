package blockchain

import "encoding/hex"

var BlockRecorder *BlockRecord

func init() {
	BlockRecorder = &BlockRecord{
		Blocks:     make(map[string]*Block),
		Signatures: make(map[string]string),
	}
}

type BlockRecord struct {
	Blocks     map[string]*Block
	Signatures map[string]string
}

func GetBlockRecordInst() *BlockRecord {
	return BlockRecorder
}

func (record BlockRecord) Record(block *Block, sign []byte) {
	record.Blocks[hex.EncodeToString(block.Hash())] = block
	record.Signatures[hex.EncodeToString(block.Hash())] = hex.EncodeToString(sign)
}
