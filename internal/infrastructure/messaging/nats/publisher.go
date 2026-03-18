package nats

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
)

const (
	streamName     = "FLEET_EVENTS"
	streamSubjects = "fleet.>"
)

// SubjectVehicleCreated is the NATS subject for vehicle creation events.
const SubjectVehicleCreated = "fleet.vehicle.created"

// SubjectVehicleUpdated is the NATS subject for vehicle update events.
const SubjectVehicleUpdated = "fleet.vehicle.updated"

type publisher struct {
	js nats.JetStreamContext
}

// NewPublisher creates a NATS JetStream publisher and ensures the FLEET_EVENTS stream exists.
func NewPublisher(nc *nats.Conn) (*publisher, error) {
	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to get JetStream context: %w", err)
	}

	if err := ensureStream(js); err != nil {
		return nil, err
	}

	return &publisher{js: js}, nil
}

// Publish sends a JSON payload to the given NATS subject.
func (p *publisher) Publish(_ context.Context, subject string, payload []byte) error {
	_, err := p.js.Publish(subject, payload)
	return err
}

// ensureStream creates the FLEET_EVENTS JetStream stream if it does not already exist.
func ensureStream(js nats.JetStreamContext) error {
	_, err := js.StreamInfo(streamName)
	if err == nil {
		return nil
	}
	if err != nats.ErrStreamNotFound {
		return fmt.Errorf("failed to check stream: %w", err)
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{streamSubjects},
		Storage:  nats.MemoryStorage,
	})
	if err != nil {
		return fmt.Errorf("failed to create stream %s: %w", streamName, err)
	}
	return nil
}
