package event

import (
	"fmt"
	"testing"
)

func TestEvent_ValidateEvent(t *testing.T) {
	newAccountEvent := &Event{NewAccountEvent, NewAccountParam{}}
	updatePublicKeyEvent := &Event{UpdatePublicKeyEvent, UpdatePublicKeyParam{}}
	if !newAccountEvent.ValidateEvent() || !updatePublicKeyEvent.ValidateEvent() {
		t.Fail()
	}
	errorEvent1 := &Event{UpdatePublicKeyEvent, NewAccountParam{}}
	errorEvent2 := &Event{NewAccountEvent, UpdatePublicKeyParam{}}
	if errorEvent1.ValidateEvent() || errorEvent2.ValidateEvent() {
		t.Fail()
	}
	fmt.Println("success")
}
