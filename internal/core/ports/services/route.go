package services

import (
	"context"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/google/uuid"
)

// RouteService defines the core business logic for Route and Stop management.
type RouteService interface {
	CreateRoute(ctx context.Context, req CreateRouteRequest) (*domain.Route, error)
	UpdateRoute(ctx context.Context, req UpdateRouteRequest) (*domain.Route, error)
	GetRoute(ctx context.Context, id uuid.UUID) (*domain.Route, error)
	ListRoutes(ctx context.Context, limit, offset int) ([]*domain.Route, int, error)
	ListRoutesBySchool(ctx context.Context, schoolID uuid.UUID, limit, offset int) ([]*domain.Route, int, error)

	AddStop(ctx context.Context, req AddStopRequest) (*domain.RouteStop, error)
	RemoveStop(ctx context.Context, stopID uuid.UUID) error
	UpdateStopOrder(ctx context.Context, stopID uuid.UUID, newOrder int) error
	GetStop(ctx context.Context, stopID uuid.UUID) (*domain.RouteStop, error)
	GetRouteStops(ctx context.Context, routeID uuid.UUID) ([]*domain.RouteStop, error)
}

type CreateRouteRequest struct {
	Name         string
	Description  string
	VehicleID    *uuid.UUID
	DriverID     *uuid.UUID
	SchoolID     uuid.UUID
	Direction    domain.RouteDirection
	ScheduleTime string
}

type UpdateRouteRequest struct {
	ID           uuid.UUID
	Name         string
	Description  string
	VehicleID    **uuid.UUID
	DriverID     **uuid.UUID
	Direction    domain.RouteDirection
	ScheduleTime string
	IsActive     *bool
}

type AddStopRequest struct {
	RouteID   uuid.UUID
	StudentID uuid.UUID
	Order     int
	Location  domain.Location
	Address   string
	EstTime   string
}
