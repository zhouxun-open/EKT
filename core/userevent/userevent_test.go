package userevent

import (
	"encoding/json"
	"fmt"
	"github.com/EducationEKT/EKT/core/types"
	"github.com/EducationEKT/EKT/crypto"
	"sort"
	"testing"
)

func TestSort(t *testing.T) {
	list := make(SortedUserEvent, 0)
	for i := 10; i > 0; i-- {
		tx := types.Transaction{Nonce: int64(i)}
		list = append(list, tx)
	}
	sort.Sort(list)
	for i := 0; i < len(list)-1; i++ {
		if list[i].GetNonce() > list[i+1].GetNonce() {
			t.Fail()
		}
	}
}

func TestSortedUserEvent_QuikInsert(t *testing.T) {
	from := crypto.Sha3_256([]byte("123"))
	to := crypto.Sha3_256([]byte("456"))
	event1 := types.NewTransaction(from, to, 0, 0, 0, 1, "", "")
	event2 := types.NewTransaction(from, to, 0, 0, 0, 2, "", "")
	event3 := types.NewTransaction(from, to, 0, 0, 0, 5, "", "")
	event4 := types.NewTransaction(from, to, 0, 0, 0, 4, "", "")
	event5 := types.NewTransaction(from, to, 0, 0, 0, 3, "", "")
	data := event5.Bytes()
	var userevent IUserEvent
	err := json.Unmarshal(data, &userevent)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(userevent.(types.Transaction))
	events := make(SortedUserEvent, 0)

	events.QuikInsert(event1)
	events.QuikInsert(event2)
	events.QuikInsert(event3)
	events.QuikInsert(event4)
	events.QuikInsert(event5)

	for i := 0; i < len(events)-1; i++ {
		if events[i].GetNonce() > events[i+1].GetNonce() {
			t.Fail()
		}
	}
}

func BenchmarkSortedUserEvent_QuikInsert(b *testing.B) {
	from := crypto.Sha3_256([]byte("123"))
	to := crypto.Sha3_256([]byte("456"))
	events := make(SortedUserEvent, 0)
	for i := b.N; i >= 0; i-- {
		event := types.NewTransaction(from, to, 0, 0, 0, int64(i), "", "")
		events.QuikInsert(event)
	}
}
