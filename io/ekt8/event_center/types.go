package event_center

type EventParam struct {
	Params  map[string]interface{}
	Sticker map[string]interface{}
}

type EventResp struct {
	Params map[string]interface{}
}

type EventHandler func(*EventParam) (resp *EventResp, err error)
type RespListener func(EventResp, error)

type EventMultiHandler struct {
	Priority int
	Handler  EventHandler
}

type Handlers []EventMultiHandler

type EventCenter struct {
	EventMapping map[string]Handlers
}

func (handlers Handlers) Len() int {
	return len(handlers)
}

func (handlers Handlers) Swap(i, j int) {
	handlers[i], handlers[j] = handlers[j], handlers[i]
}

func (handlers Handlers) Less(i, j int) bool {
	return handlers[i].Priority > handlers[j].Priority
}
