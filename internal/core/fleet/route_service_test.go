package fleet

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/mocks"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestCreateRoute(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()
	repo := new(mocks.MockRouteRepository)
	svc := NewRouteService(repo, log)

	req := services.CreateRouteRequest{
		Name:         "Ruta 1",
		SchoolID:     uuid.New(),
		ScheduleTime: "07:00",
	}

	repo.On("Create", ctx, mock.MatchedBy(func(r *domain.Route) bool {
		return r.Name == req.Name
	})).Return(nil)

	res, err := svc.CreateRoute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	repo.AssertExpectations(t)
}

func TestAddStop(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()
	repo := new(mocks.MockRouteRepository)
	svc := NewRouteService(repo, log)

	req := services.AddStopRequest{
		RouteID:   uuid.New(),
		StudentID: uuid.New(),
		Order:     1,
	}

	repo.On("CreateStop", ctx, mock.MatchedBy(func(s *domain.RouteStop) bool {
		return s.RouteID == req.RouteID && s.StudentID == req.StudentID
	})).Return(nil)

	res, err := svc.AddStop(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestUpdateStopOrder(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()
	stopID := uuid.New()

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockRouteRepository)
		svc := NewRouteService(repo, log)

		repo.On("UpdateStopOrder", ctx, stopID, 3).Return(nil)

		err := svc.UpdateStopOrder(ctx, stopID, 3)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("repo error is propagated", func(t *testing.T) {
		repo := new(mocks.MockRouteRepository)
		svc := NewRouteService(repo, log)

		repo.On("UpdateStopOrder", ctx, stopID, 3).Return(errors.New("db error"))

		err := svc.UpdateStopOrder(ctx, stopID, 3)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestGetStop(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()
	stopID := uuid.New()
	routeID := uuid.New()
	studentID := uuid.New()

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockRouteRepository)
		svc := NewRouteService(repo, log)

		expected := &domain.RouteStop{
			ID:        stopID,
			RouteID:   routeID,
			StudentID: studentID,
			Order:     2,
			CreatedAt: time.Now(),
		}
		repo.On("GetStopByID", ctx, stopID).Return(expected, nil)

		stop, err := svc.GetStop(ctx, stopID)

		assert.NoError(t, err)
		assert.Equal(t, expected, stop)
		assert.Equal(t, 2, stop.Order)
		repo.AssertExpectations(t)
	})

	t.Run("stop not found", func(t *testing.T) {
		repo := new(mocks.MockRouteRepository)
		svc := NewRouteService(repo, log)

		repo.On("GetStopByID", ctx, stopID).Return(nil, domain.ErrStopNotFound)

		stop, err := svc.GetStop(ctx, stopID)

		assert.Error(t, err)
		assert.Nil(t, stop)
		assert.ErrorIs(t, err, domain.ErrStopNotFound)
		repo.AssertExpectations(t)
	})
}

func TestRemoveStop(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()
	stopID := uuid.New()

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockRouteRepository)
		svc := NewRouteService(repo, log)

		repo.On("DeleteStop", ctx, stopID).Return(nil)

		err := svc.RemoveStop(ctx, stopID)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestGetRouteStops(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()
	routeID := uuid.New()

	t.Run("success returns stops", func(t *testing.T) {
		repo := new(mocks.MockRouteRepository)
		svc := NewRouteService(repo, log)

		stops := []*domain.RouteStop{
			{ID: uuid.New(), RouteID: routeID, Order: 1},
			{ID: uuid.New(), RouteID: routeID, Order: 2},
		}
		repo.On("GetStopsByRouteID", ctx, routeID).Return(stops, nil)

		result, err := svc.GetRouteStops(ctx, routeID)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		repo.AssertExpectations(t)
	})

	t.Run("empty route returns empty slice", func(t *testing.T) {
		repo := new(mocks.MockRouteRepository)
		svc := NewRouteService(repo, log)

		repo.On("GetStopsByRouteID", ctx, routeID).Return([]*domain.RouteStop{}, nil)

		result, err := svc.GetRouteStops(ctx, routeID)

		assert.NoError(t, err)
		assert.Empty(t, result)
		repo.AssertExpectations(t)
	})
}
