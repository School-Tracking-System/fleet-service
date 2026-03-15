package domain

import (
	"time"

	"github.com/google/uuid"
)

type VehicleStatus string

const (
	VehicleStatusActive      VehicleStatus = "active"
	VehicleStatusMaintenance VehicleStatus = "maintenance"
	VehicleStatusInactive    VehicleStatus = "inactive"
)

// Vehicle represents a physical transport unit in the fleet.
type Vehicle struct {
	ID        uuid.UUID
	Plate     string
	Brand     string
	Model     string
	Year      int
	Capacity  int
	Status    VehicleStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewVehicle is a factory function that creates a valid Vehicle instance.
func NewVehicle(plate, brand, model string, year, capacity int) (*Vehicle, error) {
	// Add business validations here (e.g., Year constraints, Capacity constraints)
	return &Vehicle{
		ID:        uuid.New(),
		Plate:     plate,
		Brand:     brand,
		Model:     model,
		Year:      year,
		Capacity:  capacity,
		Status:    VehicleStatusActive,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}
