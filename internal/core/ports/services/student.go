package services

import (
	"context"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/google/uuid"
)

// StudentService defines the core business logic contract for Student operations.
type StudentService interface {
	RegisterStudent(ctx context.Context, req RegisterStudentRequest) (*domain.Student, error)
	UpdateStudent(ctx context.Context, req UpdateStudentRequest) (*domain.Student, error)
	GetStudent(ctx context.Context, id uuid.UUID) (*domain.Student, error)
	ListStudents(ctx context.Context, limit, offset int) ([]*domain.Student, int, error)
	ListStudentsBySchool(ctx context.Context, schoolID uuid.UUID, limit, offset int) ([]*domain.Student, int, error)
	DeactivateStudent(ctx context.Context, id uuid.UUID) error
}

// RegisterStudentRequest encapsulates the data required to register a new student.
type RegisterStudentRequest struct {
	FirstName      string
	LastName       string
	Grade          string
	SchoolID       uuid.UUID
	PickupLocation *domain.Location
	PickupAddress  string
	PhotoURL       string
}

// UpdateStudentRequest encapsulates the data for a partial student update.
type UpdateStudentRequest struct {
	ID             uuid.UUID
	FirstName      string
	LastName       string
	Grade          string
	PickupLocation *domain.Location
	PickupAddress  string
	PhotoURL       string
}
