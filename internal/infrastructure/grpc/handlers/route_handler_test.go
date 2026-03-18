package handlers

import (
	"context"
	"testing"
	"time"

	pb "github.com/fercho/school-tracking/proto/gen/fleet/v1"
	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/mocks"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newRouteHandler(svc *mocks.MockRouteService) pb.RouteServiceServer {
	return NewRouteHandler(svc, zap.NewNop())
}

func TestCreateRouteHandler(t *testing.T) {
	ctx := context.Background()
	schoolID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.MockRouteService)
		h := newRouteHandler(mockSvc)

		route := &domain.Route{
			ID:           uuid.New(),
			Name:         "Ruta Norte",
			SchoolID:     schoolID,
			Direction:    domain.RouteDirectionToSchool,
			ScheduleTime: "07:00",
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		mockSvc.On("CreateRoute", ctx, mock.MatchedBy(func(r services.CreateRouteRequest) bool {
			return r.Name == "Ruta Norte" && r.SchoolID == schoolID
		})).Return(route, nil)

		resp, err := h.CreateRoute(ctx, &pb.CreateRouteRequest{
			Name:         "Ruta Norte",
			SchoolId:     schoolID.String(),
			ScheduleTime: "07:00",
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, route.ID.String(), resp.Route.Id)
		assert.Equal(t, "Ruta Norte", resp.Route.Name)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid school_id returns InvalidArgument", func(t *testing.T) {
		mockSvc := new(mocks.MockRouteService)
		h := newRouteHandler(mockSvc)

		_, err := h.CreateRoute(ctx, &pb.CreateRouteRequest{
			Name:     "Ruta Norte",
			SchoolId: "not-a-uuid",
		})

		assert.Error(t, err)
		st, _ := status.FromError(err)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})
}

func TestGetRouteHandler(t *testing.T) {
	ctx := context.Background()
	routeID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.MockRouteService)
		h := newRouteHandler(mockSvc)

		route := &domain.Route{
			ID:        routeID,
			Name:      "Ruta Sur",
			SchoolID:  uuid.New(),
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockSvc.On("GetRoute", ctx, routeID).Return(route, nil)

		resp, err := h.GetRoute(ctx, &pb.GetRouteRequest{Id: routeID.String()})

		assert.NoError(t, err)
		assert.Equal(t, routeID.String(), resp.Route.Id)
		mockSvc.AssertExpectations(t)
	})

	t.Run("not found returns NotFound", func(t *testing.T) {
		mockSvc := new(mocks.MockRouteService)
		h := newRouteHandler(mockSvc)

		mockSvc.On("GetRoute", ctx, routeID).Return(nil, domain.ErrRouteNotFound)

		_, err := h.GetRoute(ctx, &pb.GetRouteRequest{Id: routeID.String()})

		assert.Error(t, err)
		st, _ := status.FromError(err)
		assert.Equal(t, codes.NotFound, st.Code())
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid id returns InvalidArgument", func(t *testing.T) {
		mockSvc := new(mocks.MockRouteService)
		h := newRouteHandler(mockSvc)

		_, err := h.GetRoute(ctx, &pb.GetRouteRequest{Id: "bad-id"})

		assert.Error(t, err)
		st, _ := status.FromError(err)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})
}

func TestUpdateStopOrderHandler(t *testing.T) {
	ctx := context.Background()
	stopID := uuid.New()
	routeID := uuid.New()
	studentID := uuid.New()

	t.Run("success returns updated stop", func(t *testing.T) {
		mockSvc := new(mocks.MockRouteService)
		h := newRouteHandler(mockSvc)

		updatedStop := &domain.RouteStop{
			ID:        stopID,
			RouteID:   routeID,
			StudentID: studentID,
			Order:     5,
			Location:  domain.Location{Longitude: -78.1, Latitude: -0.2},
			CreatedAt: time.Now(),
		}

		mockSvc.On("UpdateStopOrder", ctx, stopID, 5).Return(nil)
		mockSvc.On("GetStop", ctx, stopID).Return(updatedStop, nil)

		resp, err := h.UpdateStopOrder(ctx, &pb.UpdateStopOrderRequest{
			Id:       stopID.String(),
			NewOrder: 5,
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, stopID.String(), resp.Stop.Id)
		assert.Equal(t, int32(5), resp.Stop.Order)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid stop_id returns InvalidArgument", func(t *testing.T) {
		mockSvc := new(mocks.MockRouteService)
		h := newRouteHandler(mockSvc)

		_, err := h.UpdateStopOrder(ctx, &pb.UpdateStopOrderRequest{
			Id:       "not-a-uuid",
			NewOrder: 1,
		})

		assert.Error(t, err)
		st, _ := status.FromError(err)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("UpdateStopOrder service error maps to gRPC error", func(t *testing.T) {
		mockSvc := new(mocks.MockRouteService)
		h := newRouteHandler(mockSvc)

		mockSvc.On("UpdateStopOrder", ctx, stopID, 2).Return(domain.ErrStopNotFound)

		_, err := h.UpdateStopOrder(ctx, &pb.UpdateStopOrderRequest{
			Id:       stopID.String(),
			NewOrder: 2,
		})

		assert.Error(t, err)
		st, _ := status.FromError(err)
		assert.Equal(t, codes.NotFound, st.Code())
		mockSvc.AssertExpectations(t)
	})

	t.Run("GetStop error after update is propagated", func(t *testing.T) {
		mockSvc := new(mocks.MockRouteService)
		h := newRouteHandler(mockSvc)

		mockSvc.On("UpdateStopOrder", ctx, stopID, 3).Return(nil)
		mockSvc.On("GetStop", ctx, stopID).Return(nil, domain.ErrStopNotFound)

		_, err := h.UpdateStopOrder(ctx, &pb.UpdateStopOrderRequest{
			Id:       stopID.String(),
			NewOrder: 3,
		})

		assert.Error(t, err)
		st, _ := status.FromError(err)
		assert.Equal(t, codes.NotFound, st.Code())
		mockSvc.AssertExpectations(t)
	})
}

func TestAddStopHandler(t *testing.T) {
	ctx := context.Background()
	routeID := uuid.New()
	studentID := uuid.New()
	stopID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.MockRouteService)
		h := newRouteHandler(mockSvc)

		stop := &domain.RouteStop{
			ID:        stopID,
			RouteID:   routeID,
			StudentID: studentID,
			Order:     1,
			CreatedAt: time.Now(),
		}

		mockSvc.On("AddStop", ctx, mock.MatchedBy(func(r services.AddStopRequest) bool {
			return r.RouteID == routeID && r.StudentID == studentID && r.Order == 1
		})).Return(stop, nil)

		resp, err := h.AddStop(ctx, &pb.AddStopRequest{
			RouteId:   routeID.String(),
			StudentId: studentID.String(),
			Order:     1,
			Location:  &pb.GeoPoint{Longitude: -78.1, Latitude: -0.2},
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, stopID.String(), resp.Stop.Id)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid route_id returns InvalidArgument", func(t *testing.T) {
		mockSvc := new(mocks.MockRouteService)
		h := newRouteHandler(mockSvc)

		_, err := h.AddStop(ctx, &pb.AddStopRequest{
			RouteId:   "bad-id",
			StudentId: studentID.String(),
		})

		assert.Error(t, err)
		st, _ := status.FromError(err)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})
}

func TestRemoveStopHandler(t *testing.T) {
	ctx := context.Background()
	stopID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.MockRouteService)
		h := newRouteHandler(mockSvc)

		mockSvc.On("RemoveStop", ctx, stopID).Return(nil)

		resp, err := h.RemoveStop(ctx, &pb.RemoveStopRequest{Id: stopID.String()})

		assert.NoError(t, err)
		assert.True(t, resp.Success)
		mockSvc.AssertExpectations(t)
	})

	t.Run("stop not found returns NotFound", func(t *testing.T) {
		mockSvc := new(mocks.MockRouteService)
		h := newRouteHandler(mockSvc)

		mockSvc.On("RemoveStop", ctx, stopID).Return(domain.ErrStopNotFound)

		_, err := h.RemoveStop(ctx, &pb.RemoveStopRequest{Id: stopID.String()})

		assert.Error(t, err)
		st, _ := status.FromError(err)
		assert.Equal(t, codes.NotFound, st.Code())
		mockSvc.AssertExpectations(t)
	})
}

func TestGetRouteStopsHandler(t *testing.T) {
	ctx := context.Background()
	routeID := uuid.New()

	t.Run("success returns list of stops", func(t *testing.T) {
		mockSvc := new(mocks.MockRouteService)
		h := newRouteHandler(mockSvc)

		stops := []*domain.RouteStop{
			{ID: uuid.New(), RouteID: routeID, Order: 1, StudentID: uuid.New(), CreatedAt: time.Now()},
			{ID: uuid.New(), RouteID: routeID, Order: 2, StudentID: uuid.New(), CreatedAt: time.Now()},
		}
		mockSvc.On("GetRouteStops", ctx, routeID).Return(stops, nil)

		resp, err := h.GetRouteStops(ctx, &pb.GetRouteStopsRequest{RouteId: routeID.String()})

		assert.NoError(t, err)
		assert.Len(t, resp.Stops, 2)
		assert.Equal(t, int32(1), resp.Stops[0].Order)
		assert.Equal(t, int32(2), resp.Stops[1].Order)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid route_id returns InvalidArgument", func(t *testing.T) {
		mockSvc := new(mocks.MockRouteService)
		h := newRouteHandler(mockSvc)

		_, err := h.GetRouteStops(ctx, &pb.GetRouteStopsRequest{RouteId: "not-a-uuid"})

		assert.Error(t, err)
		st, _ := status.FromError(err)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})
}
