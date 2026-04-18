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

type schoolContactService struct {
	repo repositories.SchoolContactRepository
	log  *zap.Logger
}

// NewSchoolContactService creates a new business logic service for school contacts.
func NewSchoolContactService(repo repositories.SchoolContactRepository, log *zap.Logger) services.SchoolContactService {
	return &schoolContactService{repo: repo, log: log}
}

func (s *schoolContactService) AddContact(ctx context.Context, req services.AddContactRequest) (*domain.SchoolContact, error) {
	s.log.Info("Adding school contact",
		zap.String("school_id", req.SchoolID.String()),
		zap.String("user_id", req.UserID.String()),
	)

	contact, err := domain.NewSchoolContact(domain.NewSchoolContactParams{
		SchoolID: req.SchoolID,
		UserID:   req.UserID,
		Position: req.Position,
	})
	if err != nil {
		return nil, fmt.Errorf("invalid contact data: %w", err)
	}

	if err := s.repo.Create(ctx, contact); err != nil {
		s.log.Error("Failed to persist school contact", zap.Error(err))
		return nil, fmt.Errorf("failed to add contact: %w", err)
	}

	s.log.Info("School contact added", zap.String("id", contact.ID.String()))
	return contact, nil
}

func (s *schoolContactService) RemoveContact(ctx context.Context, contactID uuid.UUID) error {
	s.log.Info("Removing school contact", zap.String("id", contactID.String()))

	if _, err := s.repo.GetByID(ctx, contactID); err != nil {
		return fmt.Errorf("contact not found: %w", err)
	}

	return s.repo.Delete(ctx, contactID)
}

func (s *schoolContactService) ListContacts(ctx context.Context, schoolID uuid.UUID) ([]*domain.SchoolContact, error) {
	return s.repo.GetBySchoolID(ctx, schoolID)
}

// SchoolContactModule provides the school contact service to the fx graph.
var SchoolContactModule = fx.Provide(NewSchoolContactService)
