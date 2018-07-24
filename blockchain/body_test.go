package blockchain

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestBlockBody_AddEventResult(t *testing.T) {
	body := NewBlockBody()
	bodyData, _ := json.Marshal(body)
	fmt.Println(string(bodyData))
}
