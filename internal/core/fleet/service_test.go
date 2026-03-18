package fleet

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/mocks"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestCreateVehicle(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()

	t.Run("success with minimal fields", func(t *testing.T) {
		repo := new(mocks.MockVehicleRepository)
		svc := NewVehicleService(repo, log)

		req := services.CreateVehicleRequest{
			Plate:    "ABC-1234",
			Brand:    "Toyota",
			Model:    "Hiace",
			Year:     2024,
			Capacity: 15,
		}

		repo.On("Create", ctx, mock.MatchedBy(func(v *domain.Vehicle) bool {
			return v.Plate == "ABC-1234" && v.Brand == "Toyota" && v.Status == domain.VehicleStatusActive
		})).Return(nil)

		vehicle, err := svc.CreateVehicle(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, vehicle)
		assert.Equal(t, "ABC-1234", vehicle.Plate)
		repo.AssertExpectations(t)
	})

	t.Run("success with all extended fields", func(t *testing.T) {
		repo := new(mocks.MockVehicleRepository)
		svc := NewVehicleService(repo, log)

		insExp := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
		req := services.CreateVehicleRequest{
			Plate:        "VAN-001",
			Brand:        "Mercedes",
			Model:        "Sprinter",
			Year:         2023,
			Capacity:     20,
			Color:        "Blue",
			VehicleType:  domain.VehicleTypeVan,
			ChassisNum:   "WDB9066351L123456",
			InsuranceExp: &insExp,
		}

		repo.On("Create", ctx, mock.MatchedBy(func(v *domain.Vehicle) bool {
			return v.Plate == "VAN-001" && v.Color == "Blue" && v.VehicleType == domain.VehicleTypeVan
		})).Return(nil)

		vehicle, err := svc.CreateVehicle(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, "Blue", vehicle.Color)
		assert.Equal(t, domain.VehicleTypeVan, vehicle.VehicleType)
		assert.Equal(t, &insExp, vehicle.InsuranceExp)
		repo.AssertExpectations(t)
	})

	t.Run("plate is normalized to uppercase", func(t *testing.T) {
		repo := new(mocks.MockVehicleRepository)
		svc := NewVehicleService(repo, log)

		req := services.CreateVehicleRequest{
			Plate:    "abc-1234",
			Year:     2022,
			Capacity: 10,
		}

		repo.On("Create", ctx, mock.MatchedBy(func(v *domain.Vehicle) bool {
			return v.Plate == "ABC-1234"
		})).Return(nil)

		vehicle, err := svc.CreateVehicle(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, "ABC-1234", vehicle.Plate)
		repo.AssertExpectations(t)
	})

	t.Run("empty plate fails domain validation", func(t *testing.T) {
		repo := new(mocks.MockVehicleRepository)
		svc := NewVehicleService(repo, log)

		_, err := svc.CreateVehicle(ctx, services.CreateVehicleRequest{Plate: ""})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid vehicle data")
	})

	t.Run("repo error is propagated", func(t *testing.T) {
		repo := new(mocks.MockVehicleRepository)
		svc := NewVehicleService(repo, log)

		req := services.CreateVehicleRequest{Plate: "ERR-001", Year: 2022, Capacity: 10}
		repo.On("Create", ctx, mock.Anything).Return(errors.New("db error"))

		_, err := svc.CreateVehicle(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create vehicle in repository")
	})
}

func TestUpdateVehicle(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()
	id := uuid.New()

	t.Run("success partial update", func(t *testing.T) {
		repo := new(mocks.MockVehicleRepository)
		svc := NewVehicleService(repo, log)

		existing := &domain.Vehicle{ID: id, Plate: "OLD-001", Brand: "Toyota", Year: 2020, Capacity: 10, Status: domain.VehicleStatusActive}
		repo.On("GetByID", ctx, id).Return(existing, nil)
		repo.On("Update", ctx, mock.MatchedBy(func(v *domain.Vehicle) bool {
			return v.Brand == "Ford" && v.Status == domain.VehicleStatusMaintenance
		})).Return(nil)

		updated, err := svc.UpdateVehicle(ctx, services.UpdateVehicleRequest{
			ID:     id,
			Brand:  "Ford",
			Status: domain.VehicleStatusMaintenance,
		})

		assert.NoError(t, err)
		assert.Equal(t, "Ford", updated.Brand)
		assert.Equal(t, domain.VehicleStatusMaintenance, updated.Status)
		repo.AssertExpectations(t)
	})

	t.Run("vehicle not found", func(t *testing.T) {
		repo := new(mocks.MockVehicleRepository)
		svc := NewVehicleService(repo, log)

		repo.On("GetByID", ctx, id).Return(nil, domain.ErrVehicleNotFound)

		_, err := svc.UpdateVehicle(ctx, services.UpdateVehicleRequest{ID: id})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "vehicle not found")
	})
}

func TestGetVehicle(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()
	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockVehicleRepository)
		svc := NewVehicleService(repo, log)

		expected := &domain.Vehicle{ID: id, Plate: "XYZ-789"}
		repo.On("GetByID", ctx, id).Return(expected, nil)

		vehicle, err := svc.GetVehicle(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, expected, vehicle)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(mocks.MockVehicleRepository)
		svc := NewVehicleService(repo, log)

		repo.On("GetByID", ctx, id).Return(nil, errors.New("not found"))

		vehicle, err := svc.GetVehicle(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, vehicle)
	})
}
