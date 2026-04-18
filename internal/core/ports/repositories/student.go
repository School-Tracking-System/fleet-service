package repositories

import (
	"context"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/google/uuid"
)

// StudentRepository defines the persistence contract for Student entities.
type StudentRepository interface {
	Create(ctx context.Context, student *domain.Student) error
	Update(ctx context.Context, student *domain.Student) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Student, error)
	GetByCedulaID(ctx context.Context, cedulaID string) (*domain.Student, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Student, int, error)
	ListBySchool(ctx context.Context, schoolID uuid.UUID, limit, offset int) ([]*domain.Student, int, error)
}
