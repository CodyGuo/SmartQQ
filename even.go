package smartqq

type EventHandler func(qq *QQClient)

type Event struct {
	handlers []EventHandler
}

func (this *Event) Attach(handler EventHandler) int {
	for i, h := range this.handlers {
		if h == nil {
			this.handlers[i] = handler
			return i
		}
	}

	this.handlers = append(this.handlers, handler)
	return len(this.handlers) - 1
}

func (this *Event) Detach(handle int) {
	this.handlers[handle] = nil
}

type EventPublisher struct {
	event Event
}

func (this *EventPublisher) Event() *Event {
	return &this.event
}

func (this *EventPublisher) Publish(qq *QQClient) {
	for _, handler := range this.event.handlers {
		if handler != nil {
			handler(qq)
		}
	}
}
