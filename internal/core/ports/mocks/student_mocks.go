package mocks

import (
	"context"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockStudentRepository mocks StudentRepository for testing.
type MockStudentRepository struct {
	mock.Mock
}

func (m *MockStudentRepository) Create(ctx context.Context, s *domain.Student) error {
	return m.Called(ctx, s).Error(0)
}

func (m *MockStudentRepository) Update(ctx context.Context, s *domain.Student) error {
	return m.Called(ctx, s).Error(0)
}

func (m *MockStudentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Student, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Student), args.Error(1)
}

func (m *MockStudentRepository) List(ctx context.Context, limit, offset int) ([]*domain.Student, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.Student), args.Int(1), args.Error(2)
}

func (m *MockStudentRepository) ListBySchool(ctx context.Context, schoolID uuid.UUID, limit, offset int) ([]*domain.Student, int, error) {
	args := m.Called(ctx, schoolID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.Student), args.Int(1), args.Error(2)
}

// MockStudentService mocks StudentService for testing.
type MockStudentService struct {
	mock.Mock
}

func (m *MockStudentService) RegisterStudent(ctx context.Context, req services.RegisterStudentRequest) (*domain.Student, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Student), args.Error(1)
}

func (m *MockStudentService) UpdateStudent(ctx context.Context, req services.UpdateStudentRequest) (*domain.Student, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Student), args.Error(1)
}

func (m *MockStudentService) GetStudent(ctx context.Context, id uuid.UUID) (*domain.Student, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Student), args.Error(1)
}

func (m *MockStudentService) ListStudents(ctx context.Context, limit, offset int) ([]*domain.Student, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.Student), args.Int(1), args.Error(2)
}

func (m *MockStudentService) ListStudentsBySchool(ctx context.Context, schoolID uuid.UUID, limit, offset int) ([]*domain.Student, int, error) {
	args := m.Called(ctx, schoolID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.Student), args.Int(1), args.Error(2)
}

func (m *MockStudentService) DeactivateStudent(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
