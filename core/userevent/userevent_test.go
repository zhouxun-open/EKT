package userevent

import (
	"fmt"
	"github.com/EducationEKT/EKT/crypto"
	"sort"
	"testing"
)

func TestSort(t *testing.T) {
	list := make(SortedUserEvent, 0)
	for i := 10; i > 0; i-- {
		tx := Transaction{Nonce: int64(i)}
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
	event1 := NewTransaction(from, to, 0, 0, 0, 1, "", "")
	event2 := NewTransaction(from, to, 0, 0, 0, 2, "", "")
	event3 := NewTransaction(from, to, 0, 0, 0, 5, "", "")
	event4 := NewTransaction(from, to, 0, 0, 0, 4, "", "")
	event5 := NewTransaction(from, to, 0, 0, 0, 3, "", "")

	events := make(SortedUserEvent, 0)

	events = events.QuickInsert(event1)
	events = events.QuickInsert(event2)
	events = events.QuickInsert(event3)
	events = events.QuickInsert(event4)
	events = events.QuickInsert(event5)

	fmt.Println(events)
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
		event := NewTransaction(from, to, 0, 0, 0, int64(i), "", "")
		events.QuickInsert(event)
	}
}
