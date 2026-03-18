package services

import (
	"context"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/google/uuid"
)

// GuardianService defines the core business logic contract for Guardian operations.
type GuardianService interface {
	LinkGuardian(ctx context.Context, req LinkGuardianRequest) (*domain.Guardian, error)
	GetGuardiansByStudent(ctx context.Context, studentID uuid.UUID) ([]*domain.Guardian, error)
	GetStudentsByGuardian(ctx context.Context, userID uuid.UUID) ([]*domain.Student, error)
	UnlinkGuardian(ctx context.Context, guardianID uuid.UUID) error
}

// LinkGuardianRequest encapsulates the data to link a guardian to a student.
type LinkGuardianRequest struct {
	UserID    uuid.UUID
	StudentID uuid.UUID
	Relation  domain.GuardianRelation
	IsPrimary bool
}
