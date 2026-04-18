package repositories

import (
	"context"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/google/uuid"
)

// SchoolContactRepository defines the persistence contract for SchoolContact entities.
type SchoolContactRepository interface {
	Create(ctx context.Context, contact *domain.SchoolContact) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.SchoolContact, error)
	GetBySchoolID(ctx context.Context, schoolID uuid.UUID) ([]*domain.SchoolContact, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
