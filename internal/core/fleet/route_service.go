package fleet

import (
	"context"
	"fmt"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/repositories"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/services"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type routeService struct {
	repo repositories.RouteRepository
	log  *zap.Logger
}

func NewRouteService(repo repositories.RouteRepository, log *zap.Logger) services.RouteService {
	return &routeService{repo: repo, log: log}
}

func (s *routeService) CreateRoute(ctx context.Context, req services.CreateRouteRequest) (*domain.Route, error) {
	s.log.Info("Creating route", zap.String("name", req.Name))

	route, err := domain.NewRoute(domain.NewRouteParams{
		Name:         req.Name,
		Description:  req.Description,
		VehicleID:    req.VehicleID,
		DriverID:     req.DriverID,
		SchoolID:     req.SchoolID,
		Direction:    req.Direction,
		ScheduleTime: req.ScheduleTime,
	})
	if err != nil {
		return nil, fmt.Errorf("invalid route data: %w", err)
	}

	if err := s.repo.Create(ctx, route); err != nil {
		s.log.Error("Failed to persist route", zap.Error(err))
		return nil, fmt.Errorf("failed to create route: %w", err)
	}

	return route, nil
}

func (s *routeService) UpdateRoute(ctx context.Context, req services.UpdateRouteRequest) (*domain.Route, error) {
	s.log.Info("Updating route", zap.String("id", req.ID.String()))

	route, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("route not found: %w", err)
	}

	route.Apply(domain.RoutePatch{
		Name:         nonEmptyStr(req.Name),
		Description:  nonEmptyStr(req.Description),
		VehicleID:    req.VehicleID,
		DriverID:     req.DriverID,
		Direction:    &req.Direction,
		ScheduleTime: nonEmptyStr(req.ScheduleTime),
		IsActive:     req.IsActive,
	})

	if err := s.repo.Update(ctx, route); err != nil {
		s.log.Error("Failed to update route", zap.Error(err))
		return nil, fmt.Errorf("failed to update route: %w", err)
	}

	return route, nil
}

func (s *routeService) GetRoute(ctx context.Context, id uuid.UUID) (*domain.Route, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *routeService) ListRoutes(ctx context.Context, limit, offset int) ([]*domain.Route, int, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *routeService) ListRoutesBySchool(ctx context.Context, schoolID uuid.UUID, limit, offset int) ([]*domain.Route, int, error) {
	return s.repo.ListBySchool(ctx, schoolID, limit, offset)
}

func (s *routeService) AddStop(ctx context.Context, req services.AddStopRequest) (*domain.RouteStop, error) {
	s.log.Info("Adding stop to route", zap.String("route_id", req.RouteID.String()), zap.String("student_id", req.StudentID.String()))

	stop, err := domain.NewRouteStop(domain.NewStopParams{
		RouteID:   req.RouteID,
		StudentID: req.StudentID,
		Order:     req.Order,
		Location:  req.Location,
		Address:   req.Address,
		EstTime:   req.EstTime,
	})
	if err != nil {
		return nil, fmt.Errorf("invalid stop data: %w", err)
	}

	if err := s.repo.CreateStop(ctx, stop); err != nil {
		s.log.Error("Failed to persist stop", zap.Error(err))
		return nil, fmt.Errorf("failed to add stop: %w", err)
	}

	return stop, nil
}

func (s *routeService) RemoveStop(ctx context.Context, stopID uuid.UUID) error {
	s.log.Info("Removing stop", zap.String("id", stopID.String()))
	return s.repo.DeleteStop(ctx, stopID)
}

func (s *routeService) UpdateStopOrder(ctx context.Context, stopID uuid.UUID, newOrder int) error {
	s.log.Info("Updating stop order", zap.String("id", stopID.String()), zap.Int("order", newOrder))
	return s.repo.UpdateStopOrder(ctx, stopID, newOrder)
}

func (s *routeService) GetStop(ctx context.Context, stopID uuid.UUID) (*domain.RouteStop, error) {
	return s.repo.GetStopByID(ctx, stopID)
}

func (s *routeService) GetRouteStops(ctx context.Context, routeID uuid.UUID) ([]*domain.RouteStop, error) {
	return s.repo.GetStopsByRouteID(ctx, routeID)
}

var RouteModule = fx.Provide(NewRouteService)
