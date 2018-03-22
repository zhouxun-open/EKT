package tx_pool

import (
	"testing"
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"sort"
	"fmt"
)

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
}
