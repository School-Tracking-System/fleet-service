package handlers

import (
	"context"
	"errors"
	"time"

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
	h.log.Debug("CreateVehicle request received", zap.String("plate", req.Plate))

	vehicle, err := h.service.CreateVehicle(ctx, services.CreateVehicleRequest{
		Plate:         req.Plate,
		Brand:         req.Brand,
		Model:         req.Model,
		Year:          int(req.Year),
		Capacity:      int(req.Capacity),
		Color:         req.Color,
		VehicleType:   domain.VehicleType(req.VehicleType),
		ChassisNum:    req.ChassisNum,
		InsuranceExp:  protoTimestampToTime(req.InsuranceExp),
		TechReviewExp: protoTimestampToTime(req.TechReviewExp),
	})
	if err != nil {
		return nil, h.mapError(err)
	}

	return &pb.VehicleResponse{Vehicle: domainToProto(vehicle)}, nil
}

func (h *vehicleHandler) UpdateVehicle(ctx context.Context, req *pb.UpdateVehicleRequest) (*pb.VehicleResponse, error) {
	h.log.Debug("UpdateVehicle request received", zap.String("id", req.Id))

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid vehicle ID format")
	}

	vehicle, err := h.service.UpdateVehicle(ctx, services.UpdateVehicleRequest{
		ID:            id,
		Brand:         req.Brand,
		Model:         req.Model,
		Year:          int(req.Year),
		Capacity:      int(req.Capacity),
		Status:        domain.VehicleStatus(req.Status),
		Color:         req.Color,
		VehicleType:   domain.VehicleType(req.VehicleType),
		ChassisNum:    req.ChassisNum,
		InsuranceExp:  protoTimestampToTime(req.InsuranceExp),
		TechReviewExp: protoTimestampToTime(req.TechReviewExp),
	})
	if err != nil {
		return nil, h.mapError(err)
	}

	return &pb.VehicleResponse{Vehicle: domainToProto(vehicle)}, nil
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

	return &pb.VehicleResponse{Vehicle: domainToProto(vehicle)}, nil
}

func (h *vehicleHandler) ListVehicles(ctx context.Context, req *pb.ListVehiclesRequest) (*pb.ListVehiclesResponse, error) {
	limit := int(req.Limit)
	if limit == 0 {
		limit = 10
	}

	vehicles, total, err := h.service.ListVehicles(ctx, limit, int(req.Offset))
	if err != nil {
		return nil, h.mapError(err)
	}

	var pbVehicles []*pb.Vehicle
	for _, v := range vehicles {
		pbVehicles = append(pbVehicles, domainToProto(v))
	}

	return &pb.ListVehiclesResponse{
		Vehicles:   pbVehicles,
		TotalCount: int32(total),
	}, nil
}

// domainToProto converts a domain Vehicle to its protobuf representation.
func domainToProto(v *domain.Vehicle) *pb.Vehicle {
	p := &pb.Vehicle{
		Id:          v.ID.String(),
		Plate:       v.Plate,
		Brand:       v.Brand,
		Model:       v.Model,
		Year:        int32(v.Year),
		Capacity:    int32(v.Capacity),
		Status:      string(v.Status),
		Color:       v.Color,
		VehicleType: string(v.VehicleType),
		ChassisNum:  v.ChassisNum,
		CreatedAt:   timestamppb.New(v.CreatedAt),
		UpdatedAt:   timestamppb.New(v.UpdatedAt),
	}
	if v.InsuranceExp != nil {
		p.InsuranceExp = timestamppb.New(*v.InsuranceExp)
	}
	if v.TechReviewExp != nil {
		p.TechReviewExp = timestamppb.New(*v.TechReviewExp)
	}
	return p
}

// protoTimestampToTime converts a protobuf Timestamp to a *time.Time (nil-safe).
func protoTimestampToTime(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
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
