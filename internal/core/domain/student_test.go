package domain_test

import (
	"testing"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validStudentParams() domain.NewStudentParams {
	return domain.NewStudentParams{
		FirstName: "Juan",
		LastName:  "Pérez",
		Grade:     "5to Grado",
		SchoolID:  uuid.New(),
		PickupLocation: &domain.Location{
			Longitude: -78.5,
			Latitude:  -0.2,
		},
		PickupAddress: "Calle de los Olivos 123",
	}
}

func TestNewStudent_Success(t *testing.T) {
	p := validStudentParams()
	s, err := domain.NewStudent(p)

	require.NoError(t, err)
	assert.NotNil(t, s)
	assert.NotEmpty(t, s.ID)
	assert.Equal(t, p.FirstName, s.FirstName)
	assert.Equal(t, "Juan Pérez", s.FullName())
	assert.True(t, s.IsActive)
}

func TestNewStudent_Validation(t *testing.T) {
	t.Run("empty first name", func(t *testing.T) {
		p := validStudentParams()
		p.FirstName = ""
		_, err := domain.NewStudent(p)
		assert.Error(t, err)
	})

	t.Run("empty last name", func(t *testing.T) {
		p := validStudentParams()
		p.LastName = ""
		_, err := domain.NewStudent(p)
		assert.Error(t, err)
	})

	t.Run("empty school id", func(t *testing.T) {
		p := validStudentParams()
		p.SchoolID = uuid.Nil
		_, err := domain.NewStudent(p)
		assert.Error(t, err)
	})
}

func TestStudent_ApplyPatch(t *testing.T) {
	s, _ := domain.NewStudent(validStudentParams())
	newName := "Carlos"
	newActive := false

	s.Apply(domain.StudentPatch{
		FirstName: &newName,
		IsActive:  &newActive,
	})

	assert.Equal(t, "Carlos", s.FirstName)
	assert.False(t, s.IsActive)
	assert.Equal(t, "Carlos Pérez", s.FullName())
}

func TestNewGuardian_Success(t *testing.T) {
	userID := uuid.New()
	studentID := uuid.New()
	
	g, err := domain.NewGuardian(domain.NewGuardianParams{
		UserID:    userID,
		StudentID: studentID,
		Relation:  domain.GuardianRelationFather,
		IsPrimary: true,
	})

	require.NoError(t, err)
	assert.NotNil(t, g)
	assert.Equal(t, userID, g.UserID)
	assert.Equal(t, studentID, g.StudentID)
	assert.Equal(t, domain.GuardianRelationFather, g.Relation)
	assert.True(t, g.IsPrimary)
}
