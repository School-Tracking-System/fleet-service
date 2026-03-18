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

type studentService struct {
	repo repositories.StudentRepository
	log  *zap.Logger
}

// NewStudentService creates a new business logic service for students.
func NewStudentService(repo repositories.StudentRepository, log *zap.Logger) services.StudentService {
	return &studentService{repo: repo, log: log}
}

func (s *studentService) RegisterStudent(ctx context.Context, req services.RegisterStudentRequest) (*domain.Student, error) {
	s.log.Info("Registering student", zap.String("name", req.FirstName+" "+req.LastName))

	student, err := domain.NewStudent(domain.NewStudentParams{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Grade:          req.Grade,
		SchoolID:       req.SchoolID,
		PickupLocation: req.PickupLocation,
		PickupAddress:  req.PickupAddress,
		PhotoURL:       req.PhotoURL,
	})
	if err != nil {
		return nil, fmt.Errorf("invalid student data: %w", err)
	}

	if err := s.repo.Create(ctx, student); err != nil {
		s.log.Error("Failed to persist student", zap.Error(err))
		return nil, fmt.Errorf("failed to register student: %w", err)
	}

	s.log.Info("Student registered", zap.String("id", student.ID.String()))
	return student, nil
}

func (s *studentService) UpdateStudent(ctx context.Context, req services.UpdateStudentRequest) (*domain.Student, error) {
	s.log.Info("Updating student", zap.String("id", req.ID.String()))

	student, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("student not found: %w", err)
	}

	student.Apply(toStudentPatch(req))

	if err := s.repo.Update(ctx, student); err != nil {
		s.log.Error("Failed to update student", zap.Error(err))
		return nil, fmt.Errorf("failed to update student: %w", err)
	}

	return student, nil
}

func (s *studentService) GetStudent(ctx context.Context, id uuid.UUID) (*domain.Student, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *studentService) ListStudents(ctx context.Context, limit, offset int) ([]*domain.Student, int, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.List(ctx, limit, offset)
}

func (s *studentService) ListStudentsBySchool(ctx context.Context, schoolID uuid.UUID, limit, offset int) ([]*domain.Student, int, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.ListBySchool(ctx, schoolID, limit, offset)
}

func (s *studentService) DeactivateStudent(ctx context.Context, id uuid.UUID) error {
	s.log.Info("Deactivating student", zap.String("id", id.String()))

	student, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("student not found: %w", err)
	}

	inactive := false
	student.Apply(domain.StudentPatch{IsActive: &inactive})

	if err := s.repo.Update(ctx, student); err != nil {
		return fmt.Errorf("failed to deactivate student: %w", err)
	}

	return nil
}

// toStudentPatch maps an UpdateStudentRequest to a domain StudentPatch.
func toStudentPatch(req services.UpdateStudentRequest) domain.StudentPatch {
	return domain.StudentPatch{
		FirstName:      nonEmptyStr(req.FirstName),
		LastName:       nonEmptyStr(req.LastName),
		Grade:          nonEmptyStr(req.Grade),
		PickupLocation: req.PickupLocation,
		PickupAddress:  nonEmptyStr(req.PickupAddress),
		PhotoURL:       nonEmptyStr(req.PhotoURL),
	}
}

// StudentModule provides the student service to the fx graph.
var StudentModule = fx.Provide(NewStudentService)
