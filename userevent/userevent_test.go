package userevent

import (
	"github.com/EducationEKT/EKT/core/common"
	"sort"
	"testing"
)

func TestSort(t *testing.T) {
	list := make(SortedUserEvent, 0)
	for i := 10; i > 0; i-- {
		tx := common.Transaction{Nonce: int64(i)}
		list = append(list, tx)
	}
	sort.Sort(list)
	for i := 0; i < len(list)-1; i++ {
		if list[i].GetNonce() > list[i+1].GetNonce() {
			t.Fail()
		}
	}
}
