package outbox

import "context"

// Message describes an outbox record awaiting publication.
type Message struct {
	ID          string
	AggregateID string
	EventType   string
	Payload     string
	CreatedAt   string
}

// Relay publishes unpublished events. The real implementation is deliberately kept thin here.
type Relay interface {
	Publish(context.Context, Message) error
}
