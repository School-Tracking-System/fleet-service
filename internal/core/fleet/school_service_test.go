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

func validCreateSchoolReq() services.CreateSchoolRequest {
	return services.CreateSchoolRequest{
		Name:    "Unidad Educativa Benalcázar",
		Address: "Av. República y Diego de Almagro, Quito",
		Location: &domain.Location{
			Longitude: -78.4966,
			Latitude:  -0.1865,
		},
	}
}

func TestCreateSchool(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()

	t.Run("success with location", func(t *testing.T) {
		repo := new(mocks.MockSchoolRepository)
		svc := NewSchoolService(repo, log)

		req := validCreateSchoolReq()
		repo.On("Create", ctx, mock.MatchedBy(func(s *domain.School) bool {
			return s.Name == req.Name && s.Location != nil
		})).Return(nil)

		school, err := svc.CreateSchool(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, school)
		assert.Equal(t, req.Name, school.Name)
		assert.NotNil(t, school.Location)
		repo.AssertExpectations(t)
	})

	t.Run("success without location", func(t *testing.T) {
		repo := new(mocks.MockSchoolRepository)
		svc := NewSchoolService(repo, log)

		req := validCreateSchoolReq()
		req.Location = nil
		repo.On("Create", ctx, mock.MatchedBy(func(s *domain.School) bool {
			return s.Location == nil
		})).Return(nil)

		school, err := svc.CreateSchool(ctx, req)

		assert.NoError(t, err)
		assert.Nil(t, school.Location)
	})

	t.Run("empty name fails domain validation", func(t *testing.T) {
		repo := new(mocks.MockSchoolRepository)
		svc := NewSchoolService(repo, log)

		req := validCreateSchoolReq()
		req.Name = ""

		_, err := svc.CreateSchool(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid school data")
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("repo error is propagated", func(t *testing.T) {
		repo := new(mocks.MockSchoolRepository)
		svc := NewSchoolService(repo, log)

		repo.On("Create", ctx, mock.Anything).Return(errors.New("db error"))

		_, err := svc.CreateSchool(ctx, validCreateSchoolReq())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create school")
	})
}

func TestUpdateSchool(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()
	id := uuid.New()

	t.Run("success - update name and location", func(t *testing.T) {
		repo := new(mocks.MockSchoolRepository)
		svc := NewSchoolService(repo, log)

		existing := &domain.School{ID: id, Name: "Old Name", Address: "Old Address"}
		repo.On("GetByID", ctx, id).Return(existing, nil)
		repo.On("Update", ctx, mock.MatchedBy(func(s *domain.School) bool {
			return s.Name == "New Name" && s.Location != nil
		})).Return(nil)

		updated, err := svc.UpdateSchool(ctx, services.UpdateSchoolRequest{
			ID:       id,
			Name:     "New Name",
			Location: &domain.Location{Longitude: -78.5, Latitude: -0.2},
		})

		assert.NoError(t, err)
		assert.Equal(t, "New Name", updated.Name)
		assert.NotNil(t, updated.Location)
		repo.AssertExpectations(t)
	})

	t.Run("school not found", func(t *testing.T) {
		repo := new(mocks.MockSchoolRepository)
		svc := NewSchoolService(repo, log)

		repo.On("GetByID", ctx, id).Return(nil, errors.New("not found"))

		_, err := svc.UpdateSchool(ctx, services.UpdateSchoolRequest{ID: id})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "school not found")
	})
}

func TestGetSchool(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()
	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockSchoolRepository)
		svc := NewSchoolService(repo, log)

		expected := &domain.School{ID: id, Name: "Test School"}
		repo.On("GetByID", ctx, id).Return(expected, nil)

		school, err := svc.GetSchool(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, expected, school)
	})
}
