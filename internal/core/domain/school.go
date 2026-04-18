package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Location represents a GPS coordinate (longitude, latitude).
// Maps to a PostGIS GEOMETRY(Point, 4326) column.
type Location struct {
	Longitude float64
	Latitude  float64
}

// School represents an educational institution served by the fleet.
type School struct {
	ID        uuid.UUID
	Name      string
	Address   string
	Location  *Location // nil if not yet geocoded
	Phone     string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewSchoolParams holds all data required to register a School.
type NewSchoolParams struct {
	Name     string
	Address  string
	Location *Location
	Phone    string
	Email    string
}

// NewSchool creates a valid School instance enforcing business invariants.
func NewSchool(p NewSchoolParams) (*School, error) {
	if p.Name == "" {
		return nil, errors.New("school name is required")
	}
	if p.Address == "" {
		return nil, errors.New("school address is required")
	}
	now := time.Now().UTC()
	return &School{
		ID:        uuid.New(),
		Name:      p.Name,
		Address:   p.Address,
		Location:  p.Location,
		Phone:     p.Phone,
		Email:     p.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// SchoolPatch holds optional fields for a partial update of a School.
type SchoolPatch struct {
	Name     *string
	Address  *string
	Location *Location
	Phone    *string
	Email    *string
}

// Apply merges a SchoolPatch into the School, updating only non-nil fields.
func (s *School) Apply(patch SchoolPatch) {
	if patch.Name != nil {
		s.Name = *patch.Name
	}
	if patch.Address != nil {
		s.Address = *patch.Address
	}
	if patch.Location != nil {
		s.Location = patch.Location
	}
	if patch.Phone != nil {
		s.Phone = *patch.Phone
	}
	if patch.Email != nil {
		s.Email = *patch.Email
	}
	s.UpdatedAt = time.Now().UTC()
}

// SchoolContact represents a staff member at a school.
type SchoolContact struct {
	ID        uuid.UUID
	SchoolID  uuid.UUID
	UserID    uuid.UUID // Logical FK → Auth service user
	Position  string
	IsActive  bool
	CreatedAt time.Time
}

// NewSchoolContactParams holds the data required to add a contact to a school.
type NewSchoolContactParams struct {
	SchoolID uuid.UUID
	UserID   uuid.UUID
	Position string
}

// NewSchoolContact creates a valid SchoolContact enforcing business invariants.
func NewSchoolContact(p NewSchoolContactParams) (*SchoolContact, error) {
	if p.SchoolID == uuid.Nil {
		return nil, errors.New("school_id is required")
	}
	if p.UserID == uuid.Nil {
		return nil, errors.New("user_id is required")
	}
	return &SchoolContact{
		ID:        uuid.New(),
		SchoolID:  p.SchoolID,
		UserID:    p.UserID,
		Position:  p.Position,
		IsActive:  true,
		CreatedAt: time.Now().UTC(),
	}, nil
}
