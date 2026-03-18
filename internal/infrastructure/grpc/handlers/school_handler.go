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

type schoolHandler struct {
	pb.UnimplementedSchoolServiceServer
	service services.SchoolService
	log     *zap.Logger
}

// NewSchoolHandler creates a new gRPC handler for the SchoolService.
func NewSchoolHandler(service services.SchoolService, log *zap.Logger) pb.SchoolServiceServer {
	return &schoolHandler{
		service: service,
		log:     log,
	}
}

func (h *schoolHandler) CreateSchool(ctx context.Context, req *pb.CreateSchoolRequest) (*pb.SchoolResponse, error) {
	h.log.Debug("CreateSchool request received", zap.String("name", req.Name))

	school, err := h.service.CreateSchool(ctx, services.CreateSchoolRequest{
		Name:     req.Name,
		Address:  req.Address,
		Location: protoToLocation(req.Location),
		Phone:    req.Phone,
		Email:    req.Email,
	})
	if err != nil {
		return nil, h.mapError(err)
	}

	return &pb.SchoolResponse{School: domainToProtoSchool(school)}, nil
}

func (h *schoolHandler) UpdateSchool(ctx context.Context, req *pb.UpdateSchoolRequest) (*pb.SchoolResponse, error) {
	h.log.Debug("UpdateSchool request received", zap.String("id", req.Id))

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid school ID format")
	}

	school, err := h.service.UpdateSchool(ctx, services.UpdateSchoolRequest{
		ID:       id,
		Name:     req.Name,
		Address:  req.Address,
		Location: protoToLocation(req.Location),
		Phone:    req.Phone,
		Email:    req.Email,
	})
	if err != nil {
		return nil, h.mapError(err)
	}

	return &pb.SchoolResponse{School: domainToProtoSchool(school)}, nil
}

func (h *schoolHandler) GetSchool(ctx context.Context, req *pb.GetSchoolRequest) (*pb.SchoolResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid school ID format")
	}

	school, err := h.service.GetSchool(ctx, id)
	if err != nil {
		return nil, h.mapError(err)
	}

	return &pb.SchoolResponse{School: domainToProtoSchool(school)}, nil
}

func (h *schoolHandler) ListSchools(ctx context.Context, req *pb.ListSchoolsRequest) (*pb.ListSchoolsResponse, error) {
	limit := int(req.Limit)
	if limit == 0 {
		limit = 10
	}

	schools, total, err := h.service.ListSchools(ctx, limit, int(req.Offset))
	if err != nil {
		return nil, h.mapError(err)
	}

	var pbSchools []*pb.School
	for _, s := range schools {
		pbSchools = append(pbSchools, domainToProtoSchool(s))
	}

	return &pb.ListSchoolsResponse{
		Schools:    pbSchools,
		TotalCount: int32(total),
	}, nil
}

// domainToProtoSchool converts a domain School to its protobuf representation.
func domainToProtoSchool(s *domain.School) *pb.School {
	p := &pb.School{
		Id:        s.ID.String(),
		Name:      s.Name,
		Address:   s.Address,
		Phone:     s.Phone,
		Email:     s.Email,
		CreatedAt: timestamppb.New(s.CreatedAt),
		UpdatedAt: timestamppb.New(s.UpdatedAt),
	}
	if s.Location != nil {
		p.Location = &pb.GeoPoint{
			Longitude: s.Location.Longitude,
			Latitude:  s.Location.Latitude,
		}
	}
	return p
}

// protoToLocation converts a proto GeoPoint to a domain Location (nil-safe).
func protoToLocation(geo *pb.GeoPoint) *domain.Location {
	if geo == nil {
		return nil
	}
	return &domain.Location{Longitude: geo.Longitude, Latitude: geo.Latitude}
}

// mapError converts domain School errors to gRPC status errors.
func (h *schoolHandler) mapError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, domain.ErrSchoolNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrDuplicateSchool):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrInvalidSchool):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		h.log.Error("Unexpected school error", zap.Error(err))
		return status.Errorf(codes.Internal, "an unexpected error occurred")
	}
}

var SchoolHandlerModule = fx.Provide(NewSchoolHandler)
