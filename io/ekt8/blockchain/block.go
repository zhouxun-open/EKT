package blockchain

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/EducationEKT/EKT/io/ekt8/MPTPlus"
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"github.com/EducationEKT/EKT/io/ekt8/crypto"
	"github.com/EducationEKT/EKT/io/ekt8/db"
	"github.com/EducationEKT/EKT/io/ekt8/event"
	"github.com/EducationEKT/EKT/io/ekt8/i_consensus"
)

var currentBlock *Block = nil

type Block struct {
	Height       int64              `json:"height"`
	Timestamp    int                `json:"timestamp"`
	Nonce        int64              `json:"nonce"`
	Fee          int64              `json:"fee"`
	TotalFee     int64              `json:"totalFee"`
	PreviousHash []byte             `json:"previousHash"`
	CurrentHash  []byte             `json:"currentHash"`
	BlockBody    *BlockBody         `json:"-"`
	Body         []byte             `json:"body"`
	Round        *i_consensus.Round `json:"round"`
	Locker       sync.RWMutex       `json:"-"`
	StatTree     *MPTPlus.MTP       `json:"-"`
	StatRoot     []byte             `json:"statRoot"`
	TxTree       *MPTPlus.MTP       `json:"-"`
	TxRoot       []byte             `json:"txRoot"`
	EventTree    *MPTPlus.MTP       `json:"-"`
	EventRoot    []byte             `json:"eventRoot"`
	TokenTree    *MPTPlus.MTP       `json:"-"`
	TokenRoot    []byte             `json:"tokenRoot"`
}

func (block *Block) Bytes() []byte {
	block.UpdateMPTPlusRoot()
	data, _ := json.Marshal(block)
	return data
}

func (block *Block) Hash() []byte {
	return block.CurrentHash
}

func (block *Block) CaculateHash() []byte {
	block.CurrentHash = crypto.Sha3_256(block.Bytes())
	return block.CurrentHash
}

func (block *Block) NewNonce() {
	block.Nonce++
}

// 校验区块头的hash值和其他字段是否匹配，以及签名是否正确
func (block Block) Validate(sign []byte) error {
	if !bytes.Equal(block.CurrentHash, block.CaculateHash()) {
		return errors.New("Invalid Hash")
	}
	pub, err := crypto.RecoverPubKey(block.Hash(), sign)
	if err != nil {
		return err
	}
	if !strings.EqualFold(hex.EncodeToString(crypto.Sha3_256(pub)), block.Round.Peers[block.Round.CurrentIndex].PeerId) {
		return errors.New("Invalid signature")
	}
	return nil
}

// 从网络节点过来的区块头，如果区块的body为空，则从打包节点获取
// 获取之后会对blockBody的Hash进行校验，如果不符合要求则放弃Recover
func (block Block) Recover() error {
	if !bytes.Equal(block.Body, block.BlockBody.Bytes()) {
		peer := block.Round.Peers[block.Round.CurrentIndex]
		bodyData, err := peer.GetDBValue(block.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(bodyData, block.BlockBody)
		if err != nil {
			return err
		}
		if !bytes.Equal(crypto.Sha3_256(block.BlockBody.Bytes()), block.Body) {
			return errors.New(fmt.Sprintf("Block body is wrong, want hash(body) = %s, get %s", block.Body, crypto.Sha3_256(block.BlockBody.Bytes())))
		}
	}
	return nil
}

func (block *Block) GetAccount(address []byte) (*common.Account, error) {
	value, err := block.StatTree.GetValue(address)
	if err != nil {
		return nil, err
	}
	var account common.Account
	err = json.Unmarshal(value, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (block *Block) ExistAddress(address []byte) bool {
	return block.StatTree.ContainsKey(address)
}

func (block *Block) CreateAccount(address, pubKey []byte) {
	if !block.ExistAddress(address) {
		block.newAccount(address, pubKey)
	}
}

func (block *Block) InsertAccount(account common.Account) bool {
	if !block.ExistAddress(account.Address()) {
		value, _ := json.Marshal(account)
		block.StatTree.MustInsert(account.Address(), value)
		block.UpdateMPTPlusRoot()
		return true
	}
	return false
}

func (block *Block) newAccount(address []byte, pubKey []byte) {
	account := common.NewAccount(address, pubKey)
	value, _ := json.Marshal(account)
	block.StatTree.MustInsert(address, value)
	block.UpdateMPTPlusRoot()
}

func (block *Block) NewTransaction(tx *common.Transaction, fee int64) *common.TxResult {
	block.Locker.Lock()
	defer block.Locker.Unlock()
	fromAddress, _ := hex.DecodeString(tx.From)
	toAddress, _ := hex.DecodeString(tx.To)
	account, _ := block.GetAccount(fromAddress)
	recieverAccount, _ := block.GetAccount(toAddress)
	var txResult *common.TxResult
	if account.GetAmount() < tx.Amount+fee {
		txResult = common.NewTransactionResult(tx, fee, false, "no enough amount")
	} else {
		txResult = common.NewTransactionResult(tx, fee, true, "")
		account.ReduceAmount(tx.Amount + block.Fee)
		block.TotalFee += block.Fee
		recieverAccount.AddAmount(tx.Amount)
		block.StatTree.MustInsert(fromAddress, account.ToBytes())
		block.StatTree.MustInsert(toAddress, recieverAccount.ToBytes())
	}
	txId, _ := hex.DecodeString(tx.TransactionId())
	block.TxTree.MustInsert(txId, txResult.ToBytes())
	block.UpdateMPTPlusRoot()
	return txResult
}

func (block *Block) UpdateMPTPlusRoot() {
	if block.StatTree != nil {
		block.StatTree.Lock.RLock()
		block.StatRoot = block.StatTree.Root
	}
	if block.TxTree != nil {
		block.TxTree.Lock.RLock()
		block.TxRoot = block.TxTree.Root
	}
	if block.EventTree != nil {
		block.EventTree.Lock.RLock()
		block.EventRoot = block.EventTree.Root
	}
	if block.TokenTree != nil {
		block.TokenTree.Lock.RLock()
		block.TokenRoot = block.TokenTree.Root
	}
}

func FromBytes2Block(data []byte) (*Block, error) {
	var block Block
	err := json.Unmarshal(data, block)
	if err != nil {
		return nil, err
	}
	block.EventTree = MPTPlus.MTP_Tree(db.GetDBInst(), block.EventRoot)
	block.StatTree = MPTPlus.MTP_Tree(db.GetDBInst(), block.StatRoot)
	block.TxTree = MPTPlus.MTP_Tree(db.GetDBInst(), block.TxRoot)
	block.Locker = sync.RWMutex{}
	return &block, nil
}

func NewBlock(last *Block) *Block {
	block := &Block{
		Height:       last.Height + 1,
		Nonce:        0,
		Fee:          last.Fee,
		TotalFee:     0,
		PreviousHash: last.Hash(),
		Timestamp:    time.Now().Second(),
		CurrentHash:  nil,
		BlockBody:    NewBlockBody(last.Height + 1),
		Body:         nil,
		Round:        last.Round.NextRound(last.Hash()),
		Locker:       sync.RWMutex{},
		StatTree:     last.StatTree,
		TxTree:       MPTPlus.NewMTP(db.GetDBInst()),
		EventTree:    MPTPlus.NewMTP(db.GetDBInst()),
		TokenTree:    last.TokenTree,
	}
	return block
}

func (block *Block) ValidateNextBlock(next *Block, interval int) bool {
	//如果不是当前的块的下一个区块，则返回false
	if !bytes.Equal(next.PreviousHash, block.Hash()) || block.Height+1 != next.Height {
		return false
	}
	time := next.Timestamp - block.Timestamp
	//时间差在下一个区块，说明中间没有错过区块
	// 如果前n个节点没有出块，判断当前节点是否拥有打包权限（时间）
	n := time / interval
	if n > len(block.Round.Peers) {
		// 如果已经超过一轮没有出块，则所有节点等放弃出块，等待当前轮下一个节点进行打包
		if !block.Round.IndexPlus(block.Hash()).Equal(next.Round) {
			return false
		}
	}
	remainder := time % interval
	if remainder > interval/2 {
		n++
	}
	// 需要计算下一个区块的index
	if block.Round.CurrentIndex+n >= len(block.Round.Peers) {
		// 计算当前区块的区块差
		miningNumber := len(block.Round.Peers) - block.Round.CurrentIndex + next.Round.CurrentIndex
		if miningNumber != n {
			return false
		}
	} else if block.Round.CurrentIndex+n != next.Round.CurrentIndex {
		return false
	}
	return block.ValidateBlockStat(next)
}

func (block *Block) ValidateBlockStat(next *Block) bool {
	// 从打包节点获取body
	body, err := next.Round.Peers[next.Round.CurrentIndex].GetDBValue(next.Body)
	if err != nil {
		return false
	}
	next.BlockBody, err = FromBytes(body)
	if err != nil {
		return false
	}

	//根据上一个区块头生成一个新的区块
	_next := NewBlock(block)

	//让新生成的区块执行peer传过来的body中的events进行计算
	for _, eventResult := range next.BlockBody.EventResults {
		evtId, _ := hex.DecodeString(eventResult.EventId)
		evt := event.GetEvent(evtId)
		if evt == nil {
			data, err := next.Round.Peers[next.Round.CurrentIndex].GetDBValue(evtId)
			if err != nil {
				return false
			}
			evt = event.FromBytes(data)
			if evt == nil {
				return false
			}
		}
		_next.HandlerEvent(evt)
	}

	//让新生成的区块执行peer传过来的body中的transactions进行计算
	for _, txResult := range block.BlockBody.TxResults {
		txId, _ := hex.DecodeString(txResult.TxId)
		tx := common.GetTransaction(txId)
		if tx == nil {
			data, err := next.Round.Peers[next.Round.CurrentIndex].GetDBValue(txId)
			if err != nil {
				return false
			}
			tx = common.FromBytes(data)
			if tx == nil {
				return false
			}
		}
		_next.NewTransaction(tx, block.Fee)
	}
	_next.Nonce = next.Nonce
	if !bytes.Equal(next.Hash(), _next.CaculateHash()) {
		return false
	}

	return true
}

func (block *Block) HandlerEvent(evt *event.Event) event.EventResult {
	evtResult := event.EventResult{
		EventId: hex.EncodeToString(evt.EventId()),
		Success: false,
		Reason:  "",
	}
	if evt.EventType == event.NewAccountEvent {
		param := evt.EventParam.(event.NewAccountParam)
		address, _ := hex.DecodeString(param.Address)
		pubKey, _ := hex.DecodeString(param.PubKey)
		if !block.ExistAddress(address) {
			block.newAccount(address, pubKey)
			evtResult.Success = true
		} else {
			evtResult.Reason = "AddressExist"
		}
	}
	block.EventTree.MustInsert(evt.EventId(), evtResult.Bytes())
	return evtResult
}
