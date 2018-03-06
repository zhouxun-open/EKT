package event

const (
	NewAccountEvent      = "NewAccount"
	UpdatePublicKeyEvent = "UpdatePublicKey"
)

type EventParam interface {
	EventType() string
}

type Event struct {
	EventType  string
	EventParam EventParam
}

type NewAccountParam struct {
	Address []byte
	PubKey  []byte
	Nonce   int
	EventId []byte
}

type UpdatePublicKeyParam struct {
	Address   []byte
	NewPubKey []byte
	msg       []byte
	Nonce     int
	Sign      []byte
}

func (newAccountParam NewAccountParam) EventType() string {
	return NewAccountEvent
}

func (updatePublicKeyParam UpdatePublicKeyParam) EventType() string {
	return UpdatePublicKeyEvent
}

func (event Event) ValidateEvent() bool {
	if event.EventParam.EventType() != event.EventType {
		return false
	}
	return true
}
