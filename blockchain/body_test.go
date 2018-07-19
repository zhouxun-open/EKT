package blockchain

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/EducationEKT/EKT/userevent"
)

func TestBlockBody_AddEventResult(t *testing.T) {
	body := NewBlockBody(1)
	eventResult := userevent.EventResult{
		EventId: "123",
		Success: true,
		Reason:  "",
	}
	body.AddEventResult(eventResult)
	eventResult2 := userevent.EventResult{
		EventId: "122",
		Success: false,
		Reason:  "err sign",
	}
	body.AddEventResult(eventResult2)
	bodyData, _ := json.Marshal(body)
	fmt.Println(string(bodyData))
}
