package handlers

import (
	"context"
	"errors"

	pb "github.com/fercho/school-tracking/proto/gen/fleet/v1"
	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/services"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type routeHandler struct {
	pb.UnimplementedRouteServiceServer
	service services.RouteService
	log     *zap.Logger
}

func NewRouteHandler(service services.RouteService, log *zap.Logger) pb.RouteServiceServer {
	return &routeHandler{service: service, log: log}
}

func (h *routeHandler) CreateRoute(ctx context.Context, req *pb.CreateRouteRequest) (*pb.RouteResponse, error) {
	h.log.Debug("CreateRoute request", zap.String("name", req.Name))

	schoolID, err := uuid.Parse(req.SchoolId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid school_id")
	}

	var vehicleID, driverID *uuid.UUID
	if req.VehicleId != "" {
		if id, err := uuid.Parse(req.VehicleId); err == nil {
			vehicleID = &id
		}
	}
	if req.DriverId != "" {
		if id, err := uuid.Parse(req.DriverId); err == nil {
			driverID = &id
		}
	}

	route, err := h.service.CreateRoute(ctx, services.CreateRouteRequest{
		Name:         req.Name,
		Description:  req.Description,
		VehicleID:    vehicleID,
		DriverID:     driverID,
		SchoolID:     schoolID,
		Direction:    domain.RouteDirection(req.Direction),
		ScheduleTime: req.ScheduleTime,
	})
	if err != nil {
		return nil, h.mapError(err)
	}

	return &pb.RouteResponse{Route: domainToProtoRoute(route)}, nil
}

func (h *routeHandler) UpdateRoute(ctx context.Context, req *pb.UpdateRouteRequest) (*pb.RouteResponse, error) {
	h.log.Debug("UpdateRoute request", zap.String("id", req.Id))

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid route_id")
	}

	// Logic for vehicle/driver id updates (allowing unset)
	var vIDPtr, dIDPtr **uuid.UUID
	if req.VehicleId != "" {
		vID, _ := uuid.Parse(req.VehicleId)
		vIDRef := &vID
		vIDPtr = &vIDRef
	}
	if req.DriverId != "" {
		dID, _ := uuid.Parse(req.DriverId)
		dIDRef := &dID
		dIDPtr = &dIDRef
	}

	route, err := h.service.UpdateRoute(ctx, services.UpdateRouteRequest{
		ID:           id,
		Name:         req.Name,
		Description:  req.Description,
		VehicleID:    vIDPtr,
		DriverID:     dIDPtr,
		Direction:    domain.RouteDirection(req.Direction),
		ScheduleTime: req.ScheduleTime,
		IsActive:     req.IsActive,
	})
	if err != nil {
		return nil, h.mapError(err)
	}

	return &pb.RouteResponse{Route: domainToProtoRoute(route)}, nil
}

func (h *routeHandler) GetRoute(ctx context.Context, req *pb.GetRouteRequest) (*pb.RouteResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid route_id")
	}

	route, err := h.service.GetRoute(ctx, id)
	if err != nil {
		return nil, h.mapError(err)
	}

	return &pb.RouteResponse{Route: domainToProtoRoute(route)}, nil
}

func (h *routeHandler) ListRoutes(ctx context.Context, req *pb.ListRoutesRequest) (*pb.ListRoutesResponse, error) {
	routes, total, err := h.service.ListRoutes(ctx, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, h.mapError(err)
	}

	var pbRoutes []*pb.Route
	for _, r := range routes {
		pbRoutes = append(pbRoutes, domainToProtoRoute(r))
	}
	return &pb.ListRoutesResponse{Routes: pbRoutes, TotalCount: int32(total)}, nil
}

func (h *routeHandler) ListRoutesBySchool(ctx context.Context, req *pb.ListRoutesBySchoolRequest) (*pb.ListRoutesResponse, error) {
	schoolID, err := uuid.Parse(req.SchoolId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid school_id")
	}

	routes, total, err := h.service.ListRoutesBySchool(ctx, schoolID, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, h.mapError(err)
	}

	var pbRoutes []*pb.Route
	for _, r := range routes {
		pbRoutes = append(pbRoutes, domainToProtoRoute(r))
	}
	return &pb.ListRoutesResponse{Routes: pbRoutes, TotalCount: int32(total)}, nil
}

func (h *routeHandler) AddStop(ctx context.Context, req *pb.AddStopRequest) (*pb.StopResponse, error) {
	routeID, err := uuid.Parse(req.RouteId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid route_id")
	}
	studentID, err := uuid.Parse(req.StudentId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid student_id")
	}

	stop, err := h.service.AddStop(ctx, services.AddStopRequest{
		RouteID:   routeID,
		StudentID: studentID,
		Order:     int(req.Order),
		Location:  *protoToLocation(req.Location),
		Address:   req.Address,
		EstTime:   req.EstTime,
	})
	if err != nil {
		return nil, h.mapError(err)
	}

	return &pb.StopResponse{Stop: domainToProtoStop(stop)}, nil
}

func (h *routeHandler) RemoveStop(ctx context.Context, req *pb.RemoveStopRequest) (*pb.RemoveStopResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid stop_id")
	}
	if err := h.service.RemoveStop(ctx, id); err != nil {
		return nil, h.mapError(err)
	}
	return &pb.RemoveStopResponse{Success: true}, nil
}

func (h *routeHandler) UpdateStopOrder(ctx context.Context, req *pb.UpdateStopOrderRequest) (*pb.StopResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid stop_id")
	}
	if err := h.service.UpdateStopOrder(ctx, id, int(req.NewOrder)); err != nil {
		return nil, h.mapError(err)
	}

	stop, err := h.service.GetStop(ctx, id)
	if err != nil {
		return nil, h.mapError(err)
	}

	return &pb.StopResponse{Stop: domainToProtoStop(stop)}, nil
}

func (h *routeHandler) GetRouteStops(ctx context.Context, req *pb.GetRouteStopsRequest) (*pb.ListStopsResponse, error) {
	id, err := uuid.Parse(req.RouteId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid route_id")
	}
	stops, err := h.service.GetRouteStops(ctx, id)
	if err != nil {
		return nil, h.mapError(err)
	}
	var pbStops []*pb.RouteStop
	for _, s := range stops {
		pbStops = append(pbStops, domainToProtoStop(s))
	}
	return &pb.ListStopsResponse{Stops: pbStops}, nil
}

func domainToProtoRoute(r *domain.Route) *pb.Route {
	p := &pb.Route{
		Id:           r.ID.String(),
		Name:         r.Name,
		Description:  r.Description,
		SchoolId:     r.SchoolID.String(),
		Direction:    string(r.Direction),
		ScheduleTime: r.ScheduleTime,
		IsActive:     r.IsActive,
		CreatedAt:    timestamppb.New(r.CreatedAt),
		UpdatedAt:    timestamppb.New(r.UpdatedAt),
	}
	if r.VehicleID != nil {
		p.VehicleId = r.VehicleID.String()
	}
	if r.DriverID != nil {
		p.DriverId = r.DriverID.String()
	}
	return p
}

func domainToProtoStop(s *domain.RouteStop) *pb.RouteStop {
	return &pb.RouteStop{
		Id:        s.ID.String(),
		RouteId:   s.RouteID.String(),
		StudentId: s.StudentID.String(),
		Order:     int32(s.Order),
		Location: &pb.GeoPoint{
			Longitude: s.Location.Longitude,
			Latitude:  s.Location.Latitude,
		},
		Address:   s.Address,
		EstTime:   s.EstTime,
		CreatedAt: timestamppb.New(s.CreatedAt),
	}
}

func (h *routeHandler) mapError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, domain.ErrRouteNotFound), errors.Is(err, domain.ErrStopNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrInvalidRoute), errors.Is(err, domain.ErrInvalidStop):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		h.log.Error("Unexpected route error", zap.Error(err))
		return status.Errorf(codes.Internal, "an unexpected error occurred")
	}
}

var RouteHandlerModule = fx.Provide(NewRouteHandler)
