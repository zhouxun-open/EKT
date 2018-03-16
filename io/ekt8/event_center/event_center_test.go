package event_center

import (
	"fmt"
	"testing"
)

const (
	TestEventName = "TestEvent"
)

func init() {
}

func handler1(param *EventParam) (resp *EventResp, err error) {
	fmt.Println("handler1")
	return nil, nil
}

func handler2(param *EventParam) (resp *EventResp, err error) {
	fmt.Println("handler2")
	return nil, nil
}

func handler3(param *EventParam) (resp *EventResp, err error) {
	fmt.Println("handler3")
	data := make(map[string]interface{})
	data["handler"] = "handler3"
	return &EventResp{Params: data}, nil
}

func handler4(param *EventParam) (resp *EventResp, err error) {
	fmt.Println("handler4")
	return &EventResp{param.Sticker}, nil
}

func handler5(param *EventParam) (resp *EventResp, err error) {
	param.Sticker["handler5"] = true
	return nil, nil
}

func handler6(param *EventParam) (resp *EventResp, err error) {
	param.Sticker["handler6"] = true
	return nil, nil
}

func TestPublishEventOrder(t *testing.T) {
	RegistEvent(TestEventName, handler1, 1)
	RegistEvent(TestEventName, handler2, 2)
	resp, err := PublishEvent(TestEventName, &EventParam{
		Params:  make(map[string]interface{}),
		Sticker: make(map[string]interface{}),
	})
	if err != nil || resp != nil {
		t.FailNow()
	}
}

func TestPublishEventInterrupt(t *testing.T) {
	RegistEvent(TestEventName, handler3, 1)
	RegistEvent(TestEventName, handler4, 0)
	resp, err := PublishEvent(TestEventName, &EventParam{
		Params:  make(map[string]interface{}),
		Sticker: make(map[string]interface{}),
	})
	if err != nil {
		t.FailNow()
	}
	if resp.Params["handler"] != "handler3" {
		t.FailNow()
	}
	fmt.Println(*resp)
}

func TestSticker(t *testing.T) {
	RegistEvent(TestEventName, handler5, 1)
	RegistEvent(TestEventName, handler4, 0)
	RegistEvent(TestEventName, handler6, 2)
	resp, err := PublishEvent(TestEventName, &EventParam{
		Params:  make(map[string]interface{}),
		Sticker: make(map[string]interface{}),
	})
	if err != nil {
		t.FailNow()
	}
	if resp.Params["handler5"] != true || resp.Params["handler6"] != true {
		t.FailNow()
	}
	fmt.Println(*resp)
}
