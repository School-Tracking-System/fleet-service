package fleet

import (
	"context"
	"fmt"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/repositories"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/services"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type driverService struct {
	repo repositories.DriverRepository
	log  *zap.Logger
}

// NewDriverService creates a new business logic service for drivers.
func NewDriverService(repo repositories.DriverRepository, log *zap.Logger) services.DriverService {
	return &driverService{
		repo: repo,
		log:  log,
	}
}

func (s *driverService) RegisterDriver(ctx context.Context, req services.RegisterDriverRequest) (*domain.Driver, error) {
	s.log.Info("Registering driver", zap.String("user_id", req.UserID.String()))

	driver, err := domain.NewDriver(domain.NewDriverParams{
		UserID:         req.UserID,
		LicenseNumber:  req.LicenseNumber,
		LicenseType:    req.LicenseType,
		LicenseExpiry:  req.LicenseExpiry,
		CedulaID:       req.CedulaID,
		EmergencyPhone: req.EmergencyPhone,
	})
	if err != nil {
		return nil, fmt.Errorf("invalid driver data: %w", err)
	}

	if err := s.repo.Create(ctx, driver); err != nil {
		s.log.Error("Failed to persist driver", zap.Error(err))
		return nil, fmt.Errorf("failed to register driver: %w", err)
	}

	s.log.Info("Driver successfully registered", zap.String("id", driver.ID.String()))
	return driver, nil
}

func (s *driverService) UpdateDriver(ctx context.Context, req services.UpdateDriverRequest) (*domain.Driver, error) {
	s.log.Info("Updating driver", zap.String("id", req.ID.String()))

	driver, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("driver not found: %w", err)
	}

	driver.Apply(toDriverPatch(req))

	if err := s.repo.Update(ctx, driver); err != nil {
		s.log.Error("Failed to update driver", zap.Error(err))
		return nil, fmt.Errorf("failed to update driver: %w", err)
	}

	s.log.Info("Driver successfully updated", zap.String("id", driver.ID.String()))
	return driver, nil
}

func (s *driverService) GetDriver(ctx context.Context, id uuid.UUID) (*domain.Driver, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *driverService) GetDriverByUserID(ctx context.Context, userID uuid.UUID) (*domain.Driver, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *driverService) ListDrivers(ctx context.Context, limit, offset int) ([]*domain.Driver, int, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.List(ctx, limit, offset)
}

// toDriverPatch maps an UpdateDriverRequest to a domain DriverPatch.
func toDriverPatch(req services.UpdateDriverRequest) domain.DriverPatch {
	return domain.DriverPatch{
		LicenseType:    nonEmptyStr(req.LicenseType),
		LicenseExpiry:  req.LicenseExpiry,
		EmergencyPhone: nonEmptyStr(req.EmergencyPhone),
		Status:         nonEmptyDriverStatus(req.Status),
	}
}

func nonEmptyDriverStatus(s domain.DriverStatus) *domain.DriverStatus {
	if s == "" {
		return nil
	}
	return &s
}

// DriverModule provides the driver service to the fx dependency graph.
var DriverModule = fx.Provide(NewDriverService)
