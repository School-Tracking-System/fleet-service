package fleet

import (
	"context"
	"testing"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/mocks"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func validRegisterStudentReq() services.RegisterStudentRequest {
	return services.RegisterStudentRequest{
		FirstName: "Ana",
		LastName:  "García",
		SchoolID:  uuid.New(),
	}
}

func TestRegisterStudent(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockStudentRepository)
		pub := new(mocks.MockEventPublisher)
		svc := NewStudentService(repo, pub, log)

		req := validRegisterStudentReq()
		repo.On("Create", ctx, mock.MatchedBy(func(s *domain.Student) bool {
			return s.FirstName == req.FirstName && s.IsActive
		})).Return(nil)
		pub.On("Publish", ctx, mock.AnythingOfType("string"), mock.Anything).Return(nil)

		student, err := svc.RegisterStudent(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, student)
		repo.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		repo := new(mocks.MockStudentRepository)
		pub := new(mocks.MockEventPublisher)
		svc := NewStudentService(repo, pub, log)

		req := validRegisterStudentReq()
		req.FirstName = ""

		_, err := svc.RegisterStudent(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid student data")
	})
}

func TestUpdateStudent_Service(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()
	id := uuid.New()

	repo := new(mocks.MockStudentRepository)
	pub := new(mocks.MockEventPublisher)
	svc := NewStudentService(repo, pub, log)

	existing := &domain.Student{ID: id, FirstName: "Old", LastName: "Name"}
	repo.On("GetByID", ctx, id).Return(existing, nil)
	repo.On("Update", ctx, mock.MatchedBy(func(s *domain.Student) bool {
		return s.FirstName == "New"
	})).Return(nil)

	updated, err := svc.UpdateStudent(ctx, services.UpdateStudentRequest{
		ID:        id,
		FirstName: "New",
	})

	assert.NoError(t, err)
	assert.Equal(t, "New", updated.FirstName)
}

func TestDeactivateStudent(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()
	id := uuid.New()

	repo := new(mocks.MockStudentRepository)
	pub := new(mocks.MockEventPublisher)
	svc := NewStudentService(repo, pub, log)

	existing := &domain.Student{ID: id, IsActive: true}
	repo.On("GetByID", ctx, id).Return(existing, nil)
	repo.On("Update", ctx, mock.MatchedBy(func(s *domain.Student) bool {
		return !s.IsActive
	})).Return(nil)

	err := svc.DeactivateStudent(ctx, id)

	assert.NoError(t, err)
}
