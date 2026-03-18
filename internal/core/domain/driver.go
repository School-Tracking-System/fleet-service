package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// DriverStatus represents the operational state of a driver.
type DriverStatus string

const (
	DriverStatusActive    DriverStatus = "active"
	DriverStatusSuspended DriverStatus = "suspended"
	DriverStatusInactive  DriverStatus = "inactive"
)

// Driver represents a person authorized to operate a vehicle in the fleet.
// user_id is a logical FK to the Auth service's users table (no physical constraint).
type Driver struct {
	ID             uuid.UUID
	UserID         uuid.UUID    // Logical reference to Auth service user
	LicenseNumber  string
	LicenseType    string       // e.g. "B", "C", "D"
	LicenseExpiry  time.Time
	CedulaID       string       // National ID document
	EmergencyPhone string
	Status         DriverStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewDriverParams holds all data required to register a new Driver.
type NewDriverParams struct {
	UserID         uuid.UUID
	LicenseNumber  string
	LicenseType    string
	LicenseExpiry  time.Time
	CedulaID       string
	EmergencyPhone string
}

// NewDriver creates a valid Driver instance enforcing business invariants.
func NewDriver(p NewDriverParams) (*Driver, error) {
	if p.UserID == uuid.Nil {
		return nil, errors.New("user_id is required")
	}
	if p.LicenseNumber == "" {
		return nil, errors.New("license number is required")
	}
	if p.LicenseType == "" {
		return nil, errors.New("license type is required")
	}
	if p.LicenseExpiry.IsZero() {
		return nil, errors.New("license expiry date is required")
	}
	if p.LicenseExpiry.Before(time.Now()) {
		return nil, errors.New("license is already expired")
	}
	if p.CedulaID == "" {
		return nil, errors.New("cedula ID is required")
	}

	now := time.Now().UTC()
	return &Driver{
		ID:             uuid.New(),
		UserID:         p.UserID,
		LicenseNumber:  p.LicenseNumber,
		LicenseType:    p.LicenseType,
		LicenseExpiry:  p.LicenseExpiry,
		CedulaID:       p.CedulaID,
		EmergencyPhone: p.EmergencyPhone,
		Status:         DriverStatusActive,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// DriverPatch holds optional fields for a partial update of a Driver.
// Uses the Patch Object pattern: nil = no change, non-nil = apply.
type DriverPatch struct {
	LicenseType    *string
	LicenseExpiry  *time.Time
	EmergencyPhone *string
	Status         *DriverStatus
}

// Apply merges a DriverPatch into the Driver, updating only non-nil fields.
func (d *Driver) Apply(patch DriverPatch) {
	if patch.LicenseType != nil {
		d.LicenseType = *patch.LicenseType
	}
	if patch.LicenseExpiry != nil {
		d.LicenseExpiry = *patch.LicenseExpiry
	}
	if patch.EmergencyPhone != nil {
		d.EmergencyPhone = *patch.EmergencyPhone
	}
	if patch.Status != nil {
		d.Status = *patch.Status
	}
	d.UpdatedAt = time.Now().UTC()
}
