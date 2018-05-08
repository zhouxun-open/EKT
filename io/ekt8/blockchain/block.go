package blockchain

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/EducationEKT/EKT/io/ekt8/MPTPlus"
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"github.com/EducationEKT/EKT/io/ekt8/crypto"
	"github.com/EducationEKT/EKT/io/ekt8/db"
	"github.com/EducationEKT/EKT/io/ekt8/event"
	"github.com/EducationEKT/EKT/io/ekt8/i_consensus"
	"github.com/EducationEKT/EKT/io/ekt8/param"
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

// 校验区块头的hash值和其他字段是否匹配，以及签名是否正确
func (block Block) Validate(sign []byte) error {
	if !bytes.Equal(block.CurrentHash, block.CaculateHash()) {
		return errors.New("Invalid Hash")
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

func NewBlock(last *Block) *Block {
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
		Round:        nil,
		Locker:       sync.RWMutex{},
		StatTree:     MPTPlus.MTP_Tree(db.GetDBInst(), last.StatRoot),
		TxTree:       MPTPlus.NewMTP(db.GetDBInst()),
		EventTree:    MPTPlus.NewMTP(db.GetDBInst()),
		TokenTree:    last.TokenTree,
	}
	if last.Height == 0 {
		block.Round = &i_consensus.Round{
			Peers:        param.MainChainDPosNode,
			CurrentIndex: 0,
		}
	} else {
		block.Round = last.Round.NewRandom(last.CurrentHash)
		block.Round.CurrentIndex = block.Round.MyIndex()
	}
	return block
}

func (block *Block) ValidateNextBlock(next Block, interval time.Duration) bool {
	round := &i_consensus.Round{
		Peers:        param.MainChainDPosNode,
		CurrentIndex: -1,
	}
	if block.Height > 0 {
		round = block.Round
	}
	// 如果不是当前的块的下一个区块，则返回false
	if !bytes.Equal(next.PreviousHash, block.Hash()) || block.Height+1 != next.Height {
		fmt.Printf("This block's previous hash is unexpected, want %s, get %s. \n", hex.EncodeToString(block.Hash()), hex.EncodeToString(next.PreviousHash))
		return false
	}
	time := next.Timestamp - block.Timestamp
	fmt.Printf("block.Timestamp=%d, next.Timestamp=%d, time=%d\n", block.Timestamp, next.Timestamp, time)
	// 时间差在下一个区块，说明中间没有错过区块
	// 如果前n个节点没有出块，判断当前节点是否拥有打包权限（时间）
	n := int(int64(time) / int64(interval/1e6))
	if n > len(round.Peers) {
		// 如果已经超过一轮没有出块，则所有节点等放弃出块，等待当前轮下一个节点进行打包
		if !round.IndexPlus(block.Hash()).Equal(next.Round) {
			// 如果当前节点不是下一个节点，返回错误
			fmt.Println("Time goes more than a round interval, and current node is not the next node, validate false.")
			return false
		}
	}
	remainder := int(time) % int(interval/1e6)
	if remainder > int(interval/1e6)/2 {
		n++
	}
	// 需要计算下一个区块的index
	if round.CurrentIndex+n >= len(round.Peers) {
		// 计算当前区块的区块差
		miningNumber := (round.CurrentIndex + n) % round.Len()
		if miningNumber != next.Round.CurrentIndex {
			fmt.Println("Current peer is not the next block candidate, validate false.")
			return false
		}
	} else if round.CurrentIndex+n != next.Round.CurrentIndex {
		fmt.Printf("Round: %s, n:%d, next.Round.CurrentIndex:%d current node is not the next block candidate, validate false. \n", round.String(), n, next.Round.CurrentIndex)
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
	_next := NewBlock(block)
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

func (block *Block) Sign(privKey []byte) string {
	data, _ := crypto.Crypto(crypto.Sha3_256(block.Hash()), privKey)
	return hex.EncodeToString(data)
}
