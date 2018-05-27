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
	"github.com/EducationEKT/EKT/io/ekt8/conf"
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"github.com/EducationEKT/EKT/io/ekt8/crypto"
	"github.com/EducationEKT/EKT/io/ekt8/db"
	"github.com/EducationEKT/EKT/io/ekt8/event"
	"github.com/EducationEKT/EKT/io/ekt8/i_consensus"
)

var currentBlock *Block = nil

type Block struct {
	Height       int64              `json:"height"`
	Timestamp    int64              `json:"timestamp"`
	Nonce        int64              `json:"nonce"`
	Fee          int64              `json:"fee"`
	TotalFee     int64              `json:"totalFee"`
	PreviousHash []byte             `json:"previousHash"`
	CurrentHash  []byte             `json:"currentHash"`
	Signature    string             `json:"signature"`
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

func (block *Block) Data() []byte {
	round := ""
	if block.Height > 0 {
		round = block.Round.String()
	}
	return []byte(fmt.Sprintf(
		`{"height": %d, "timestamp": %d, "nonce": %d, "fee": %d, "totalFee": %d, "previousHash": "%s", "body": "%s", "round": %s, "statRoot": "%s", "txRoot": "%s", "eventRoot": "%s", "tokenRoot": "%s"}`,
		block.Height, block.Timestamp, block.Nonce, block.Fee, block.TotalFee, hex.EncodeToString(block.PreviousHash), hex.EncodeToString(block.Body),
		round, hex.EncodeToString(block.StatRoot), hex.EncodeToString(block.TxRoot),
		hex.EncodeToString(block.EventRoot), hex.EncodeToString(block.TokenRoot),
	))
}

func (block *Block) Hash() []byte {
	return block.CurrentHash
}

func (block *Block) CaculateHash() []byte {
	block.CurrentHash = crypto.Sha3_256(block.Data())
	return block.CurrentHash
}

func (block *Block) NewNonce() {
	block.Nonce++
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
		err := block.StatTree.MustInsert(account.Address(), value)
		if err != nil {
			return false
		}
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
	fromAddress, _ := hex.DecodeString(tx.From)
	toAddress, _ := hex.DecodeString(tx.To)
	account, _ := block.GetAccount(fromAddress)
	recieverAccount, _ := block.GetAccount(toAddress)
	var txResult *common.TxResult
	if fee < block.Fee {
		return common.NewTransactionResult(tx, fee, false, "fee is too less")
	}
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
		block.StatTree.Lock.RUnlock()
	}
	if block.TxTree != nil {
		block.TxTree.Lock.RLock()
		block.TxRoot = block.TxTree.Root
		block.TxTree.Lock.RUnlock()
	}
	if block.EventTree != nil {
		block.EventTree.Lock.RLock()
		block.EventRoot = block.EventTree.Root
		block.EventTree.Lock.RUnlock()
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
	block.EventTree = MPTPlus.MTP_Tree(db.GetDBInst(), block.EventRoot)
	block.StatTree = MPTPlus.MTP_Tree(db.GetDBInst(), block.StatRoot)
	block.TxTree = MPTPlus.MTP_Tree(db.GetDBInst(), block.TxRoot)
	block.Locker = sync.RWMutex{}
	return &block, nil
}

func NewBlock(last *Block, newRound *i_consensus.Round) *Block {
	block := &Block{
		Height:       last.Height + 1,
		Nonce:        0,
		Fee:          last.Fee,
		TotalFee:     0,
		PreviousHash: last.Hash(),
		Timestamp:    time.Now().UnixNano() / 1e6,
		CurrentHash:  nil,
		BlockBody:    NewBlockBody(last.Height + 1),
		Body:         nil,
		Round:        newRound,
		Locker:       sync.RWMutex{},
		StatTree:     MPTPlus.MTP_Tree(db.GetDBInst(), last.StatRoot),
		TxTree:       MPTPlus.NewMTP(db.GetDBInst()),
		EventTree:    MPTPlus.NewMTP(db.GetDBInst()),
		TokenTree:    MPTPlus.MTP_Tree(db.GetDBInst(), last.TokenRoot),
	}
	return block
}

func (block *Block) ValidateNextBlock(next Block, interval time.Duration) bool {
	// 如果不是当前的块的下一个区块，则返回false
	if !bytes.Equal(next.PreviousHash, block.Hash()) || block.Height+1 != next.Height {
		fmt.Printf("This block's previous hash is unexpected, want %s, get %s. \n", hex.EncodeToString(block.Hash()), hex.EncodeToString(next.PreviousHash))
		return false
	}
	return block.ValidateBlockStat(next)
}

func (block *Block) ValidateBlockStat(next Block) bool {
	fmt.Println("Validating block stat merkler proof.")
	// 从打包节点获取body
	body, err := next.Round.Peers[next.Round.CurrentIndex].GetDBValue(next.Body)
	if err != nil {
		fmt.Println("Can not get body from mining node, return false.")
		return false
	}
	next.BlockBody, err = FromBytes(body)
	if err != nil {
		fmt.Println("Get an error body, return false.")
		return false
	}
	//根据上一个区块头生成一个新的区块
	_next := NewBlock(block, next.Round)
	_next.Round = next.Round
	//让新生成的区块执行peer传过来的body中的events进行计算
	for _, eventResult := range next.BlockBody.EventResults {
		evtId, _ := hex.DecodeString(eventResult.EventId)
		evt := event.GetEvent(evtId)
		if evt == nil {
			data, err := next.Round.Peers[next.Round.CurrentIndex].GetDBValue(evtId)
			if err != nil {
				fmt.Println("Can not get this event, validate false.")
				return false
			}
			evt = event.FromBytes(data)
			if evt == nil {
				fmt.Println("Can not get this event, validate false.")
				return false
			}
		}
		_next.HandlerEvent(evt)
	}

	//让新生成的区块执行peer传过来的body中的transactions进行计算
	for _, txResult := range next.BlockBody.TxResults {
		txId, _ := hex.DecodeString(txResult.TxId)
		tx := common.GetTransaction(txId)
		if tx == nil {
			data, err := next.Round.Peers[next.Round.CurrentIndex].GetDBValue(txId)
			if err != nil {
				fmt.Println("Can not get this transaction, validate false.")
				return false
			}
			tx = common.FromBytes(data)
			if tx == nil {
				fmt.Println("Can not get this transaction, validate false.")
				return false
			}
		}
		_next.NewTransaction(tx, block.Fee)
	}
	_next.UpdateMPTPlusRoot()
	if !bytes.Equal(next.TxRoot, _next.TxRoot) ||
		!bytes.Equal(next.EventRoot, _next.EventRoot) ||
		!bytes.Equal(next.StatRoot, _next.StatRoot) ||
		!bytes.Equal(next.TokenRoot, _next.TokenRoot) {
		fmt.Printf("next.Data  = %s, \n_next.Data = %s", next.Data(), block.Data())
		fmt.Printf("next.Hash  = %s, \n_next.Hash = %s \n", hex.EncodeToString(next.Hash()), hex.EncodeToString(_next.CaculateHash()))
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

func (block *Block) Sign() error {
	Signature, err := crypto.Crypto(crypto.Sha3_256(block.Hash()), conf.EKTConfig.PrivateKey)
	block.Signature = hex.EncodeToString(Signature)
	return err
}

// 校验区块头的hash值和其他字段是否匹配，以及签名是否正确
func (block Block) Validate() error {
	if !bytes.Equal(block.CurrentHash, block.CaculateHash()) {
		return errors.New("Invalid Hash")
	}
	sign, err := hex.DecodeString(block.Signature)
	if err != nil {
		return err
	}
	if pubkey, err := crypto.RecoverPubKey(crypto.Sha3_256(block.CurrentHash), sign); err != nil {
		fmt.Println("Recover public key failed.", err)
		return err
	} else {
		if !strings.EqualFold(hex.EncodeToString(crypto.Sha3_256(pubkey)), block.Round.Peers[block.Round.CurrentIndex].PeerId) {
			return errors.New("Invalid signature")
		}
	}
	return nil
}
