package repositories

import (
	"context"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/google/uuid"
)

// DriverRepository defines the persistence contract for Driver entities.
type DriverRepository interface {
	Create(ctx context.Context, driver *domain.Driver) error
	Update(ctx context.Context, driver *domain.Driver) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Driver, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Driver, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Driver, int, error)
}
