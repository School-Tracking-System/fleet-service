package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockEventPublisher is a mock implementation of resources.EventPublisher.
type MockEventPublisher struct {
	mock.Mock
}

// Publish records the call and returns a mocked error.
func (m *MockEventPublisher) Publish(ctx context.Context, subject string, payload []byte) error {
	args := m.Called(ctx, subject, payload)
	return args.Error(0)
}
