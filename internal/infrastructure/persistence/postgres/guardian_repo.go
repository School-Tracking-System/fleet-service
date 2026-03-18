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

type guardianRepo struct {
	db  *sql.DB
	log *zap.Logger
}

// NewGuardianRepository creates a new PostgreSQL-backed GuardianRepository.
func NewGuardianRepository(db *sql.DB, log *zap.Logger) repositories.GuardianRepository {
	return &guardianRepo{db: db, log: log}
}

func (r *guardianRepo) Create(ctx context.Context, g *domain.Guardian) error {
	query := `
		INSERT INTO guardians (id, user_id, student_id, relation, is_primary, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		g.ID, g.UserID, g.StudentID, string(g.Relation), g.IsPrimary, g.CreatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "23505") { // unique (user_id, student_id)
			return domain.ErrDuplicateGuardian
		}
		if strings.Contains(err.Error(), "23503") { // FK violation on student_id
			return fmt.Errorf("student not found: %w", domain.ErrStudentNotFound)
		}
		return fmt.Errorf("failed to link guardian: %w", err)
	}
	return nil
}

func (r *guardianRepo) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM guardians WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete guardian: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrGuardianNotFound
	}
	return nil
}

func (r *guardianRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Guardian, error) {
	g := &domain.Guardian{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, student_id, relation, is_primary, created_at FROM guardians WHERE id = $1`,
		id,
	).Scan(&g.ID, &g.UserID, &g.StudentID, &g.Relation, &g.IsPrimary, &g.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrGuardianNotFound
		}
		return nil, fmt.Errorf("failed to get guardian: %w", err)
	}
	return g, nil
}

func (r *guardianRepo) GetByStudentID(ctx context.Context, studentID uuid.UUID) ([]*domain.Guardian, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, student_id, relation, is_primary, created_at FROM guardians WHERE student_id = $1 ORDER BY is_primary DESC`,
		studentID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get guardians: %w", err)
	}
	defer rows.Close()

	var guardians []*domain.Guardian
	for rows.Next() {
		g := &domain.Guardian{}
		if err := rows.Scan(&g.ID, &g.UserID, &g.StudentID, &g.Relation, &g.IsPrimary, &g.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan guardian: %w", err)
		}
		guardians = append(guardians, g)
	}
	return guardians, nil
}

func (r *guardianRepo) GetStudentsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Student, error) {
	// Joins guardians → students to get all students for a given guardian user
	query := `
		SELECT s.id, s.first_name, s.last_name, s.grade, s.school_id,
		       ST_AsText(s.pickup_location), s.pickup_address, s.photo_url, s.is_active,
		       s.created_at, s.updated_at
		FROM students s
		INNER JOIN guardians g ON g.student_id = s.id
		WHERE g.user_id = $1
		ORDER BY s.last_name, s.first_name
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get students by guardian: %w", err)
	}
	defer rows.Close()

	students, _, err := scanStudentRows(rows, 0)
	return students, err
}
