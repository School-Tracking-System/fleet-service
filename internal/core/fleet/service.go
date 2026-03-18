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

	vehicle, err := domain.NewVehicle(domain.NewVehicleParams{
		Plate:         strings.ToUpper(strings.TrimSpace(req.Plate)),
		Brand:         req.Brand,
		Model:         req.Model,
		Year:          req.Year,
		Capacity:      req.Capacity,
		Color:         req.Color,
		VehicleType:   req.VehicleType,
		ChassisNum:    req.ChassisNum,
		InsuranceExp:  req.InsuranceExp,
		TechReviewExp: req.TechReviewExp,
	})
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

func (s *vehicleService) UpdateVehicle(ctx context.Context, req services.UpdateVehicleRequest) (*domain.Vehicle, error) {
	s.log.Info("Updating vehicle", zap.String("id", req.ID.String()))

	vehicle, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("vehicle not found: %w", err)
	}

	vehicle.Apply(toPatch(req))

	if err := s.repo.Update(ctx, vehicle); err != nil {
		s.log.Error("Failed to update vehicle", zap.Error(err))
		return nil, fmt.Errorf("failed to update vehicle in repository: %w", err)
	}

	s.log.Info("Vehicle successfully updated", zap.String("id", vehicle.ID.String()))
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

// toPatch maps an UpdateVehicleRequest to a domain VehiclePatch.
// It uses helper functions to convert zero-values to nil, so Apply only
// touches the fields the caller explicitly provided.
func toPatch(req services.UpdateVehicleRequest) domain.VehiclePatch {
	return domain.VehiclePatch{
		Brand:         nonEmptyStr(req.Brand),
		Model:         nonEmptyStr(req.Model),
		Year:          nonZeroInt(req.Year),
		Capacity:      nonZeroInt(req.Capacity),
		Status:        nonEmptyStatus(req.Status),
		Color:         nonEmptyStr(req.Color),
		VehicleType:   nonEmptyType(req.VehicleType),
		ChassisNum:    nonEmptyStr(req.ChassisNum),
		InsuranceExp:  req.InsuranceExp,
		TechReviewExp: req.TechReviewExp,
	}
}

func nonEmptyStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func nonZeroInt(n int) *int {
	if n == 0 {
		return nil
	}
	return &n
}

func nonEmptyStatus(s domain.VehicleStatus) *domain.VehicleStatus {
	if s == "" {
		return nil
	}
	return &s
}

func nonEmptyType(t domain.VehicleType) *domain.VehicleType {
	if t == "" {
		return nil
	}
	return &t
}

// Module provides the fleet core service to fx graph.
var Module = fx.Provide(NewVehicleService)
