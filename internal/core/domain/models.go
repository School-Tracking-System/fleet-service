package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type VehicleStatus string
type VehicleType string

const (
	VehicleStatusActive      VehicleStatus = "active"
	VehicleStatusMaintenance VehicleStatus = "maintenance"
	VehicleStatusInactive    VehicleStatus = "inactive"

	VehicleTypeVan     VehicleType = "van"
	VehicleTypeBus     VehicleType = "bus"
	VehicleTypeMinibus VehicleType = "minibus"
)

// Vehicle represents a physical transport unit in the fleet.
type Vehicle struct {
	ID            uuid.UUID
	Plate         string
	Brand         string
	Model         string
	Year          int
	Capacity      int
	Status        VehicleStatus
	Color         string
	VehicleType   VehicleType
	ChassisNum    string
	InsuranceExp  *time.Time // nullable
	TechReviewExp *time.Time // nullable
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewVehicleParams holds all parameters needed to create a Vehicle.
type NewVehicleParams struct {
	Plate         string
	Brand         string
	Model         string
	Year          int
	Capacity      int
	Color         string
	VehicleType   VehicleType
	ChassisNum    string
	InsuranceExp  *time.Time
	TechReviewExp *time.Time
}

// NewVehicle is a factory function that creates a valid, new Vehicle instance.
// It enforces minimal business invariants.
func NewVehicle(p NewVehicleParams) (*Vehicle, error) {
	if p.Plate == "" {
		return nil, errors.New("plate is required")
	}
	if p.Year < 1900 || p.Year > time.Now().Year()+1 {
		return nil, errors.New("vehicle year is invalid")
	}
	if p.Capacity <= 0 {
		return nil, errors.New("capacity must be greater than 0")
	}
	now := time.Now().UTC()
	return &Vehicle{
		ID:            uuid.New(),
		Plate:         p.Plate,
		Brand:         p.Brand,
		Model:         p.Model,
		Year:          p.Year,
		Capacity:      p.Capacity,
		Status:        VehicleStatusActive,
		Color:         p.Color,
		VehicleType:   p.VehicleType,
		ChassisNum:    p.ChassisNum,
		InsuranceExp:  p.InsuranceExp,
		TechReviewExp: p.TechReviewExp,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// VehiclePatch holds optional fields to be applied to an existing Vehicle.
// Only non-nil pointer fields are applied. This is the canonical pattern for
// partial updates in DDD: the domain owns its own mutation logic.
type VehiclePatch struct {
	Brand         *string
	Model         *string
	Year          *int
	Capacity      *int
	Status        *VehicleStatus
	Color         *string
	VehicleType   *VehicleType
	ChassisNum    *string
	InsuranceExp  *time.Time
	TechReviewExp *time.Time
}

// Apply merges a VehiclePatch into the Vehicle, updating only the fields
// that are explicitly set (non-nil). All mutation logic stays in the domain.
func (v *Vehicle) Apply(patch VehiclePatch) {
	if patch.Brand != nil {
		v.Brand = *patch.Brand
	}
	if patch.Model != nil {
		v.Model = *patch.Model
	}
	if patch.Year != nil {
		v.Year = *patch.Year
	}
	if patch.Capacity != nil {
		v.Capacity = *patch.Capacity
	}
	if patch.Status != nil {
		v.Status = *patch.Status
	}
	if patch.Color != nil {
		v.Color = *patch.Color
	}
	if patch.VehicleType != nil {
		v.VehicleType = *patch.VehicleType
	}
	if patch.ChassisNum != nil {
		v.ChassisNum = *patch.ChassisNum
	}
	if patch.InsuranceExp != nil {
		v.InsuranceExp = patch.InsuranceExp
	}
	if patch.TechReviewExp != nil {
		v.TechReviewExp = patch.TechReviewExp
	}
	v.UpdatedAt = time.Now().UTC()
}
