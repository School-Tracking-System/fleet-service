package domain_test

import (
	"testing"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validSchoolParams() domain.NewSchoolParams {
	return domain.NewSchoolParams{
		Name:    "Unidad Educativa Benalcázar",
		Address: "Av. República y Diego de Almagro, Quito",
		Location: &domain.Location{
			Longitude: -78.4966,
			Latitude:  -0.1865,
		},
		Phone: "+593-2-255-0000",
		Email: "info@benalcazar.edu.ec",
	}
}

func TestNewSchool_Success(t *testing.T) {
	p := validSchoolParams()
	s, err := domain.NewSchool(p)

	require.NoError(t, err)
	assert.NotNil(t, s)
	assert.NotEmpty(t, s.ID)
	assert.Equal(t, p.Name, s.Name)
	assert.Equal(t, p.Address, s.Address)
	assert.NotNil(t, s.Location)
	assert.InDelta(t, -78.4966, s.Location.Longitude, 0.0001)
	assert.InDelta(t, -0.1865, s.Location.Latitude, 0.0001)
}

func TestNewSchool_WithoutLocation(t *testing.T) {
	p := validSchoolParams()
	p.Location = nil

	s, err := domain.NewSchool(p)

	require.NoError(t, err)
	assert.Nil(t, s.Location)
}

func TestNewSchool_EmptyName(t *testing.T) {
	p := validSchoolParams()
	p.Name = ""

	_, err := domain.NewSchool(p)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "school name is required")
}

func TestNewSchool_EmptyAddress(t *testing.T) {
	p := validSchoolParams()
	p.Address = ""

	_, err := domain.NewSchool(p)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "school address is required")
}

func TestSchool_Apply_UpdatesOnlySetFields(t *testing.T) {
	s, _ := domain.NewSchool(validSchoolParams())
	originalAddress := s.Address
	newName := "Colegio Nacional Mejía"

	s.Apply(domain.SchoolPatch{
		Name: &newName,
	})

	assert.Equal(t, "Colegio Nacional Mejía", s.Name)
	assert.Equal(t, originalAddress, s.Address) // unchanged
}

func TestSchool_Apply_UpdatesLocation(t *testing.T) {
	s, _ := domain.NewSchool(validSchoolParams())
	newLoc := &domain.Location{Longitude: -78.5000, Latitude: -0.2000}

	s.Apply(domain.SchoolPatch{Location: newLoc})

	assert.Equal(t, -78.5000, s.Location.Longitude)
}
