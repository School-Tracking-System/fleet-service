package domain_test

import (
	"testing"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validRouteParams() domain.NewRouteParams {
	return domain.NewRouteParams{
		Name:         "Ruta Mañana - Norte",
		SchoolID:     uuid.New(),
		Direction:    domain.RouteDirectionToSchool,
		ScheduleTime: "06:30",
	}
}

func TestNewRoute_Success(t *testing.T) {
	p := validRouteParams()
	r, err := domain.NewRoute(p)

	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, p.Name, r.Name)
	assert.True(t, r.IsActive)
}

func TestRoute_ApplyPatch(t *testing.T) {
	r, _ := domain.NewRoute(validRouteParams())
	newName := "Ruta Editada"
	newActive := false

	r.Apply(domain.RoutePatch{
		Name:     &newName,
		IsActive: &newActive,
	})

	assert.Equal(t, "Ruta Editada", r.Name)
	assert.False(t, r.IsActive)
}

func TestNewRouteStop(t *testing.T) {
	routeID := uuid.New()
	studentID := uuid.New()
	p := domain.NewStopParams{
		RouteID:   routeID,
		StudentID: studentID,
		Order:     1,
		Location:  domain.Location{Longitude: -78.1, Latitude: -0.1},
		Address:   "Parada central",
		EstTime:   "06:45",
	}

	s, err := domain.NewRouteStop(p)

	require.NoError(t, err)
	assert.NotNil(t, s)
	assert.Equal(t, routeID, s.RouteID)
	assert.Equal(t, 1, s.Order)
}
