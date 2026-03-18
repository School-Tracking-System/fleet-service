package domain_test

import (
	"testing"
	"time"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validDriverParams() domain.NewDriverParams {
	return domain.NewDriverParams{
		UserID:         uuid.New(),
		LicenseNumber:  "DRV-001",
		LicenseType:    "C",
		LicenseExpiry:  time.Now().Add(365 * 24 * time.Hour),
		CedulaID:       "1712345678",
		EmergencyPhone: "+593987654321",
	}
}

func TestNewDriver_Success(t *testing.T) {
	p := validDriverParams()
	d, err := domain.NewDriver(p)

	require.NoError(t, err)
	assert.NotNil(t, d)
	assert.NotEqual(t, uuid.Nil, d.ID)
	assert.Equal(t, p.UserID, d.UserID)
	assert.Equal(t, p.LicenseNumber, d.LicenseNumber)
	assert.Equal(t, "C", d.LicenseType)
	assert.Equal(t, domain.DriverStatusActive, d.Status)
	assert.Equal(t, "+593987654321", d.EmergencyPhone)
}

func TestNewDriver_EmptyUserID(t *testing.T) {
	p := validDriverParams()
	p.UserID = uuid.Nil

	_, err := domain.NewDriver(p)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user_id is required")
}

func TestNewDriver_EmptyLicenseNumber(t *testing.T) {
	p := validDriverParams()
	p.LicenseNumber = ""

	_, err := domain.NewDriver(p)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "license number is required")
}

func TestNewDriver_EmptyLicenseType(t *testing.T) {
	p := validDriverParams()
	p.LicenseType = ""

	_, err := domain.NewDriver(p)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "license type is required")
}

func TestNewDriver_ExpiredLicense(t *testing.T) {
	p := validDriverParams()
	p.LicenseExpiry = time.Now().Add(-24 * time.Hour) // yesterday

	_, err := domain.NewDriver(p)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "license is already expired")
}

func TestNewDriver_EmptyCedula(t *testing.T) {
	p := validDriverParams()
	p.CedulaID = ""

	_, err := domain.NewDriver(p)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cedula ID is required")
}

func TestDriver_Apply_PartialUpdate(t *testing.T) {
	d, _ := domain.NewDriver(validDriverParams())
	originalExpiry := d.LicenseExpiry

	newStatus := domain.DriverStatusSuspended
	d.Apply(domain.DriverPatch{
		Status: &newStatus,
	})

	assert.Equal(t, domain.DriverStatusSuspended, d.Status)
	assert.Equal(t, originalExpiry, d.LicenseExpiry) // unchanged
}

func TestDriver_Apply_UpdateExpiry(t *testing.T) {
	d, _ := domain.NewDriver(validDriverParams())
	newExpiry := time.Now().Add(700 * 24 * time.Hour)

	d.Apply(domain.DriverPatch{
		LicenseExpiry: &newExpiry,
	})

	assert.Equal(t, newExpiry, d.LicenseExpiry)
}
