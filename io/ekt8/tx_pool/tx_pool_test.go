package tx_pool

import (
	"testing"
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"sort"
	"fmt"
)

var txPool=TxPool{ready:make(map[string]*common.Transaction),block:make(map[string]UserTransactions)}

var txarr =[10]common.Transaction{
	common.Transaction{TransactionId:"1",From:"bob",To:"alice",TimeStamp:001,Amount:99,Nonce:01,Sign:"bob"},
	common.Transaction{TransactionId:"2",From:"bob",To:"alice",TimeStamp:002,Amount:99,Nonce:02,Sign:"bob"},
	common.Transaction{TransactionId:"3",From:"bob",To:"alice",TimeStamp:003,Amount:99,Nonce:03,Sign:"bob"},
	common.Transaction{TransactionId:"4",From:"bob",To:"alice",TimeStamp:004,Amount:99,Nonce:04,Sign:"bob"},
	common.Transaction{TransactionId:"5",From:"bob",To:"alice",TimeStamp:005,Amount:99,Nonce:05,Sign:"bob"},
	common.Transaction{TransactionId:"6",From:"bob",To:"alice",TimeStamp:006,Amount:99,Nonce:06,Sign:"bob"},
	common.Transaction{TransactionId:"7",From:"bob",To:"alice",TimeStamp:007,Amount:99,Nonce:07,Sign:"bob"},
	common.Transaction{TransactionId:"8",From:"bob",To:"alice",TimeStamp:017,Amount:99,Nonce:10,Sign:"bob"},
	common.Transaction{TransactionId:"9",From:"bob",To:"alice",TimeStamp:027,Amount:99,Nonce:11,Sign:"bob"},
	common.Transaction{TransactionId:"10",From:"bob",To:"alice",TimeStamp:037,Amount:99,Nonce:12,Sign:"bob"},
}


func TestTxPool_Fetch(t *testing.T) {
	//println(txPool.Fetch(0))
	txPool.Fetch(1)
	size := len(txPool.ready)
	txPool.Fetch(size)
	txPool.Fetch(size + 1)
}

func BenchmarkTxPool_Fetch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		txPool.Fetch(1)
	}
}

func TestUserTransactions_sort(t *testing.T) {
	fmt.Println("----TestUserTransactions_sort----")
	t1:=common.Transaction{Nonce:1}
	t2:=common.Transaction{Nonce:2}
	t3:=common.Transaction{Nonce:3}
	u:=UserTransactions{}
	u=append(u,&t3,&t2,&t1)
	fmt.Println("before sort")
	fmt.Println(u)
	fmt.Println("after sort")
	sort.Sort(u)
	fmt.Println(u)
	fmt.Println("----TestUserTransactions_sort----")
}

func TestTxPool_Notify(t *testing.T) {
	fmt.Println("----TestTxPool_Notify----")
	txPool.Park(&txarr[0],1)//ready
	txPool.Park(&txarr[1],0)//block
	txPool.Park(&txarr[2],0)//block
	txPool.Notify(&txarr[0])
	txPool.Notify(&txarr[1])
	k,e:=txPool.ready[txarr[2].TransactionId]
	if e == false{
		fmt.Println(e)
		t.Fatal()
	}
	fmt.Println(k)
	fmt.Println("----TestTxPool_Notify----")
}
