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

type guardianHandler struct {
	pb.UnimplementedGuardianServiceServer
	service services.GuardianService
	log     *zap.Logger
}

// NewGuardianHandler creates a new gRPC handler for the GuardianService.
func NewGuardianHandler(service services.GuardianService, log *zap.Logger) pb.GuardianServiceServer {
	return &guardianHandler{
		service: service,
		log:     log,
	}
}

func (h *guardianHandler) LinkGuardian(ctx context.Context, req *pb.LinkGuardianRequest) (*pb.GuardianResponse, error) {
	h.log.Debug("LinkGuardian request received", zap.String("user_id", req.UserId), zap.String("student_id", req.StudentId))

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id format")
	}

	studentID, err := uuid.Parse(req.StudentId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid student_id format")
	}

	guardian, err := h.service.LinkGuardian(ctx, services.LinkGuardianRequest{
		UserID:    userID,
		StudentID: studentID,
		Relation:  domain.GuardianRelation(req.Relation),
		IsPrimary: req.IsPrimary,
	})
	if err != nil {
		return nil, h.mapError(err)
	}

	return &pb.GuardianResponse{Guardian: domainToProtoGuardian(guardian)}, nil
}

func (h *guardianHandler) UnlinkGuardian(ctx context.Context, req *pb.UnlinkGuardianRequest) (*pb.UnlinkGuardianResponse, error) {
	h.log.Debug("UnlinkGuardian request received", zap.String("id", req.Id))

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid guardian ID format")
	}

	if err := h.service.UnlinkGuardian(ctx, id); err != nil {
		return nil, h.mapError(err)
	}

	return &pb.UnlinkGuardianResponse{Success: true}, nil
}

func (h *guardianHandler) GetGuardiansByStudent(ctx context.Context, req *pb.GetGuardiansByStudentRequest) (*pb.GuardianListResponse, error) {
	studentID, err := uuid.Parse(req.StudentId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid student_id format")
	}

	guardians, err := h.service.GetGuardiansByStudent(ctx, studentID)
	if err != nil {
		return nil, h.mapError(err)
	}

	var pbGuardians []*pb.Guardian
	for _, g := range guardians {
		pbGuardians = append(pbGuardians, domainToProtoGuardian(g))
	}

	return &pb.GuardianListResponse{Guardians: pbGuardians}, nil
}

func (h *guardianHandler) GetStudentsByGuardian(ctx context.Context, req *pb.GetStudentsByGuardianRequest) (*pb.ListStudentsResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id format")
	}

	students, err := h.service.GetStudentsByGuardian(ctx, userID)
	if err != nil {
		return nil, h.mapError(err)
	}

	var pbStudents []*pb.Student
	for _, s := range students {
		pbStudents = append(pbStudents, domainToProtoStudent(s))
	}

	return &pb.ListStudentsResponse{
		Students:    pbStudents,
		TotalCount: int32(len(students)),
	}, nil
}

// domainToProtoGuardian converts a domain Guardian to its protobuf representation.
func domainToProtoGuardian(g *domain.Guardian) *pb.Guardian {
	return &pb.Guardian{
		Id:        g.ID.String(),
		UserId:    g.UserID.String(),
		StudentId: g.StudentID.String(),
		Relation:  string(g.Relation),
		IsPrimary: g.IsPrimary,
		CreatedAt: timestamppb.New(g.CreatedAt),
	}
}

// mapError converts domain Guardian errors to gRPC status errors.
func (h *guardianHandler) mapError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, domain.ErrGuardianNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrDuplicateGuardian):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrInvalidGuardian):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrStudentNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		h.log.Error("Unexpected guardian error", zap.Error(err))
		return status.Errorf(codes.Internal, "an unexpected error occurred")
	}
}

var GuardianHandlerModule = fx.Provide(NewGuardianHandler)
