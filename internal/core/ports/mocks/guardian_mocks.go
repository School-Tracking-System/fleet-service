package mocks

import (
	"context"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockGuardianRepository mocks GuardianRepository for testing.
type MockGuardianRepository struct {
	mock.Mock
}

func (m *MockGuardianRepository) Create(ctx context.Context, g *domain.Guardian) error {
	return m.Called(ctx, g).Error(0)
}

func (m *MockGuardianRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockGuardianRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Guardian, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Guardian), args.Error(1)
}

func (m *MockGuardianRepository) GetByStudentID(ctx context.Context, studentID uuid.UUID) ([]*domain.Guardian, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Guardian), args.Error(1)
}

func (m *MockGuardianRepository) GetStudentsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Student, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Student), args.Error(1)
}

// MockGuardianService mocks GuardianService for testing.
type MockGuardianService struct {
	mock.Mock
}

func (m *MockGuardianService) LinkGuardian(ctx context.Context, req services.LinkGuardianRequest) (*domain.Guardian, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Guardian), args.Error(1)
}

func (m *MockGuardianService) GetGuardiansByStudent(ctx context.Context, studentID uuid.UUID) ([]*domain.Guardian, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Guardian), args.Error(1)
}

func (m *MockGuardianService) GetStudentsByGuardian(ctx context.Context, userID uuid.UUID) ([]*domain.Student, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Student), args.Error(1)
}

func (m *MockGuardianService) UnlinkGuardian(ctx context.Context, guardianID uuid.UUID) error {
	return m.Called(ctx, guardianID).Error(0)
}
