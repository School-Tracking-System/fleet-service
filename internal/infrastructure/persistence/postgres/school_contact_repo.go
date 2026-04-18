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
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type schoolContactRepo struct {
	db  *sql.DB
	log *zap.Logger
}

// NewSchoolContactRepository creates a new PostgreSQL-backed SchoolContactRepository.
func NewSchoolContactRepository(db *sql.DB, log *zap.Logger) repositories.SchoolContactRepository {
	return &schoolContactRepo{db: db, log: log}
}

func (r *schoolContactRepo) Create(ctx context.Context, c *domain.SchoolContact) error {
	query := `
		INSERT INTO school_contacts (id, school_id, user_id, position, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		c.ID, c.SchoolID, c.UserID, nullableString(c.Position), c.IsActive, c.CreatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "unique constraint") {
			return domain.ErrDuplicateContact
		}
		return fmt.Errorf("failed to insert school contact: %w", err)
	}
	return nil
}

func (r *schoolContactRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.SchoolContact, error) {
	query := `
		SELECT id, school_id, user_id, position, is_active, created_at
		FROM school_contacts WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return scanContact(row)
}

func (r *schoolContactRepo) GetBySchoolID(ctx context.Context, schoolID uuid.UUID) ([]*domain.SchoolContact, error) {
	query := `
		SELECT id, school_id, user_id, position, is_active, created_at
		FROM school_contacts WHERE school_id = $1 AND is_active = TRUE
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, fmt.Errorf("failed to query school contacts: %w", err)
	}
	defer rows.Close()

	var contacts []*domain.SchoolContact
	for rows.Next() {
		c := &domain.SchoolContact{}
		var position sql.NullString
		if err := rows.Scan(&c.ID, &c.SchoolID, &c.UserID, &position, &c.IsActive, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan contact row: %w", err)
		}
		c.Position = position.String
		contacts = append(contacts, c)
	}
	return contacts, nil
}

func (r *schoolContactRepo) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE school_contacts SET is_active = FALSE WHERE id = $1`, id,
	)
	if err != nil {
		return fmt.Errorf("failed to deactivate school contact: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return domain.ErrContactNotFound
	}
	return nil
}

func scanContact(row *sql.Row) (*domain.SchoolContact, error) {
	c := &domain.SchoolContact{}
	var position sql.NullString
	err := row.Scan(&c.ID, &c.SchoolID, &c.UserID, &position, &c.IsActive, &c.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrContactNotFound
		}
		return nil, fmt.Errorf("failed to scan school contact: %w", err)
	}
	c.Position = position.String
	return c, nil
}

// SchoolContactRepoModule provides the repository to the fx graph.
var SchoolContactRepoModule = fx.Provide(NewSchoolContactRepository)
