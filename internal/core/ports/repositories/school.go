package repositories

import (
	"context"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/google/uuid"
)

// SchoolRepository defines the persistence contract for School entities.
type SchoolRepository interface {
	Create(ctx context.Context, school *domain.School) error
	Update(ctx context.Context, school *domain.School) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.School, error)
	List(ctx context.Context, limit, offset int) ([]*domain.School, int, error)
}
