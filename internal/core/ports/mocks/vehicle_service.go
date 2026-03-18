package mocks

import (
	"context"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockVehicleService struct {
	mock.Mock
}

func (m *MockVehicleService) CreateVehicle(ctx context.Context, req services.CreateVehicleRequest) (*domain.Vehicle, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Vehicle), args.Error(1)
}

func (m *MockVehicleService) UpdateVehicle(ctx context.Context, req services.UpdateVehicleRequest) (*domain.Vehicle, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Vehicle), args.Error(1)
}

func (m *MockVehicleService) GetVehicle(ctx context.Context, id uuid.UUID) (*domain.Vehicle, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Vehicle), args.Error(1)
}

func (m *MockVehicleService) ListVehicles(ctx context.Context, limit, offset int) ([]*domain.Vehicle, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.Vehicle), args.Int(1), args.Error(2)
}
