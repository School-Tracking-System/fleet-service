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

type studentHandler struct {
	pb.UnimplementedStudentServiceServer
	service services.StudentService
	log     *zap.Logger
}

// NewStudentHandler creates a new gRPC handler for the StudentService.
func NewStudentHandler(service services.StudentService, log *zap.Logger) pb.StudentServiceServer {
	return &studentHandler{
		service: service,
		log:     log,
	}
}

func (h *studentHandler) RegisterStudent(ctx context.Context, req *pb.RegisterStudentRequest) (*pb.StudentResponse, error) {
	h.log.Debug("RegisterStudent request received", zap.String("first_name", req.FirstName), zap.String("last_name", req.LastName))

	schoolID, err := uuid.Parse(req.SchoolId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid school_id format")
	}

	student, err := h.service.RegisterStudent(ctx, services.RegisterStudentRequest{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Grade:          req.Grade,
		SchoolID:       schoolID,
		PickupLocation: protoToLocation(req.PickupLocation),
		PickupAddress:  req.PickupAddress,
		PhotoURL:       req.PhotoUrl,
	})
	if err != nil {
		return nil, h.mapError(err)
	}

	return &pb.StudentResponse{Student: domainToProtoStudent(student)}, nil
}

func (h *studentHandler) UpdateStudent(ctx context.Context, req *pb.UpdateStudentRequest) (*pb.StudentResponse, error) {
	h.log.Debug("UpdateStudent request received", zap.String("id", req.Id))

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid student ID format")
	}

	student, err := h.service.UpdateStudent(ctx, services.UpdateStudentRequest{
		ID:             id,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Grade:          req.Grade,
		PickupLocation: protoToLocation(req.PickupLocation),
		PickupAddress:  req.PickupAddress,
		PhotoURL:       req.PhotoUrl,
	})
	if err != nil {
		return nil, h.mapError(err)
	}

	return &pb.StudentResponse{Student: domainToProtoStudent(student)}, nil
}

func (h *studentHandler) GetStudent(ctx context.Context, req *pb.GetStudentRequest) (*pb.StudentResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid student ID format")
	}

	student, err := h.service.GetStudent(ctx, id)
	if err != nil {
		return nil, h.mapError(err)
	}

	return &pb.StudentResponse{Student: domainToProtoStudent(student)}, nil
}

func (h *studentHandler) ListStudents(ctx context.Context, req *pb.ListStudentsRequest) (*pb.ListStudentsResponse, error) {
	limit := int(req.Limit)
	if limit == 0 {
		limit = 10
	}

	students, total, err := h.service.ListStudents(ctx, limit, int(req.Offset))
	if err != nil {
		return nil, h.mapError(err)
	}

	var pbStudents []*pb.Student
	for _, s := range students {
		pbStudents = append(pbStudents, domainToProtoStudent(s))
	}

	return &pb.ListStudentsResponse{
		Students:    pbStudents,
		TotalCount: int32(total),
	}, nil
}

func (h *studentHandler) ListStudentsBySchool(ctx context.Context, req *pb.ListStudentsBySchoolRequest) (*pb.ListStudentsResponse, error) {
	schoolID, err := uuid.Parse(req.SchoolId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid school_id format")
	}

	limit := int(req.Limit)
	if limit == 0 {
		limit = 10
	}

	students, total, err := h.service.ListStudentsBySchool(ctx, schoolID, limit, int(req.Offset))
	if err != nil {
		return nil, h.mapError(err)
	}

	var pbStudents []*pb.Student
	for _, s := range students {
		pbStudents = append(pbStudents, domainToProtoStudent(s))
	}

	return &pb.ListStudentsResponse{
		Students:    pbStudents,
		TotalCount: int32(total),
	}, nil
}

func (h *studentHandler) DeactivateStudent(ctx context.Context, req *pb.DeactivateStudentRequest) (*pb.StudentResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid student ID format")
	}

	if err := h.service.DeactivateStudent(ctx, id); err != nil {
		return nil, h.mapError(err)
	}

	student, err := h.service.GetStudent(ctx, id)
	if err != nil {
		return nil, h.mapError(err)
	}

	return &pb.StudentResponse{Student: domainToProtoStudent(student)}, nil
}

// domainToProtoStudent converts a domain Student to its protobuf representation.
func domainToProtoStudent(s *domain.Student) *pb.Student {
	p := &pb.Student{
		Id:            s.ID.String(),
		FirstName:     s.FirstName,
		LastName:      s.LastName,
		Grade:         s.Grade,
		SchoolId:      s.SchoolID.String(),
		PickupAddress: s.PickupAddress,
		PhotoUrl:      s.PhotoURL,
		IsActive:      s.IsActive,
		CreatedAt:     timestamppb.New(s.CreatedAt),
		UpdatedAt:     timestamppb.New(s.UpdatedAt),
	}
	if s.PickupLocation != nil {
		p.PickupLocation = &pb.GeoPoint{
			Longitude: s.PickupLocation.Longitude,
			Latitude:  s.PickupLocation.Latitude,
		}
	}
	return p
}

// mapError converts domain Student errors to gRPC status errors.
func (h *studentHandler) mapError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, domain.ErrStudentNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrInvalidStudent):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrSchoolNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		h.log.Error("Unexpected student error", zap.Error(err))
		return status.Errorf(codes.Internal, "an unexpected error occurred")
	}
}

var StudentHandlerModule = fx.Provide(NewStudentHandler)
