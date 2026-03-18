package mocks

import (
	"context"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockVehicleRepository struct {
	mock.Mock
}

func (m *MockVehicleRepository) Create(ctx context.Context, v *domain.Vehicle) error {
	args := m.Called(ctx, v)
	return args.Error(0)
}

func (m *MockVehicleRepository) Update(ctx context.Context, v *domain.Vehicle) error {
	args := m.Called(ctx, v)
	return args.Error(0)
}

func (m *MockVehicleRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Vehicle, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Vehicle), args.Error(1)
}

func (m *MockVehicleRepository) List(ctx context.Context, limit, offset int) ([]*domain.Vehicle, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.Vehicle), args.Int(1), args.Error(2)
}
