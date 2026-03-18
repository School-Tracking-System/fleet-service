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

func TestLinkGuardian(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()
	repo := new(mocks.MockGuardianRepository)
	svc := NewGuardianService(repo, log)

	req := services.LinkGuardianRequest{
		UserID:    uuid.New(),
		StudentID: uuid.New(),
		Relation:  domain.GuardianRelationMother,
	}

	repo.On("Create", ctx, mock.MatchedBy(func(g *domain.Guardian) bool {
		return g.UserID == req.UserID && g.StudentID == req.StudentID
	})).Return(nil)

	guardian, err := svc.LinkGuardian(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, guardian)
	repo.AssertExpectations(t)
}

func TestUnlinkGuardian(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()
	repo := new(mocks.MockGuardianRepository)
	svc := NewGuardianService(repo, log)
	id := uuid.New()

	repo.On("GetByID", ctx, id).Return(&domain.Guardian{ID: id}, nil)
	repo.On("Delete", ctx, id).Return(nil)

	err := svc.UnlinkGuardian(ctx, id)

	assert.NoError(t, err)
}

func TestGetStudentsByGuardian_Service(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()
	repo := new(mocks.MockGuardianRepository)
	svc := NewGuardianService(repo, log)
	userID := uuid.New()

	expected := []*domain.Student{{ID: uuid.New(), FirstName: "Hijo"}}
	repo.On("GetStudentsByUserID", ctx, userID).Return(expected, nil)

	students, err := svc.GetStudentsByGuardian(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expected, students)
}
