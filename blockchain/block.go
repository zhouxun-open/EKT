package blockchain

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/EducationEKT/EKT/MPTPlus"
	"github.com/EducationEKT/EKT/core/types"
	"github.com/EducationEKT/EKT/core/userevent"
	"github.com/EducationEKT/EKT/crypto"
	"github.com/EducationEKT/EKT/db"
	"github.com/EducationEKT/EKT/log"
	"github.com/EducationEKT/EKT/round"
)

var currentBlock *Block = nil

type Block struct {
	Height       int64          `json:"height"`
	Timestamp    int64          `json:"timestamp"`
	Nonce        int64          `json:"nonce"`
	Fee          int64          `json:"fee"`
	TotalFee     int64          `json:"totalFee"`
	PreviousHash types.HexBytes `json:"previousHash"`
	CurrentHash  types.HexBytes `json:"currentHash"`
	Signature    types.HexBytes `json:"signature"`
	BlockBody    *BlockBody     `json:"-"`
	Body         types.HexBytes `json:"body"`
	Round        *round.Round   `json:"round"`
	Locker       sync.RWMutex   `json:"-"`
	StatTree     *MPTPlus.MTP   `json:"-"`
	StatRoot     types.HexBytes `json:"statRoot"`
	TxTree       *MPTPlus.MTP   `json:"-"`
	TxRoot       types.HexBytes `json:"txRoot"`
	TokenTree    *MPTPlus.MTP   `json:"-"`
	TokenRoot    types.HexBytes `json:"tokenRoot"`
}

func (block Block) GetRound() *round.Round {
	return block.Round.Clone()
}

func (block *Block) Bytes() []byte {
	block.UpdateMPTPlusRoot()
	data, _ := json.Marshal(block)
	return data
}

func (block *Block) Data() []byte {
	round := ""
	if block.Height > 0 {
		round = block.GetRound().String()
	}
	return []byte(fmt.Sprintf(
		`{"height": %d, "timestamp": %d, "nonce": %d, "fee": %d, "totalFee": %d, "previousHash": "%s", "body": "%s", "round": %s, "statRoot": "%s", "txRoot": "%s", "eventRoot": "%s", "tokenRoot": "%s"}`,
		block.Height, block.Timestamp, block.Nonce, block.Fee, block.TotalFee, hex.EncodeToString(block.PreviousHash), hex.EncodeToString(block.Body),
		round, hex.EncodeToString(block.StatRoot), hex.EncodeToString(block.TxRoot), hex.EncodeToString(block.TokenRoot),
	))
}

func (block *Block) Hash() []byte {
	return block.CurrentHash
}

func (block *Block) CaculateHash() []byte {
	block.CurrentHash = crypto.Sha3_256(block.Data())
	return block.CurrentHash
}

func (block Block) ValidateHash() bool {
	return bytes.EqualFold(block.CurrentHash, block.CaculateHash())
}

func (block *Block) NewNonce() {
	block.Nonce++
}

func (block Block) GetAccount(address []byte) (*types.Account, error) {
	if block.StatTree == nil {
		block.StatTree = MPTPlus.MTP_Tree(db.GetDBInst(), block.StatRoot)
	}
	value, err := block.StatTree.GetValue(address)
	if err != nil {
		return nil, err
	}
	var account types.Account
	err = json.Unmarshal(value, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (block Block) ExistAddress(address []byte) bool {
	if block.StatTree == nil {
		block.StatTree = MPTPlus.MTP_Tree(db.GetDBInst(), block.StatRoot)
	}
	return block.StatTree.ContainsKey(address)
}

func (block *Block) CreateGenesisAccount(account types.Account) bool {
	err := block.StatTree.MustInsert(account.Address, account.ToBytes())
	if err != nil {
		return false
	}
	block.UpdateMPTPlusRoot()
	return true
}

func (block *Block) NewTransaction(tx userevent.Transaction, fee int64) *userevent.TxResult {
	account, _ := block.GetAccount(tx.GetFrom())
	recieverAccount, err := block.GetAccount(tx.GetTo())
	if err != nil || nil == recieverAccount {
		*recieverAccount = types.CreateAccount(tx.GetTo(), 0)
	}
	var txResult *userevent.TxResult

	// 如果fee太少，默认使用系统最少值
	if fee < block.Fee {
		fee = block.Fee
	}

	if tx.Nonce != account.Nonce+1 {
		txResult = userevent.NewTransactionResult(tx, fee, false, "invalid nonce")
	} else if tx.TokenAddress == "" {
		if account.GetAmount() < tx.Amount+fee {
			txResult = userevent.NewTransactionResult(tx, fee, false, "no enough gas")
		} else {
			block.TotalFee++
			account.ReduceAmount(tx.Amount + fee)
			recieverAccount.AddAmount(tx.Amount)
			block.StatTree.MustInsert(tx.GetFrom(), account.ToBytes())
			block.StatTree.MustInsert(tx.GetTo(), recieverAccount.ToBytes())
			txResult = userevent.NewTransactionResult(tx, fee, true, "")
		}
	} else {
		if account.Balances[tx.TokenAddress] < tx.Amount {
			txResult = userevent.NewTransactionResult(tx, fee, false, "no enough amount")
		} else if account.GetAmount() < fee {
			txResult = userevent.NewTransactionResult(tx, fee, false, "no enough gas")
		} else {
			block.TotalFee++
			account.Balances[tx.TokenAddress] -= tx.Amount
			account.ReduceAmount(fee)
			if recieverAccount.Balances == nil {
				recieverAccount.Balances = make(map[string]int64)
				recieverAccount.Balances[tx.TokenAddress] = 0
			}
			recieverAccount.Balances[tx.TokenAddress] += tx.Amount
			block.StatTree.MustInsert(tx.GetFrom(), account.ToBytes())
			block.StatTree.MustInsert(tx.GetTo(), recieverAccount.ToBytes())
			txResult = userevent.NewTransactionResult(tx, fee, true, "")
		}
	}
	txId, _ := hex.DecodeString(tx.TransactionId())

	block.TxTree.MustInsert(txId, txResult.ToBytes())
	return txResult
}

func (block *Block) UpdateMPTPlusRoot() {
	if block.StatTree != nil {
		block.StatTree.Lock.RLock()
		block.StatRoot = block.StatTree.Root
		block.StatTree.Lock.RUnlock()
	}
	if block.TxTree != nil {
		block.TxTree.Lock.RLock()
		block.TxRoot = block.TxTree.Root
		block.TxTree.Lock.RUnlock()
	}
	if block.TokenTree != nil {
		block.TokenTree.Lock.RLock()
		block.TokenRoot = block.TokenTree.Root
		block.TokenTree.Lock.RUnlock()
	}
}

func FromBytes2Block(data []byte) (*Block, error) {
	var block Block
	err := json.Unmarshal(data, &block)
	if err != nil {
		return nil, err
	}
	block.StatTree = MPTPlus.MTP_Tree(db.GetDBInst(), block.StatRoot)
	block.TxTree = MPTPlus.MTP_Tree(db.GetDBInst(), block.TxRoot)
	block.Locker = sync.RWMutex{}
	return &block, nil
}

func NewBlock(last Block) *Block {
	block := &Block{
		Height:       last.Height + 1,
		Nonce:        0,
		Fee:          last.Fee,
		TotalFee:     0,
		PreviousHash: last.Hash(),
		Timestamp:    time.Now().UnixNano() / 1e6,
		CurrentHash:  nil,
		BlockBody:    NewBlockBody(),
		Body:         nil,
		Locker:       sync.RWMutex{},
		StatTree:     MPTPlus.MTP_Tree(db.GetDBInst(), last.StatRoot),
		TxTree:       MPTPlus.NewMTP(db.GetDBInst()),
		TokenTree:    MPTPlus.MTP_Tree(db.GetDBInst(), last.TokenRoot),
	}
	return block
}

func (block Block) ValidateNextBlock(next Block, events []userevent.IUserEvent) bool {
	// 如果不是当前的块的下一个区块，则返回false
	if !bytes.Equal(next.PreviousHash, block.Hash()) || block.Height+1 != next.Height {
		return false
	}
	return block.ValidateBlockStat(next, events)
}

// consensus 模块调用这个函数，获得一个block对象之后发送给其他节点，其他节点同意之后调用上面的NewBlock方法
func (block *Block) Pack(difficulty []byte) {
	block.Locker.Lock()
	defer block.Locker.Unlock()
	start := time.Now().Nanosecond()
	log.Info("Caculating block hash.")
	for ; !bytes.HasPrefix(block.CaculateHash(), difficulty); block.NewNonce() {
	}
	end := time.Now().Nanosecond()
	log.Info("Caculated block hash, cost %d ms. \n", (end-start+1e9)%1e9/1e6)
}

func (block Block) ValidateBlockStat(next Block, events []userevent.IUserEvent) bool {
	BlockRecorder.SetBlock(&next)
	log.Info("Validating block stat merkler proof.")

	//根据上一个区块头生成一个新的区块
	_next := NewBlock(block)

	//让新生成的区块执行peer传过来的body中的user events进行计算
	if len(block.BlockBody.Events) > 0 {
		//for _, eventId := range block.BlockBody.Events {
		//}
	}
	//if block.BlockBody != nil {
	//	block.BlockBody.Events.Range(func(key, value interface{}) bool {
	//		_, ok1 := key.(string)
	//		list, ok2 := value.([]string)
	//		if ok1 && ok2 && len(list) > 0 {
	//			for _, eventId := range list {
	//				txId, err := hex.DecodeString(eventId)
	//				if err != nil {
	//					return false
	//				}
	//				tx := types.GetTransaction(txId)
	//				_next.NewTransaction(*tx, tx.Fee)
	//			}
	//		}
	//		return true
	//	})
	//}

	// 更新默克尔树根
	_next.UpdateMPTPlusRoot()

	// 判断默克尔根是否相同
	if !bytes.Equal(next.TxRoot, _next.TxRoot) ||
		!bytes.Equal(next.StatRoot, _next.StatRoot) ||
		!bytes.Equal(next.TokenRoot, _next.TokenRoot) {
		log.Info("next.Data  = %s, \n_next.Data = %s", next.Data(), block.Data())
		log.Info("next.Hash  = %s, \n_next.Hash = %s \n", hex.EncodeToString(next.Hash()), hex.EncodeToString(_next.CaculateHash()))
		return false
	}

	BlockRecorder.SetStatus(hex.EncodeToString(next.CurrentHash), 100)
	return true
}

func (block *Block) Sign(privKey []byte) error {
	Signature, err := crypto.Crypto(crypto.Sha3_256(block.CurrentHash), privKey)
	block.Signature = Signature
	return err
}
