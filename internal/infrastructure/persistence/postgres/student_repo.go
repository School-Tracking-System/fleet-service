package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type studentRepo struct {
	db  *sql.DB
	log *zap.Logger
}

// NewStudentRepository creates a new PostgreSQL-backed StudentRepository.
func NewStudentRepository(db *sql.DB, log *zap.Logger) repositories.StudentRepository {
	return &studentRepo{db: db, log: log}
}

const studentSelectCols = `
	id, first_name, last_name, grade, school_id,
	ST_AsText(pickup_location), pickup_address, photo_url, is_active,
	cedula_id, created_at, updated_at`

func (r *studentRepo) Create(ctx context.Context, s *domain.Student) error {
	query := `
		INSERT INTO students (id, first_name, last_name, grade, school_id, pickup_location, pickup_address, photo_url, is_active, cedula_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.ExecContext(ctx, query,
		s.ID, s.FirstName, s.LastName,
		nullableString(s.Grade), s.SchoolID, locationToWKT(s.PickupLocation),
		nullableString(s.PickupAddress), nullableString(s.PhotoURL),
		s.IsActive, s.CedulaID, s.CreatedAt, s.UpdatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "23503") { // FK violation
			return fmt.Errorf("school not found for student: %w", domain.ErrSchoolNotFound)
		}
		if strings.Contains(err.Error(), "23505") { // Unique violation
			return domain.ErrDuplicateStudent
		}
		return fmt.Errorf("failed to insert student: %w", err)
	}
	return nil
}

func (r *studentRepo) Update(ctx context.Context, s *domain.Student) error {
	query := `
		UPDATE students
		SET first_name = $1, last_name = $2, grade = $3,
		    pickup_location = $4, pickup_address = $5, photo_url = $6,
		    is_active = $7, cedula_id = $8, updated_at = NOW()
		WHERE id = $9
	`
	result, err := r.db.ExecContext(ctx, query,
		s.FirstName, s.LastName, nullableString(s.Grade),
		locationToWKT(s.PickupLocation), nullableString(s.PickupAddress),
		nullableString(s.PhotoURL), s.IsActive, s.CedulaID, s.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update student: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrStudentNotFound
	}
	return nil
}

func (r *studentRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Student, error) {
	query := `SELECT ` + studentSelectCols + ` FROM students WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	return scanStudent(row)
}

func (r *studentRepo) GetByCedulaID(ctx context.Context, cedulaID string) (*domain.Student, error) {
	query := `SELECT ` + studentSelectCols + ` FROM students WHERE cedula_id = $1`
	row := r.db.QueryRowContext(ctx, query, cedulaID)
	return scanStudent(row)
}

func (r *studentRepo) List(ctx context.Context, limit, offset int) ([]*domain.Student, int, error) {
	return r.listStudents(ctx, `SELECT COUNT(*) FROM students`, `SELECT ` + studentSelectCols + ` FROM students ORDER BY last_name, first_name LIMIT $1 OFFSET $2`, limit, offset)
}

func (r *studentRepo) ListBySchool(ctx context.Context, schoolID uuid.UUID, limit, offset int) ([]*domain.Student, int, error) {
	countQuery := `SELECT COUNT(*) FROM students WHERE school_id = $1`
	listQuery := `SELECT ` + studentSelectCols + ` FROM students WHERE school_id = $1 ORDER BY last_name, first_name LIMIT $2 OFFSET $3`

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, schoolID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count students by school: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, listQuery, schoolID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list students by school: %w", err)
	}
	defer rows.Close()

	return scanStudentRows(rows, total)
}

// listStudents is a shared helper for list queries without extra parameters.
func (r *studentRepo) listStudents(ctx context.Context, countSQL, listSQL string, limit, offset int) ([]*domain.Student, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, countSQL).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count students: %w", err)
	}
	rows, err := r.db.QueryContext(ctx, listSQL, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list students: %w", err)
	}
	defer rows.Close()
	return scanStudentRows(rows, total)
}

func scanStudent(row *sql.Row) (*domain.Student, error) {
	s := &domain.Student{}
	var grade, pickupAddr, photoURL, locationWKT sql.NullString
	err := row.Scan(
		&s.ID, &s.FirstName, &s.LastName, &grade, &s.SchoolID,
		&locationWKT, &pickupAddr, &photoURL, &s.IsActive,
		&s.CedulaID, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrStudentNotFound
		}
		return nil, fmt.Errorf("failed to scan student: %w", err)
	}
	s.Grade = grade.String
	s.PickupAddress = pickupAddr.String
	s.PhotoURL = photoURL.String
	s.PickupLocation = parseWKT(locationWKT)
	return s, nil
}

func scanStudentRows(rows *sql.Rows, total int) ([]*domain.Student, int, error) {
	var students []*domain.Student
	for rows.Next() {
		s := &domain.Student{}
		var grade, pickupAddr, photoURL, locationWKT sql.NullString
		if err := rows.Scan(
			&s.ID, &s.FirstName, &s.LastName, &grade, &s.SchoolID,
			&locationWKT, &pickupAddr, &photoURL, &s.IsActive,
			&s.CedulaID, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan student row: %w", err)
		}
		s.Grade = grade.String
		s.PickupAddress = pickupAddr.String
		s.PhotoURL = photoURL.String
		s.PickupLocation = parseWKT(locationWKT)
		students = append(students, s)
	}
	return students, total, nil
}
