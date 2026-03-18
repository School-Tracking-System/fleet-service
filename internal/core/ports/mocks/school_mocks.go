package mocks

import (
	"context"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockSchoolRepository mocks SchoolRepository for testing.
type MockSchoolRepository struct {
	mock.Mock
}

func (m *MockSchoolRepository) Create(ctx context.Context, s *domain.School) error {
	return m.Called(ctx, s).Error(0)
}

func (m *MockSchoolRepository) Update(ctx context.Context, s *domain.School) error {
	return m.Called(ctx, s).Error(0)
}

func (m *MockSchoolRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.School, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.School), args.Error(1)
}

func (m *MockSchoolRepository) List(ctx context.Context, limit, offset int) ([]*domain.School, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.School), args.Int(1), args.Error(2)
}

// MockSchoolService mocks SchoolService for handler testing.
type MockSchoolService struct {
	mock.Mock
}

func (m *MockSchoolService) CreateSchool(ctx context.Context, req services.CreateSchoolRequest) (*domain.School, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.School), args.Error(1)
}

func (m *MockSchoolService) UpdateSchool(ctx context.Context, req services.UpdateSchoolRequest) (*domain.School, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.School), args.Error(1)
}

func (m *MockSchoolService) GetSchool(ctx context.Context, id uuid.UUID) (*domain.School, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.School), args.Error(1)
}

func (m *MockSchoolService) ListSchools(ctx context.Context, limit, offset int) ([]*domain.School, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.School), args.Int(1), args.Error(2)
}
