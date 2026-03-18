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

func validRegisterRequest() services.RegisterDriverRequest {
	return services.RegisterDriverRequest{
		UserID:        uuid.New(),
		LicenseNumber: "DRV-TEST-001",
		LicenseType:   "C",
		LicenseExpiry: time.Now().Add(365 * 24 * time.Hour),
		CedulaID:      "1712345678",
	}
}

func TestRegisterDriver(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockDriverRepository)
		svc := NewDriverService(repo, log)

		req := validRegisterRequest()
		repo.On("Create", ctx, mock.MatchedBy(func(d *domain.Driver) bool {
			return d.LicenseNumber == req.LicenseNumber && d.Status == domain.DriverStatusActive
		})).Return(nil)

		driver, err := svc.RegisterDriver(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, driver)
		assert.Equal(t, req.LicenseNumber, driver.LicenseNumber)
		repo.AssertExpectations(t)
	})

	t.Run("invalid data returns domain error", func(t *testing.T) {
		repo := new(mocks.MockDriverRepository)
		svc := NewDriverService(repo, log)

		req := validRegisterRequest()
		req.LicenseNumber = ""

		_, err := svc.RegisterDriver(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid driver data")
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("duplicate driver returns repo error", func(t *testing.T) {
		repo := new(mocks.MockDriverRepository)
		svc := NewDriverService(repo, log)

		repo.On("Create", ctx, mock.Anything).Return(domain.ErrDuplicateDriver)

		_, err := svc.RegisterDriver(ctx, validRegisterRequest())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to register driver")
	})
}

func TestUpdateDriver(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()
	id := uuid.New()

	t.Run("success - status change", func(t *testing.T) {
		repo := new(mocks.MockDriverRepository)
		svc := NewDriverService(repo, log)

		existing := &domain.Driver{ID: id, LicenseType: "C", Status: domain.DriverStatusActive}
		repo.On("GetByID", ctx, id).Return(existing, nil)
		repo.On("Update", ctx, mock.MatchedBy(func(d *domain.Driver) bool {
			return d.Status == domain.DriverStatusSuspended
		})).Return(nil)

		updated, err := svc.UpdateDriver(ctx, services.UpdateDriverRequest{
			ID:     id,
			Status: domain.DriverStatusSuspended,
		})

		assert.NoError(t, err)
		assert.Equal(t, domain.DriverStatusSuspended, updated.Status)
		repo.AssertExpectations(t)
	})

	t.Run("driver not found", func(t *testing.T) {
		repo := new(mocks.MockDriverRepository)
		svc := NewDriverService(repo, log)

		repo.On("GetByID", ctx, id).Return(nil, errors.New("not found"))

		_, err := svc.UpdateDriver(ctx, services.UpdateDriverRequest{ID: id})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "driver not found")
	})
}

func TestGetDriverByUserID(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()
	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockDriverRepository)
		svc := NewDriverService(repo, log)

		expected := &domain.Driver{UserID: userID, LicenseNumber: "DRV-001"}
		repo.On("GetByUserID", ctx, userID).Return(expected, nil)

		driver, err := svc.GetDriverByUserID(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, expected, driver)
	})
}
