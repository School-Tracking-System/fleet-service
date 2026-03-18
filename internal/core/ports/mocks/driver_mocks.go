package mocks

import (
	"context"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockDriverRepository mocks the DriverRepository for testing.
type MockDriverRepository struct {
	mock.Mock
}

func (m *MockDriverRepository) Create(ctx context.Context, d *domain.Driver) error {
	return m.Called(ctx, d).Error(0)
}

func (m *MockDriverRepository) Update(ctx context.Context, d *domain.Driver) error {
	return m.Called(ctx, d).Error(0)
}

func (m *MockDriverRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Driver, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Driver), args.Error(1)
}

func (m *MockDriverRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Driver, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Driver), args.Error(1)
}

func (m *MockDriverRepository) List(ctx context.Context, limit, offset int) ([]*domain.Driver, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.Driver), args.Int(1), args.Error(2)
}

// MockDriverService mocks the DriverService for handler testing.
type MockDriverService struct {
	mock.Mock
}

func (m *MockDriverService) RegisterDriver(ctx context.Context, req services.RegisterDriverRequest) (*domain.Driver, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Driver), args.Error(1)
}

func (m *MockDriverService) UpdateDriver(ctx context.Context, req services.UpdateDriverRequest) (*domain.Driver, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Driver), args.Error(1)
}

func (m *MockDriverService) GetDriver(ctx context.Context, id uuid.UUID) (*domain.Driver, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Driver), args.Error(1)
}

func (m *MockDriverService) GetDriverByUserID(ctx context.Context, userID uuid.UUID) (*domain.Driver, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Driver), args.Error(1)
}

func (m *MockDriverService) ListDrivers(ctx context.Context, limit, offset int) ([]*domain.Driver, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.Driver), args.Int(1), args.Error(2)
}
