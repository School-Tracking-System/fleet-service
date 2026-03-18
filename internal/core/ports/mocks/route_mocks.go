package mocks

import (
	"context"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRouteRepository struct {
	mock.Mock
}

func (m *MockRouteRepository) Create(ctx context.Context, rt *domain.Route) error {
	return m.Called(ctx, rt).Error(0)
}

func (m *MockRouteRepository) Update(ctx context.Context, rt *domain.Route) error {
	return m.Called(ctx, rt).Error(0)
}

func (m *MockRouteRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Route, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Route), args.Error(1)
}

func (m *MockRouteRepository) List(ctx context.Context, limit, offset int) ([]*domain.Route, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.Route), args.Int(1), args.Error(2)
}

func (m *MockRouteRepository) ListBySchool(ctx context.Context, schoolID uuid.UUID, limit, offset int) ([]*domain.Route, int, error) {
	args := m.Called(ctx, schoolID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.Route), args.Int(1), args.Error(2)
}

func (m *MockRouteRepository) CreateStop(ctx context.Context, s *domain.RouteStop) error {
	return m.Called(ctx, s).Error(0)
}

func (m *MockRouteRepository) DeleteStop(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockRouteRepository) UpdateStopOrder(ctx context.Context, id uuid.UUID, order int) error {
	return m.Called(ctx, id, order).Error(0)
}

func (m *MockRouteRepository) GetStopsByRouteID(ctx context.Context, routeID uuid.UUID) ([]*domain.RouteStop, error) {
	args := m.Called(ctx, routeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.RouteStop), args.Error(1)
}

func (m *MockRouteRepository) GetStopByID(ctx context.Context, id uuid.UUID) (*domain.RouteStop, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RouteStop), args.Error(1)
}

type MockRouteService struct {
	mock.Mock
}

func (m *MockRouteService) CreateRoute(ctx context.Context, req services.CreateRouteRequest) (*domain.Route, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Route), args.Error(1)
}

func (m *MockRouteService) UpdateRoute(ctx context.Context, req services.UpdateRouteRequest) (*domain.Route, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Route), args.Error(1)
}

func (m *MockRouteService) GetRoute(ctx context.Context, id uuid.UUID) (*domain.Route, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Route), args.Error(1)
}

func (m *MockRouteService) ListRoutes(ctx context.Context, limit, offset int) ([]*domain.Route, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.Route), args.Int(1), args.Error(2)
}

func (m *MockRouteService) ListRoutesBySchool(ctx context.Context, schoolID uuid.UUID, limit, offset int) ([]*domain.Route, int, error) {
	args := m.Called(ctx, schoolID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.Route), args.Int(1), args.Error(2)
}

func (m *MockRouteService) AddStop(ctx context.Context, req services.AddStopRequest) (*domain.RouteStop, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RouteStop), args.Error(1)
}

func (m *MockRouteService) RemoveStop(ctx context.Context, stopID uuid.UUID) error {
	return m.Called(ctx, stopID).Error(0)
}

func (m *MockRouteService) UpdateStopOrder(ctx context.Context, stopID uuid.UUID, newOrder int) error {
	return m.Called(ctx, stopID, newOrder).Error(0)
}

func (m *MockRouteService) GetStop(ctx context.Context, stopID uuid.UUID) (*domain.RouteStop, error) {
	args := m.Called(ctx, stopID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RouteStop), args.Error(1)
}

func (m *MockRouteService) GetRouteStops(ctx context.Context, routeID uuid.UUID) ([]*domain.RouteStop, error) {
	args := m.Called(ctx, routeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.RouteStop), args.Error(1)
}
