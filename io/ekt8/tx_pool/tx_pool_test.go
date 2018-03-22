package tx_pool

import (
	"testing"
)

func TestTxPool_Fetch(t *testing.T) {
	println(txPool.Fetch(0))
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
