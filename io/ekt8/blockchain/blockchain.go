package blockchain

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"errors"
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"github.com/EducationEKT/EKT/io/ekt8/crypto"
	"github.com/EducationEKT/EKT/io/ekt8/db"
	"github.com/EducationEKT/EKT/io/ekt8/event"
	"github.com/EducationEKT/EKT/io/ekt8/i_consensus"
	"github.com/EducationEKT/EKT/io/ekt8/pool"
	"strings"
)

var BackboneChainId []byte
var EKTTokenId []byte

func init() {
	BackboneChainId, _ = hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000001")
	EKTTokenId, _ = hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000001")
}

const (
	CurrentBlockKey       = "CurrentBlock"
	BackboneConsensus     = i_consensus.DPOS
	BackboneBlockInterval = 3
	InitStatus            = 0
	OpenStatus            = 100
	StartPackStatus       = 110
	CaculateHashStatus    = 150
)

type BlockChain struct {
	ChainId         []byte
	Consensus       i_consensus.ConsensusType
	CurrentBlock    *Block
	CurrentBody     *BlockBody
	Locker          sync.RWMutex
	Status          int // 100 正在计算MTProot, 150停止计算root,开始计算block Hash
	Fee             int64
	Difficulty      []byte
	Pool            *pool.Pool
	CurrentHeight   int64
	LastBlockHeader *Block
}

func (blockchain *BlockChain) PackSignal(cb func(block *Block)) {
	blockchain.Locker.Lock()
	if blockchain.Status != StartPackStatus {
		blockchain.Status = StartPackStatus
		block := blockchain.WaitAndPack()
		db.GetDBInst().Set(block.CaculateHash(), block.Bytes())
		db.GetDBInst().Set(blockchain.GetBlockByHeightKey(block.Height), block.CaculateHash())
		db.GetDBInst().Set(blockchain.GetBlockBodyByHeightKey(block.Height), block.Body)
		db.GetDBInst().Set(blockchain.CurrentBlockKey(), block.CaculateHash())
		blockchain.CurrentBlock = block
		blockchain.CurrentBody = block.BlockBody
		blockchain.CurrentHeight = block.Height
		cb(block)
		blockchain.Locker.Unlock()
	}
}

func (blockchain *BlockChain) GetStatus() int {
	blockchain.Locker.RLock()
	defer blockchain.Locker.RUnlock()
	return blockchain.Status
}

//func (blockchain *BlockChain) StartPack() {
//	blockchain.Locker.Lock()
//	defer blockchain.Locker.Unlock()
//	if blockchain.Status == 100 {
//
//	}
//}

func (blockchain *BlockChain) ValidateBlock(block *Block) bool {
	return false
}

//返回从指定高度到当前区块的区块头
func (blockchain *BlockChain) GetBlockHeaders(fromHeight int64) []*Block {
	headers := make([]*Block, 0)
	var lastHeight int64 = 0
	lastBlock, err := blockchain.LastBlock()
	if err == nil && lastBlock != nil {
		lastHeight = lastBlock.Height
	}
	var block *Block = lastBlock
	for height := lastHeight; height >= fromHeight; height-- {
		if block != nil {
			headers = append(headers, block)
		}
		data, err := db.GetDBInst().Get(block.PreviousHash)
		if err != nil {
			return headers
		}
		var header Block
		err = json.Unmarshal(data, &header)
		if err != nil {
			return headers
		}
		block = &header
	}
	return headers
}

func (blockchain *BlockChain) GetBlockByHeight(height int64) (*Block, error) {
	if height > blockchain.CurrentHeight {
		return nil, errors.New("Invalid height")
	}
	key := blockchain.GetBlockByHeightKey(height)
	data, err := db.GetDBInst().Get(key)
	if err != nil {
		return nil, err
	}
	return FromBytes2Block(data)
}

func (blockchain *BlockChain) GetBlockBodyByHeight(height int64) (*BlockBody, error) {
	if height > blockchain.CurrentHeight {
		return nil, errors.New("Invalid height")
	}
	key := blockchain.GetBlockBodyByHeightKey(height)
	data, err := db.GetDBInst().Get(key)
	if err != nil {
		return nil, err
	}
	var body BlockBody
	err = json.Unmarshal(data, &body)
	return &body, err
}

func (blockchain *BlockChain) GetBlockBodyByHeightKey(height int64) []byte {
	return []byte(fmt.Sprint(`GetBlockBodyByHeight_%s_%d`, hex.EncodeToString(blockchain.ChainId), height))
}

func (blockchain *BlockChain) GetBlockByHeightKey(height int64) []byte {
	return []byte(fmt.Sprint(`GetBlockByHeight_%s_%d`, hex.EncodeToString(blockchain.ChainId), height))
}

// 即将废除
//func (blockchain *BlockChain) NewBlock(block *Block) error {
//	blockchain.Locker.Lock()
//	defer blockchain.Locker.Unlock()
//	if err := block.Validate(); err != nil {
//		return err
//	}
//	db.GetDBInst().Set(block.Hash(), block.Bytes())
//	// TODO sync tx and stat
//	// TODO refact block的产生和交易模块
//	block.UpdateMPTPlusRoot()
//	return db.GetDBInst().Set(blockchain.CurrentBlockKey(), block.Hash())
//}

func (blockchain *BlockChain) SaveBlock(block *Block) {
	block.UpdateMPTPlusRoot()
	fmt.Println(block)
	err := db.GetDBInst().Set(block.CaculateHash(), block.Bytes())
	if err != nil {
		panic(err)
	}
	err = db.GetDBInst().Set(blockchain.CurrentBlockKey(), block.Hash())
	if err != nil {
		panic(err)
	}
}

func (blockchain *BlockChain) LastBlock() (*Block, error) {
	var err error = nil
	var block *Block
	if currentBlock == nil {
		key := blockchain.CurrentBlockKey()
		blockHash, err := db.GetDBInst().Get(key)
		if err != nil {
			return nil, err
		}
		blockValue, err := db.GetDBInst().Get(blockHash)
		if err == nil {
			err = json.Unmarshal(blockValue, &block)
			currentBlock = block
			return block, err
		}
	}
	return currentBlock, err
}

func (blockchain *BlockChain) CurrentBlockKey() []byte {
	buffer := bytes.Buffer{}
	buffer.WriteString(CurrentBlockKey)
	buffer.Write(blockchain.ChainId)
	return buffer.Bytes()
}

func (blockchain *BlockChain) WaitAndPack() *Block {
	eventTimeout := time.After(1 * time.Second)
	block := NewBlock(blockchain.CurrentBlock)
	fmt.Println("Packing transaction and other events.")
	for {
		flag := false
		select {
		case <-eventTimeout:
			flag = true
			break
		default:
			if block.BlockBody.Size() >= 3500*3 {
				flag = true
				break
			}
			evt := blockchain.Pool.FetchEvent()
			if evt != nil {
				if strings.EqualFold(evt.EventType, event.NewAccountEvent) {
					param := evt.EventParam.(event.NewAccountParam)
					address, _ := hex.DecodeString(param.Address)
					pubKey, _ := hex.DecodeString(param.PubKey)
					if block.InsertAccount(*common.NewAccount(address, pubKey)) {
						block.BlockBody.AddEventResult(event.EventResult{Success: true, EventId: evt.EventParam.Id()})
					} else {
						block.BlockBody.AddEventResult(event.EventResult{Success: false, Reason: "address exist", EventId: evt.EventParam.Id()})
					}
				}
				blockchain.Pool.NotifyEvent(*evt)
			} else {
				tx := blockchain.Pool.FetchTx()
				if tx != nil {
					txResult := block.NewTransaction(tx, block.Fee)
					blockchain.Pool.Notify(tx)
					block.BlockBody.AddTxResult(*txResult)
				}
			}
		}
		if flag {
			break
		}
	}
	blockchain.Pack(block)
	return block
}

func (blockchain *BlockChain) NewTransaction(tx *common.Transaction) {
	blockchain.Pool.ParkTx(tx, pool.Ready)
	//blockchain.Locker.Lock()
	//defer blockchain.Locker.Unlock()
	//if blockchain.Status == OpenStatus {
	//	blockchain.CurrentBlock.NewTransaction(tx, blockchain.Fee)
	//} else {
	//	blockchain.Pool.ParkTx(tx, pool.Ready)
	//}
}

// consensus 模块调用这个函数，获得一个block对象之后发送给其他节点，其他节点同意之后调用上面的NewBlock方法
func (blockchain *BlockChain) Pack(block *Block) {
	block.Locker.Lock()
	defer block.Locker.Unlock()
	bodyData, _ := json.Marshal(block.BlockBody)
	block.Body = crypto.Sha3_256(bodyData)
	db.GetDBInst().Set(block.Body, bodyData)
	start := time.Now().Nanosecond()
	fmt.Println(start)
	fmt.Println("Caculating block hash.")
	for ; !bytes.HasPrefix(block.CaculateHash(), blockchain.Difficulty); block.NewNonce() {
	}
	end := time.Now().Nanosecond()
	fmt.Println(end)
	fmt.Printf("Caculated block hash, cost %d ms\n", (end-start+1e9)%1e9/1e6)
}
