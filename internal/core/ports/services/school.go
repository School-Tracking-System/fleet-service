package services

import (
	"context"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/google/uuid"
)

// SchoolService defines the core business logic contract for School operations.
type SchoolService interface {
	CreateSchool(ctx context.Context, req CreateSchoolRequest) (*domain.School, error)
	UpdateSchool(ctx context.Context, req UpdateSchoolRequest) (*domain.School, error)
	GetSchool(ctx context.Context, id uuid.UUID) (*domain.School, error)
	ListSchools(ctx context.Context, limit, offset int) ([]*domain.School, int, error)
}

// CreateSchoolRequest encapsulates the data required to register a new school.
type CreateSchoolRequest struct {
	Name     string
	Address  string
	Location *domain.Location
	Phone    string
	Email    string
}

// UpdateSchoolRequest encapsulates the data for a partial school update.
type UpdateSchoolRequest struct {
	ID       uuid.UUID
	Name     string
	Address  string
	Location *domain.Location
	Phone    string
	Email    string
}
