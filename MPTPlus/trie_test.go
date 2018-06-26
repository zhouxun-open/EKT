package MPTPlus

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"testing"
	"time"

	"xserver/x_utils/x_random"

	"github.com/EducationEKT/EKT/db"
)

type KeyValue struct {
	Key   []byte
	Value []byte
}

type RandomKeyValues []KeyValue

func (this RandomKeyValues) Len() int {
	return len(this)
}

func (this RandomKeyValues) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this RandomKeyValues) Less(i, j int) bool {
	numstr := x_random.RandomNumber(2)
	num, _ := strconv.Atoi(numstr)
	return num%2 == 0
}

func TestMTPInsertAndGet(t *testing.T) {
	db, err := db.NewLevelDB("testTrie11")
	if err != nil {
		t.Fail()
	}
	defer db.DB.Close()
	var randomKeyValues RandomKeyValues = []KeyValue{
		KeyValue{[]byte("x"), []byte("this is value9")},
		KeyValue{[]byte("HZhouWorld"), []byte("this is value4")},
		KeyValue{[]byte("HZhouxun"), []byte("this is value3")},
		KeyValue{[]byte("helloworld"), []byte("this is value7")},
		KeyValue{[]byte("HelloWorld"), []byte("this is value1")},
		KeyValue{[]byte("HZhouWorld"), []byte("this is value6")},
		KeyValue{[]byte("HelloX"), []byte("this is value2")},
		KeyValue{[]byte("HelloWorld"), []byte("this is value5")},
		KeyValue{[]byte("zhouxun"), []byte("this is value8")},
	}

	trie1 := NewMTP(db)
	kvMap := make(map[string][]byte)
	for _, kv := range randomKeyValues {
		err = trie1.MustInsert(kv.Key, kv.Value)
		kvMap[string(kv.Key)] = kv.Value
		if err != nil {
			fmt.Println(err)
			t.Fail()
		}
	}
	for _, kv := range randomKeyValues {
		v, err := trie1.GetValue(kv.Key)
		if !bytes.Equal(v, kvMap[string(kv.Key)]) || err != nil {
			fmt.Printf("err=%v, key=%s, value=%s, dbValue=%s\n", err, string(kv.Key), string(kv.Value), string(v))
			t.Fail()
		}
	}
	fmt.Println("success")
}

func TestMTPRandomInsert(t *testing.T) {
	db, err := db.NewLevelDB("testTrie11")
	if err != nil {
		t.Fail()
	}
	defer db.DB.Close()
	var randomKeyValues RandomKeyValues = []KeyValue{
		KeyValue{[]byte("HZhouWorld1"), []byte("this is value4")},
		KeyValue{[]byte("HZhouxun"), []byte("this is value3")},
		KeyValue{[]byte("x"), []byte("this is value9")},
		KeyValue{[]byte("helloworld"), []byte("this is value7")},
		KeyValue{[]byte("HelloWorld2"), []byte("this is value1")},
		KeyValue{[]byte("HZhouWorld2"), []byte("this is value6")},
		KeyValue{[]byte("HelloX"), []byte("this is value2")},
		KeyValue{[]byte("HelloWorld1"), []byte("this is value5")},
		KeyValue{[]byte("zhouxun"), []byte("this is value8")},
	}

	start := time.Now()
	trie1 := NewMTP(db)
	for _, kv := range randomKeyValues {
		err = trie1.MustInsert(kv.Key, kv.Value)
		if err != nil {
			t.Fail()
		}
	}
	half := time.Now()
	fmt.Printf("half=%d\n", half.Nanosecond()-start.Nanosecond())

	sort.Sort(randomKeyValues)
	trie2 := NewMTP(db)
	for _, kv := range randomKeyValues {
		err = trie2.MustInsert(kv.Key, kv.Value)
		if err != nil {
			t.Fail()
		}
	}
	fmt.Printf("trie1.Root=%s, trie2.Root=%s, trie1==trie2=%v\n", hex.EncodeToString(trie1.Root), hex.EncodeToString(trie2.Root), bytes.Equal(trie2.Root, trie1.Root))
	if !bytes.Equal(trie1.Root, trie2.Root) {
		t.Fail()
	} else {
		fmt.Println("success")
	}
	finishTime := time.Now()
	fmt.Printf("finishTime=%d\n", finishTime.Nanosecond()-start.Nanosecond())
}
