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

type guardianService struct {
	repo repositories.GuardianRepository
	log  *zap.Logger
}

// NewGuardianService creates a new business logic service for guardians.
func NewGuardianService(repo repositories.GuardianRepository, log *zap.Logger) services.GuardianService {
	return &guardianService{repo: repo, log: log}
}

func (s *guardianService) LinkGuardian(ctx context.Context, req services.LinkGuardianRequest) (*domain.Guardian, error) {
	s.log.Info("Linking guardian to student",
		zap.String("user_id", req.UserID.String()),
		zap.String("student_id", req.StudentID.String()),
	)

	guardian, err := domain.NewGuardian(domain.NewGuardianParams{
		UserID:    req.UserID,
		StudentID: req.StudentID,
		Relation:  req.Relation,
		IsPrimary: req.IsPrimary,
	})
	if err != nil {
		return nil, fmt.Errorf("invalid guardian data: %w", err)
	}

	if err := s.repo.Create(ctx, guardian); err != nil {
		s.log.Error("Failed to link guardian", zap.Error(err))
		return nil, fmt.Errorf("failed to link guardian: %w", err)
	}

	s.log.Info("Guardian linked", zap.String("id", guardian.ID.String()))
	return guardian, nil
}

func (s *guardianService) GetGuardiansByStudent(ctx context.Context, studentID uuid.UUID) ([]*domain.Guardian, error) {
	return s.repo.GetByStudentID(ctx, studentID)
}

func (s *guardianService) GetStudentsByGuardian(ctx context.Context, userID uuid.UUID) ([]*domain.Student, error) {
	return s.repo.GetStudentsByUserID(ctx, userID)
}

func (s *guardianService) UnlinkGuardian(ctx context.Context, guardianID uuid.UUID) error {
	s.log.Info("Unlinking guardian", zap.String("id", guardianID.String()))

	if _, err := s.repo.GetByID(ctx, guardianID); err != nil {
		return fmt.Errorf("guardian not found: %w", err)
	}

	return s.repo.Delete(ctx, guardianID)
}

// GuardianModule provides the guardian service to the fx graph.
var GuardianModule = fx.Provide(NewGuardianService)
