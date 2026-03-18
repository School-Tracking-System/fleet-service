package resources

import "context"

// EventPublisher defines the contract for publishing domain events to a message broker.
type EventPublisher interface {
	Publish(ctx context.Context, subject string, payload []byte) error
}
