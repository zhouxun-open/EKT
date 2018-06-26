package pool

import (
	"fmt"
	"github.com/EducationEKT/EKT/core/common"
	"sort"
	"testing"
)

var pool = Pool{txReady: make(map[string]*common.Transaction), txBlock: make(map[string]UserTransactions)}

var txarr = [10]common.Transaction{
	common.Transaction{From: "bob", To: "alice", TimeStamp: 001, Amount: 99, Nonce: 01, Sign: "bob"},
	common.Transaction{From: "bob", To: "alice", TimeStamp: 002, Amount: 99, Nonce: 02, Sign: "bob"},
	common.Transaction{From: "bob", To: "alice", TimeStamp: 003, Amount: 99, Nonce: 03, Sign: "bob"},
	common.Transaction{From: "bob", To: "alice", TimeStamp: 004, Amount: 99, Nonce: 04, Sign: "bob"},
	common.Transaction{From: "bob", To: "alice", TimeStamp: 005, Amount: 99, Nonce: 05, Sign: "bob"},
	common.Transaction{From: "bob", To: "alice", TimeStamp: 006, Amount: 99, Nonce: 06, Sign: "bob"},
	common.Transaction{From: "bob", To: "alice", TimeStamp: 007, Amount: 99, Nonce: 07, Sign: "bob"},
	common.Transaction{From: "bob", To: "alice", TimeStamp: 017, Amount: 99, Nonce: 10, Sign: "bob"},
	common.Transaction{From: "bob", To: "alice", TimeStamp: 027, Amount: 99, Nonce: 11, Sign: "bob"},
	common.Transaction{From: "bob", To: "alice", TimeStamp: 037, Amount: 99, Nonce: 12, Sign: "bob"},
}

func TestTxPool_Fetch(t *testing.T) {
	//println(pool.Fetch(0))
	pool.Fetch(1)
	size := len(pool.txReady)
	pool.Fetch(size)
	pool.Fetch(size + 1)
}

func BenchmarkTxPool_Fetch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pool.Fetch(1)
	}
}

func TestUserTransactions_sort(t *testing.T) {
	fmt.Println("----TestUserTransactions_sort----")
	t1 := common.Transaction{Nonce: 1}
	t2 := common.Transaction{Nonce: 2}
	t3 := common.Transaction{Nonce: 3}
	u := UserTransactions{}
	u = append(u, &t3, &t2, &t1)
	fmt.Println("before sort")
	fmt.Println(u)
	fmt.Println("after sort")
	sort.Sort(u)
	fmt.Println(u)
	fmt.Println("----TestUserTransactions_sort----")
}

func TestTxPool_Notify(t *testing.T) {
	fmt.Println("----TestTxPool_Notify----")
	pool.ParkTx(&txarr[0], 1) //txReady
	pool.ParkTx(&txarr[1], 0) //txBlock
	pool.ParkTx(&txarr[2], 0) //txBlock
	pool.Notify(&txarr[0])
	pool.Notify(&txarr[1])
	k, e := pool.txReady[txarr[2].TransactionId()]
	if e == false {
		fmt.Println(e)
		t.Fatal()
	}
	fmt.Println(k)
	fmt.Println("----TestTxPool_Notify----")
}
