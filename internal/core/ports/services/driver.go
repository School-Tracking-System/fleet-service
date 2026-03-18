package services

import (
	"context"
	"time"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/google/uuid"
)

// DriverService defines the core business logic contract for Driver operations.
type DriverService interface {
	RegisterDriver(ctx context.Context, req RegisterDriverRequest) (*domain.Driver, error)
	UpdateDriver(ctx context.Context, req UpdateDriverRequest) (*domain.Driver, error)
	GetDriver(ctx context.Context, id uuid.UUID) (*domain.Driver, error)
	GetDriverByUserID(ctx context.Context, userID uuid.UUID) (*domain.Driver, error)
	ListDrivers(ctx context.Context, limit, offset int) ([]*domain.Driver, int, error)
}

// RegisterDriverRequest encapsulates the data required to register a new driver.
type RegisterDriverRequest struct {
	UserID         uuid.UUID
	LicenseNumber  string
	LicenseType    string
	LicenseExpiry  time.Time
	CedulaID       string
	EmergencyPhone string
}

// UpdateDriverRequest encapsulates the data for a partial driver update.
type UpdateDriverRequest struct {
	ID             uuid.UUID
	LicenseType    string
	LicenseExpiry  *time.Time
	EmergencyPhone string
	Status         domain.DriverStatus
}
