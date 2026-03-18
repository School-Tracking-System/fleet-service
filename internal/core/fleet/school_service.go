package fleet

import (
	"context"
	"fmt"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/repositories"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/services"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type schoolService struct {
	repo repositories.SchoolRepository
	log  *zap.Logger
}

// NewSchoolService creates a new business logic service for schools.
func NewSchoolService(repo repositories.SchoolRepository, log *zap.Logger) services.SchoolService {
	return &schoolService{
		repo: repo,
		log:  log,
	}
}

func (s *schoolService) CreateSchool(ctx context.Context, req services.CreateSchoolRequest) (*domain.School, error) {
	s.log.Info("Creating school", zap.String("name", req.Name))

	school, err := domain.NewSchool(domain.NewSchoolParams{
		Name:     req.Name,
		Address:  req.Address,
		Location: req.Location,
		Phone:    req.Phone,
		Email:    req.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("invalid school data: %w", err)
	}

	if err := s.repo.Create(ctx, school); err != nil {
		s.log.Error("Failed to persist school", zap.Error(err))
		return nil, fmt.Errorf("failed to create school: %w", err)
	}

	s.log.Info("School created", zap.String("id", school.ID.String()))
	return school, nil
}

func (s *schoolService) UpdateSchool(ctx context.Context, req services.UpdateSchoolRequest) (*domain.School, error) {
	s.log.Info("Updating school", zap.String("id", req.ID.String()))

	school, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("school not found: %w", err)
	}

	school.Apply(toSchoolPatch(req))

	if err := s.repo.Update(ctx, school); err != nil {
		s.log.Error("Failed to update school", zap.Error(err))
		return nil, fmt.Errorf("failed to update school: %w", err)
	}

	return school, nil
}

func (s *schoolService) GetSchool(ctx context.Context, id uuid.UUID) (*domain.School, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *schoolService) ListSchools(ctx context.Context, limit, offset int) ([]*domain.School, int, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.List(ctx, limit, offset)
}

// toSchoolPatch maps an UpdateSchoolRequest to a domain SchoolPatch.
func toSchoolPatch(req services.UpdateSchoolRequest) domain.SchoolPatch {
	return domain.SchoolPatch{
		Name:     nonEmptyStr(req.Name),
		Address:  nonEmptyStr(req.Address),
		Location: req.Location,
		Phone:    nonEmptyStr(req.Phone),
		Email:    nonEmptyStr(req.Email),
	}
}

// SchoolModule provides the school service to the fx dependency graph.
var SchoolModule = fx.Provide(NewSchoolService)
