package handlers

import (
	"context"
	"testing"
	"time"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/mocks"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/services"
	pb "github.com/fercho/school-tracking/proto/gen/fleet/v1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCreateVehicleHandler(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()
	mockSvc := new(mocks.MockVehicleService)
	handler := NewVehicleHandler(mockSvc, log)

	t.Run("success", func(t *testing.T) {
		req := &pb.CreateVehicleRequest{
			Plate:    "ABC-1234",
			Brand:    "Toyota",
			Model:    "Hiace",
			Year:     2024,
			Capacity: 15,
		}

		vID := uuid.New()
		expectedVehicle := &domain.Vehicle{
			ID:        vID,
			Plate:     "ABC-1234",
			Brand:     "Toyota",
			Model:     "Hiace",
			Year:      2024,
			Capacity:  15,
			Status:    domain.VehicleStatusActive,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockSvc.On("CreateVehicle", ctx, services.CreateVehicleRequest{
			Plate:    "ABC-1234",
			Brand:    "Toyota",
			Model:    "Hiace",
			Year:     2024,
			Capacity: 15,
		}).Return(expectedVehicle, nil)

		resp, err := handler.CreateVehicle(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, vID.String(), resp.Vehicle.Id)
		assert.Equal(t, "ABC-1234", resp.Vehicle.Plate)
		mockSvc.AssertExpectations(t)
	})
}
