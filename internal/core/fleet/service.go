package fleet

import (
	"context"
	"fmt"
	"strings"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/repositories"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/services"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type vehicleService struct {
	repo repositories.VehicleRepository
	log  *zap.Logger
}

// NewVehicleService creates a new business logic service for vehicles.
func NewVehicleService(repo repositories.VehicleRepository, log *zap.Logger) services.VehicleService {
	return &vehicleService{
		repo: repo,
		log:  log,
	}
}

func (s *vehicleService) CreateVehicle(ctx context.Context, req services.CreateVehicleRequest) (*domain.Vehicle, error) {
	s.log.Info("Creating new vehicle", zap.String("plate", req.Plate))

	// Validate minimal business rules at the application layer constraints if needed
	if strings.TrimSpace(req.Plate) == "" {
		return nil, fmt.Errorf("plate is required")
	}

	vehicle, err := domain.NewVehicle(
		strings.ToUpper(req.Plate),
		req.Brand,
		req.Model,
		req.Year,
		req.Capacity,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid vehicle data: %w", err)
	}

	if err := s.repo.Create(ctx, vehicle); err != nil {
		s.log.Error("Failed to persist vehicle", zap.Error(err), zap.String("plate", vehicle.Plate))
		return nil, fmt.Errorf("failed to create vehicle in repository: %w", err)
	}

	s.log.Info("Vehicle successfully created", zap.String("id", vehicle.ID.String()))
	return vehicle, nil
}

func (s *vehicleService) GetVehicle(ctx context.Context, id uuid.UUID) (*domain.Vehicle, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *vehicleService) ListVehicles(ctx context.Context, limit, offset int) ([]*domain.Vehicle, int, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.List(ctx, limit, offset)
}

// Module provides the fleet core service to fx graph.
var Module = fx.Provide(NewVehicleService)
