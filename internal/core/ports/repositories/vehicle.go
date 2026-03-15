package repositories

import (
	"context"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/google/uuid"
)

// VehicleRepository defines the persistence contract for Vehicle entities.
type VehicleRepository interface {
	Create(ctx context.Context, vehicle *domain.Vehicle) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Vehicle, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Vehicle, int, error)
}
