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

type driverHandler struct {
	pb.UnimplementedDriverServiceServer
	service services.DriverService
	log     *zap.Logger
}

// NewDriverHandler creates a new gRPC handler for the DriverService.
func NewDriverHandler(service services.DriverService, log *zap.Logger) pb.DriverServiceServer {
	return &driverHandler{
		service: service,
		log:     log,
	}
}

func (h *driverHandler) RegisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.DriverResponse, error) {
	h.log.Debug("RegisterDriver request received", zap.String("user_id", req.UserId))

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id format")
	}

	if req.LicenseExpiry == nil {
		return nil, status.Errorf(codes.InvalidArgument, "license_expiry is required")
	}

	driver, err := h.service.RegisterDriver(ctx, services.RegisterDriverRequest{
		UserID:         userID,
		LicenseNumber:  req.LicenseNumber,
		LicenseType:    req.LicenseType,
		LicenseExpiry:  req.LicenseExpiry.AsTime(),
		CedulaID:       req.CedulaId,
		EmergencyPhone: req.EmergencyPhone,
	})
	if err != nil {
		return nil, h.mapError(err)
	}

	return &pb.DriverResponse{Driver: domainToProtoDriver(driver)}, nil
}

func (h *driverHandler) UpdateDriver(ctx context.Context, req *pb.UpdateDriverRequest) (*pb.DriverResponse, error) {
	h.log.Debug("UpdateDriver request received", zap.String("id", req.Id))

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid driver ID format")
	}

	updateReq := services.UpdateDriverRequest{
		ID:             id,
		LicenseType:    req.LicenseType,
		EmergencyPhone: req.EmergencyPhone,
		Status:         domain.DriverStatus(req.Status),
	}
	if req.LicenseExpiry != nil {
		t := req.LicenseExpiry.AsTime()
		updateReq.LicenseExpiry = &t
	}

	driver, err := h.service.UpdateDriver(ctx, updateReq)
	if err != nil {
		return nil, h.mapError(err)
	}

	return &pb.DriverResponse{Driver: domainToProtoDriver(driver)}, nil
}

func (h *driverHandler) GetDriver(ctx context.Context, req *pb.GetDriverRequest) (*pb.DriverResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid driver ID format")
	}

	driver, err := h.service.GetDriver(ctx, id)
	if err != nil {
		return nil, h.mapError(err)
	}

	return &pb.DriverResponse{Driver: domainToProtoDriver(driver)}, nil
}

func (h *driverHandler) GetDriverByUserID(ctx context.Context, req *pb.GetDriverByUserIDRequest) (*pb.DriverResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id format")
	}

	driver, err := h.service.GetDriverByUserID(ctx, userID)
	if err != nil {
		return nil, h.mapError(err)
	}

	return &pb.DriverResponse{Driver: domainToProtoDriver(driver)}, nil
}

func (h *driverHandler) ListDrivers(ctx context.Context, req *pb.ListDriversRequest) (*pb.ListDriversResponse, error) {
	limit := int(req.Limit)
	if limit == 0 {
		limit = 10
	}

	drivers, total, err := h.service.ListDrivers(ctx, limit, int(req.Offset))
	if err != nil {
		return nil, h.mapError(err)
	}

	var pbDrivers []*pb.Driver
	for _, d := range drivers {
		pbDrivers = append(pbDrivers, domainToProtoDriver(d))
	}

	return &pb.ListDriversResponse{
		Drivers:    pbDrivers,
		TotalCount: int32(total),
	}, nil
}

// domainToProtoDriver converts a domain Driver to its protobuf representation.
func domainToProtoDriver(d *domain.Driver) *pb.Driver {
	return &pb.Driver{
		Id:             d.ID.String(),
		UserId:         d.UserID.String(),
		LicenseNumber:  d.LicenseNumber,
		LicenseType:    d.LicenseType,
		LicenseExpiry:  timestamppb.New(d.LicenseExpiry),
		CedulaId:       d.CedulaID,
		EmergencyPhone: d.EmergencyPhone,
		Status:         string(d.Status),
		CreatedAt:      timestamppb.New(d.CreatedAt),
		UpdatedAt:      timestamppb.New(d.UpdatedAt),
	}
}

// mapError converts domain Driver errors to gRPC status errors.
func (h *driverHandler) mapError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, domain.ErrDriverNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrDuplicateDriver):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrInvalidDriver):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrDriverAlreadyLinked):
		return status.Error(codes.AlreadyExists, err.Error())
	default:
		h.log.Error("Unexpected driver error", zap.Error(err))
		return status.Errorf(codes.Internal, "an unexpected error occurred")
	}
}

var DriverHandlerModule = fx.Provide(NewDriverHandler)
