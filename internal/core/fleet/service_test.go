package fleet

import (
	"context"
	"errors"
	"testing"

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
	
	t.Run("success", func(t *testing.T) {
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
			return v.Plate == "ABC-1234" && v.Brand == "Toyota"
		})).Return(nil)
		
		vehicle, err := svc.CreateVehicle(ctx, req)
		
		assert.NoError(t, err)
		assert.NotNil(t, vehicle)
		assert.Equal(t, "ABC-1234", vehicle.Plate)
		repo.AssertExpectations(t)
	})

	t.Run("empty plate", func(t *testing.T) {
		repo := new(mocks.MockVehicleRepository)
		svc := NewVehicleService(repo, log)
		
		req := services.CreateVehicleRequest{
			Plate: "",
		}
		
		vehicle, err := svc.CreateVehicle(ctx, req)
		
		assert.Error(t, err)
		assert.Nil(t, vehicle)
		assert.Contains(t, err.Error(), "plate is required")
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(mocks.MockVehicleRepository)
		svc := NewVehicleService(repo, log)
		
		req := services.CreateVehicleRequest{
			Plate:    "ABC-1234",
			Brand:    "Toyota",
			Year:     2024,
			Capacity: 15,
		}
		
		repo.On("Create", ctx, mock.Anything).Return(errors.New("db error"))
		
		vehicle, err := svc.CreateVehicle(ctx, req)
		
		assert.Error(t, err)
		assert.Nil(t, vehicle)
		assert.Contains(t, err.Error(), "failed to create vehicle in repository")
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
