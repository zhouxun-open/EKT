package rlp

import (
	"fmt"
	"testing"

	"github.com/EducationEKT/EKT/io/ekt8/MPTPlus"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func TestEncodeStruct(t *testing.T) {
	s1 := MPTPlus.TrieSonInfo{[]byte("A"), "A"}
	s2 := MPTPlus.TrieSonInfo{[]byte("World"), "World"}
	node := &MPTPlus.TrieNode{Sons: []MPTPlus.TrieSonInfo{s1, s2}, Leaf: false, PathValue: "Hello"}
	bts, err := Encode(node)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	} else {
		fmt.Println(hexutil.Encode(bts))
	}
}

func TestDecodeStruct(t *testing.T) {
	bts, err := hexutil.Decode("0xd8d0c24141cc85576f726c6485576f726c64808548656c6c6f")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	} else {
		var node MPTPlus.TrieNode
		err = Decode(bts, &node)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		} else {
			fmt.Println(node)
		}
	}
}

func TestEncode(t *testing.T) {
	ary := []string{"hello", "world"}
	bts, err := Encode(ary)
	if err != nil || "0xcc8568656c6c6f85776f726c64" != hexutil.Encode(bts) {
		t.Fail()
	} else {
		fmt.Println(hexutil.Encode(bts))
	}
}

func TestDecode(t *testing.T) {
	bts, err := hexutil.Decode("0xcc8568656c6c6f85776f726c64")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	} else {
		var ary []string
		err = Decode(bts, &ary)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		} else {
			fmt.Println(ary)
		}
	}
}
