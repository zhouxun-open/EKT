package event_center

import "sort"

var eventCenter EventCenter

func init() {
	eventCenter = EventCenter{make(map[string]Handlers)}
}

func GetInst() EventCenter {
	return eventCenter
}

func RegistEvent(eventName string, handler EventHandler, Priority int) {
	handlers, exist := eventCenter.EventMapping[eventName]
	if !exist {
		handlers = make([]EventMultiHandler, 0)
	}
	handlers = append(handlers, EventMultiHandler{Priority: Priority, Handler: handler})
	sort.Sort(handlers)
	eventCenter.EventMapping[eventName] = handlers
}

func PublishEvent(eventName string, param *EventParam) (resp *EventResp, err error) {
	handlers, exist := eventCenter.EventMapping[eventName]
	if !exist {
		return nil, NoSuchEvent
	}
	for _, handler := range handlers {
		resp, err = handler.Handler(param)
		if err != nil || resp != nil {
			return
		}
	}
	return
}
