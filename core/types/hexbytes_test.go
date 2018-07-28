package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
)

type TestStruct struct {
	Name    HexBytes `json:"name"`
	Address HexBytes `json:"address"`
}

func TestHexBytes_MarshalJSON(t *testing.T) {
	name, _ := hex.DecodeString(hex.EncodeToString([]byte("name")))
	address, _ := hex.DecodeString(hex.EncodeToString([]byte("address")))
	struct1 := TestStruct{Name: name, Address: address}
	data, err := json.Marshal(struct1)
	fmt.Println(string(data), err)
	var struct2 TestStruct
	err = json.Unmarshal(data, &struct2)
	fmt.Println(struct2, err)
	data, err = json.Marshal(struct2)
	fmt.Println(string(data), err)
}
