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

type schoolRepo struct {
	db  *sql.DB
	log *zap.Logger
}

// NewSchoolRepository creates a new PostgreSQL-backed SchoolRepository.
func NewSchoolRepository(db *sql.DB, log *zap.Logger) repositories.SchoolRepository {
	return &schoolRepo{db: db, log: log}
}

func (r *schoolRepo) Create(ctx context.Context, s *domain.School) error {
	query := `
		INSERT INTO schools (id, name, address, location, phone, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		s.ID, s.Name, s.Address, locationToWKT(s.Location),
		nullableString(s.Phone), nullableString(s.Email), s.CreatedAt, s.UpdatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "unique constraint") {
			return domain.ErrDuplicateSchool
		}
		return fmt.Errorf("failed to insert school: %w", err)
	}
	return nil
}

func (r *schoolRepo) Update(ctx context.Context, s *domain.School) error {
	query := `
		UPDATE schools
		SET name = $1, address = $2, location = $3, phone = $4, email = $5, updated_at = NOW()
		WHERE id = $6
	`
	result, err := r.db.ExecContext(ctx, query,
		s.Name, s.Address, locationToWKT(s.Location),
		nullableString(s.Phone), nullableString(s.Email), s.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update school: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return domain.ErrSchoolNotFound
	}
	return nil
}

func (r *schoolRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.School, error) {
	// ST_AsText returns WKT: "POINT(lon lat)"
	query := `
		SELECT id, name, address, ST_AsText(location), phone, email, created_at, updated_at
		FROM schools WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return scanSchool(row)
}

func (r *schoolRepo) List(ctx context.Context, limit, offset int) ([]*domain.School, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM schools`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count schools: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, address, ST_AsText(location), phone, email, created_at, updated_at
		FROM schools ORDER BY name ASC LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query schools: %w", err)
	}
	defer rows.Close()

	var schools []*domain.School
	for rows.Next() {
		s := &domain.School{}
		var locationWKT, phone, email sql.NullString
		if err := rows.Scan(&s.ID, &s.Name, &s.Address, &locationWKT, &phone, &email, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan school row: %w", err)
		}
		s.Phone = phone.String
		s.Email = email.String
		s.Location = parseWKT(locationWKT)
		schools = append(schools, s)
	}
	return schools, total, nil
}

// scanSchool scans a single *sql.Row into a School.
func scanSchool(row *sql.Row) (*domain.School, error) {
	s := &domain.School{}
	var locationWKT, phone, email sql.NullString
	err := row.Scan(&s.ID, &s.Name, &s.Address, &locationWKT, &phone, &email, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrSchoolNotFound
		}
		return nil, fmt.Errorf("failed to scan school: %w", err)
	}
	s.Phone = phone.String
	s.Email = email.String
	s.Location = parseWKT(locationWKT)
	return s, nil
}

// locationToWKT converts a *Location to a WKT string for PostGIS.
// Returns nil if the location is not set (stored as NULL in the DB).
func locationToWKT(loc *domain.Location) *string {
	if loc == nil {
		return nil
	}
	wkt := fmt.Sprintf("SRID=4326;POINT(%f %f)", loc.Longitude, loc.Latitude)
	return &wkt
}

// parseWKT parses a PostGIS WKT string "POINT(lon lat)" into a *Location.
// Returns nil if the value is NULL or malformed.
func parseWKT(wkt sql.NullString) *domain.Location {
	if !wkt.Valid || wkt.String == "" {
		return nil
	}
	var lon, lat float64
	_, err := fmt.Sscanf(wkt.String, "POINT(%f %f)", &lon, &lat)
	if err != nil {
		return nil
	}
	return &domain.Location{Longitude: lon, Latitude: lat}
}
