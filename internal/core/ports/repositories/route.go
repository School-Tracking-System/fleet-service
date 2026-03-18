package repositories

import (
	"context"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/google/uuid"
)

// RouteRepository defines the persistence contract for Route and Stop entities.
type RouteRepository interface {
	Create(ctx context.Context, route *domain.Route) error
	Update(ctx context.Context, route *domain.Route) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Route, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Route, int, error)
	ListBySchool(ctx context.Context, schoolID uuid.UUID, limit, offset int) ([]*domain.Route, int, error)

	CreateStop(ctx context.Context, stop *domain.RouteStop) error
	DeleteStop(ctx context.Context, id uuid.UUID) error
	UpdateStopOrder(ctx context.Context, id uuid.UUID, order int) error
	GetStopsByRouteID(ctx context.Context, routeID uuid.UUID) ([]*domain.RouteStop, error)
	GetStopByID(ctx context.Context, id uuid.UUID) (*domain.RouteStop, error)
}
