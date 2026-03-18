package domain_test

import (
	"testing"
	"time"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewVehicle_Success(t *testing.T) {
	insExp := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	techExp := time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)

	v, err := domain.NewVehicle(domain.NewVehicleParams{
		Plate:         "ABC-1234",
		Brand:         "Toyota",
		Model:         "Hiace",
		Year:          2022,
		Capacity:      15,
		Color:         "White",
		VehicleType:   domain.VehicleTypeVan,
		ChassisNum:    "9FBSS12H4XBB70256",
		InsuranceExp:  &insExp,
		TechReviewExp: &techExp,
	})

	require.NoError(t, err)
	assert.NotNil(t, v)
	assert.NotEmpty(t, v.ID)
	assert.Equal(t, "ABC-1234", v.Plate)
	assert.Equal(t, "Toyota", v.Brand)
	assert.Equal(t, domain.VehicleStatusActive, v.Status)
	assert.Equal(t, "White", v.Color)
	assert.Equal(t, domain.VehicleTypeVan, v.VehicleType)
	assert.Equal(t, "9FBSS12H4XBB70256", v.ChassisNum)
	assert.Equal(t, &insExp, v.InsuranceExp)
	assert.Equal(t, &techExp, v.TechReviewExp)
}

func TestNewVehicle_EmptyPlate(t *testing.T) {
	_, err := domain.NewVehicle(domain.NewVehicleParams{
		Plate:    "",
		Year:     2022,
		Capacity: 15,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plate is required")
}

func TestNewVehicle_InvalidYear(t *testing.T) {
	_, err := domain.NewVehicle(domain.NewVehicleParams{
		Plate:    "XYZ-999",
		Year:     1800,
		Capacity: 10,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "year is invalid")
}

func TestNewVehicle_ZeroCapacity(t *testing.T) {
	_, err := domain.NewVehicle(domain.NewVehicleParams{
		Plate:    "XYZ-999",
		Year:     2022,
		Capacity: 0,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "capacity must be greater than 0")
}

func TestNewVehicle_NilableDatesAreOptional(t *testing.T) {
	v, err := domain.NewVehicle(domain.NewVehicleParams{
		Plate:    "MIN-001",
		Brand:    "Ford",
		Model:    "Transit",
		Year:     2020,
		Capacity: 12,
	})

	require.NoError(t, err)
	assert.Nil(t, v.InsuranceExp)
	assert.Nil(t, v.TechReviewExp)
}
