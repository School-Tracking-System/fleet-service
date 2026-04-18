package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// GuardianRelation describes the relationship between a guardian and a student.
type GuardianRelation string

const (
	GuardianRelationFather        GuardianRelation = "father"
	GuardianRelationMother        GuardianRelation = "mother"
	GuardianRelationLegalGuardian GuardianRelation = "legal_guardian"
	GuardianRelationOther         GuardianRelation = "other"
)

// Student represents a child registered in the school transportation system.
type Student struct {
	ID             uuid.UUID
	FirstName      string
	LastName       string
	Grade          string
	SchoolID       uuid.UUID  // FK → schools(id)
	PickupLocation *Location  // PostGIS point (reuses Location from school.go)
	PickupAddress  string
	PhotoURL       string
	IsActive       bool
	CedulaID       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// FullName returns the student's full name as a convenience method.
func (s *Student) FullName() string {
	return s.FirstName + " " + s.LastName
}

// NewStudentParams holds all data required to register a Student.
type NewStudentParams struct {
	FirstName      string
	LastName       string
	Grade          string
	SchoolID       uuid.UUID
	PickupLocation *Location
	PickupAddress  string
	PhotoURL       string
	CedulaID       string
}

// NewStudent creates a valid Student instance enforcing business invariants.
func NewStudent(p NewStudentParams) (*Student, error) {
	if p.FirstName == "" {
		return nil, errors.New("first name is required")
	}
	if p.LastName == "" {
		return nil, errors.New("last name is required")
	}
	if p.SchoolID == uuid.Nil {
		return nil, errors.New("school_id is required")
	}
	if p.CedulaID == "" {
		return nil, errors.New("cedula_id is required")
	}

	now := time.Now().UTC()
	return &Student{
		ID:             uuid.New(),
		FirstName:      p.FirstName,
		LastName:       p.LastName,
		Grade:          p.Grade,
		SchoolID:       p.SchoolID,
		PickupLocation: p.PickupLocation,
		PickupAddress:  p.PickupAddress,
		PhotoURL:       p.PhotoURL,
		IsActive:       true,
		CedulaID:       p.CedulaID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// StudentPatch holds optional fields for a partial update of a Student.
type StudentPatch struct {
	FirstName      *string
	LastName       *string
	Grade          *string
	PickupLocation *Location // nil = keep current; set to explicit zero struct to remove
	PickupAddress  *string
	PhotoURL       *string
	IsActive       *bool
	CedulaID       *string
}

// Apply merges a StudentPatch into the Student, updating only non-nil fields.
func (s *Student) Apply(patch StudentPatch) {
	if patch.FirstName != nil {
		s.FirstName = *patch.FirstName
	}
	if patch.LastName != nil {
		s.LastName = *patch.LastName
	}
	if patch.Grade != nil {
		s.Grade = *patch.Grade
	}
	if patch.PickupLocation != nil {
		s.PickupLocation = patch.PickupLocation
	}
	if patch.PickupAddress != nil {
		s.PickupAddress = *patch.PickupAddress
	}
	if patch.PhotoURL != nil {
		s.PhotoURL = *patch.PhotoURL
	}
	if patch.IsActive != nil {
		s.IsActive = *patch.IsActive
	}
	if patch.CedulaID != nil {
		s.CedulaID = *patch.CedulaID
	}
	s.UpdatedAt = time.Now().UTC()
}

// Guardian represents a parent or legal guardian of a Student.
type Guardian struct {
	ID        uuid.UUID
	UserID    uuid.UUID        // Logical FK → Auth service user
	StudentID uuid.UUID        // FK → students(id)
	Relation  GuardianRelation
	IsPrimary bool
	CreatedAt time.Time
}

// NewGuardianParams holds all data required to link a guardian to a student.
type NewGuardianParams struct {
	UserID    uuid.UUID
	StudentID uuid.UUID
	Relation  GuardianRelation
	IsPrimary bool
}

// NewGuardian creates a valid Guardian instance enforcing business invariants.
func NewGuardian(p NewGuardianParams) (*Guardian, error) {
	if p.UserID == uuid.Nil {
		return nil, errors.New("user_id is required")
	}
	if p.StudentID == uuid.Nil {
		return nil, errors.New("student_id is required")
	}
	relation := p.Relation
	if relation == "" {
		relation = GuardianRelationOther
	}

	return &Guardian{
		ID:        uuid.New(),
		UserID:    p.UserID,
		StudentID: p.StudentID,
		Relation:  relation,
		IsPrimary: p.IsPrimary,
		CreatedAt: time.Now().UTC(),
	}, nil
}
