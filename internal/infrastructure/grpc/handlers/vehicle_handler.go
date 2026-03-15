package handlers

import (
	"context"
	"errors"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/services"
	pb "github.com/fercho/school-tracking/proto/gen/fleet/v1"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type vehicleHandler struct {
	pb.UnimplementedVehicleServiceServer
	service services.VehicleService
	log     *zap.Logger
}

func NewVehicleHandler(service services.VehicleService, log *zap.Logger) pb.VehicleServiceServer {
	return &vehicleHandler{
		service: service,
		log:     log,
	}
}

func (h *vehicleHandler) CreateVehicle(ctx context.Context, req *pb.CreateVehicleRequest) (*pb.VehicleResponse, error) {
	h.log.Debug("CreateVehicle request received",
		zap.String("plate", req.Plate),
		zap.String("brand", req.Brand),
		zap.String("model", req.Model),
		zap.Int32("year", req.Year),
		zap.Int32("capacity", req.Capacity),
	)

	vehicle, err := h.service.CreateVehicle(ctx, services.CreateVehicleRequest{
		Plate:    req.Plate,
		Brand:    req.Brand,
		Model:    req.Model,
		Year:     int(req.Year),
		Capacity: int(req.Capacity),
	})

	if err != nil {
		return nil, h.mapError(err)
	}

	return &pb.VehicleResponse{
		Vehicle: &pb.Vehicle{
			Id:        vehicle.ID.String(),
			Plate:     vehicle.Plate,
			Brand:     vehicle.Brand,
			Model:     vehicle.Model,
			Year:      int32(vehicle.Year),
			Capacity:  int32(vehicle.Capacity),
			Status:    string(vehicle.Status),
			CreatedAt: timestamppb.New(vehicle.CreatedAt),
			UpdatedAt: timestamppb.New(vehicle.UpdatedAt),
		},
	}, nil
}

func (h *vehicleHandler) GetVehicle(ctx context.Context, req *pb.GetVehicleRequest) (*pb.VehicleResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid vehicle ID format")
	}

	vehicle, err := h.service.GetVehicle(ctx, id)
	if err != nil {
		return nil, h.mapError(err)
	}

	return &pb.VehicleResponse{
		Vehicle: &pb.Vehicle{
			Id:        vehicle.ID.String(),
			Plate:     vehicle.Plate,
			Brand:     vehicle.Brand,
			Model:     vehicle.Model,
			Year:      int32(vehicle.Year),
			Capacity:  int32(vehicle.Capacity),
			Status:    string(vehicle.Status),
			CreatedAt: timestamppb.New(vehicle.CreatedAt),
			UpdatedAt: timestamppb.New(vehicle.UpdatedAt),
		},
	}, nil
}

func (h *vehicleHandler) ListVehicles(ctx context.Context, req *pb.ListVehiclesRequest) (*pb.ListVehiclesResponse, error) {
	limit := int(req.Limit)
	if limit == 0 {
		limit = 10
	}
	offset := int(req.Offset)

	vehicles, total, err := h.service.ListVehicles(ctx, limit, offset)
	if err != nil {
		return nil, h.mapError(err)
	}

	var pbVehicles []*pb.Vehicle
	for _, v := range vehicles {
		pbVehicles = append(pbVehicles, &pb.Vehicle{
			Id:        v.ID.String(),
			Plate:     v.Plate,
			Brand:     v.Brand,
			Model:     v.Model,
			Year:      int32(v.Year),
			Capacity:  int32(v.Capacity),
			Status:    string(v.Status),
			CreatedAt: timestamppb.New(v.CreatedAt),
			UpdatedAt: timestamppb.New(v.UpdatedAt),
		})
	}

	return &pb.ListVehiclesResponse{
		Vehicles:   pbVehicles,
		TotalCount: int32(total),
	}, nil
}

// mapError converts domain errors into appropriate gRPC status errors.
func (h *vehicleHandler) mapError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, domain.ErrVehicleNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrDuplicateVehicle):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrInvalidVehicle):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		h.log.Error("Unexpected error occurred", zap.Error(err))
		return status.Errorf(codes.Internal, "an unexpected error occurred")
	}
}

var Module = fx.Provide(NewVehicleHandler)
