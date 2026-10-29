package eventbus

import "context"

// Event is the generic message contract used by the event bus.
type Event interface {
	Type() string
	Payload() any
}

// EventBus abstracts the messaging transport used by services.
type EventBus interface {
	Publish(context.Context, Event) error
}
