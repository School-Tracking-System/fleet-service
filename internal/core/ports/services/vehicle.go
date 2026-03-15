package services

import (
	"context"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/google/uuid"
)

// VehicleService defines the core business logic contract for Vehicle operations.
type VehicleService interface {
	CreateVehicle(ctx context.Context, req CreateVehicleRequest) (*domain.Vehicle, error)
	GetVehicle(ctx context.Context, id uuid.UUID) (*domain.Vehicle, error)
	ListVehicles(ctx context.Context, limit, offset int) ([]*domain.Vehicle, int, error)
}

// CreateVehicleRequest encapsulates the data required to create a new vehicle.
// This is specific to the Application layer (Service), decoupling it from gRPC/HTTP payloads.
type CreateVehicleRequest struct {
	Plate    string
	Brand    string
	Model    string
	Year     int
	Capacity int
}
