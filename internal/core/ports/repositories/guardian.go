package repositories

import (
	"context"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/google/uuid"
)

// GuardianRepository defines the persistence contract for Guardian entities.
type GuardianRepository interface {
	Create(ctx context.Context, guardian *domain.Guardian) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Guardian, error)
	GetByStudentID(ctx context.Context, studentID uuid.UUID) ([]*domain.Guardian, error)
	GetStudentsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Student, error)
}
