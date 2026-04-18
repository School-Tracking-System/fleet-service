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
	service        services.SchoolService
	contactService services.SchoolContactService
	log            *zap.Logger
}

// NewSchoolHandler creates a new gRPC handler for the SchoolService.
func NewSchoolHandler(service services.SchoolService, contactService services.SchoolContactService, log *zap.Logger) pb.SchoolServiceServer {
	return &schoolHandler{
		service:        service,
		contactService: contactService,
		log:            log,
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

func (h *schoolHandler) AddContact(ctx context.Context, req *pb.AddContactRequest) (*pb.ContactResponse, error) {
	schoolID, err := uuid.Parse(req.SchoolId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid school_id format")
	}
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id format")
	}

	contact, err := h.contactService.AddContact(ctx, services.AddContactRequest{
		SchoolID: schoolID,
		UserID:   userID,
		Position: req.Position,
	})
	if err != nil {
		return nil, h.mapContactError(err)
	}

	return &pb.ContactResponse{Contact: domainToProtoContact(contact)}, nil
}

func (h *schoolHandler) RemoveContact(ctx context.Context, req *pb.RemoveContactRequest) (*pb.RemoveContactResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid contact ID format")
	}

	if err := h.contactService.RemoveContact(ctx, id); err != nil {
		return nil, h.mapContactError(err)
	}

	return &pb.RemoveContactResponse{Success: true}, nil
}

func (h *schoolHandler) ListContacts(ctx context.Context, req *pb.ListContactsRequest) (*pb.ListContactsResponse, error) {
	schoolID, err := uuid.Parse(req.SchoolId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid school_id format")
	}

	contacts, err := h.contactService.ListContacts(ctx, schoolID)
	if err != nil {
		return nil, h.mapContactError(err)
	}

	var pbContacts []*pb.SchoolContact
	for _, c := range contacts {
		pbContacts = append(pbContacts, domainToProtoContact(c))
	}

	return &pb.ListContactsResponse{Contacts: pbContacts}, nil
}

func domainToProtoContact(c *domain.SchoolContact) *pb.SchoolContact {
	return &pb.SchoolContact{
		Id:        c.ID.String(),
		SchoolId:  c.SchoolID.String(),
		UserId:    c.UserID.String(),
		Position:  c.Position,
		IsActive:  c.IsActive,
		CreatedAt: timestamppb.New(c.CreatedAt),
	}
}

func (h *schoolHandler) mapContactError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, domain.ErrContactNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrDuplicateContact):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrInvalidContact):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		h.log.Error("Unexpected contact error", zap.Error(err))
		return status.Errorf(codes.Internal, "an unexpected error occurred")
	}
}

var SchoolHandlerModule = fx.Provide(NewSchoolHandler)
