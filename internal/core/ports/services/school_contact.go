package services

import (
	"context"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/google/uuid"
)

// SchoolContactService defines business logic for managing school staff contacts.
type SchoolContactService interface {
	AddContact(ctx context.Context, req AddContactRequest) (*domain.SchoolContact, error)
	RemoveContact(ctx context.Context, contactID uuid.UUID) error
	ListContacts(ctx context.Context, schoolID uuid.UUID) ([]*domain.SchoolContact, error)
}

// AddContactRequest holds the data needed to add a contact to a school.
type AddContactRequest struct {
	SchoolID uuid.UUID
	UserID   uuid.UUID
	Position string
}
