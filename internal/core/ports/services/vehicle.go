package services

import (
	"context"
	"time"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/google/uuid"
)

// VehicleService defines the core business logic contract for Vehicle operations.
type VehicleService interface {
	CreateVehicle(ctx context.Context, req CreateVehicleRequest) (*domain.Vehicle, error)
	UpdateVehicle(ctx context.Context, req UpdateVehicleRequest) (*domain.Vehicle, error)
	GetVehicle(ctx context.Context, id uuid.UUID) (*domain.Vehicle, error)
	ListVehicles(ctx context.Context, limit, offset int) ([]*domain.Vehicle, int, error)
}

// CreateVehicleRequest encapsulates the data required to create a new vehicle.
type CreateVehicleRequest struct {
	Plate         string
	Brand         string
	Model         string
	Year          int
	Capacity      int
	Color         string
	VehicleType   domain.VehicleType
	ChassisNum    string
	InsuranceExp  *time.Time
	TechReviewExp *time.Time
}

// UpdateVehicleRequest encapsulates the data required to update an existing vehicle.
type UpdateVehicleRequest struct {
	ID            uuid.UUID
	Brand         string
	Model         string
	Year          int
	Capacity      int
	Status        domain.VehicleStatus
	Color         string
	VehicleType   domain.VehicleType
	ChassisNum    string
	InsuranceExp  *time.Time
	TechReviewExp *time.Time
}
