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

type driverRepo struct {
	db  *sql.DB
	log *zap.Logger
}

// NewDriverRepository creates a new PostgreSQL-backed DriverRepository.
func NewDriverRepository(db *sql.DB, log *zap.Logger) repositories.DriverRepository {
	return &driverRepo{db: db, log: log}
}

func (r *driverRepo) Create(ctx context.Context, d *domain.Driver) error {
	query := `
		INSERT INTO drivers (id, user_id, license_number, license_type, license_expiry, cedula_id, emergency_phone, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		d.ID, d.UserID, d.LicenseNumber, d.LicenseType, d.LicenseExpiry,
		d.CedulaID, nullableString(d.EmergencyPhone), d.Status, d.CreatedAt, d.UpdatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "unique constraint") {
			return domain.ErrDuplicateDriver
		}
		return fmt.Errorf("failed to insert driver: %w", err)
	}
	return nil
}

func (r *driverRepo) Update(ctx context.Context, d *domain.Driver) error {
	query := `
		UPDATE drivers
		SET license_type = $1, license_expiry = $2, emergency_phone = $3, status = $4, updated_at = NOW()
		WHERE id = $5
	`
	result, err := r.db.ExecContext(ctx, query,
		d.LicenseType, d.LicenseExpiry, nullableString(d.EmergencyPhone), d.Status, d.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update driver: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return domain.ErrDriverNotFound
	}
	return nil
}

func (r *driverRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Driver, error) {
	return r.scanOne(ctx, `
		SELECT id, user_id, license_number, license_type, license_expiry, cedula_id, emergency_phone, status, created_at, updated_at
		FROM drivers WHERE id = $1
	`, id)
}

func (r *driverRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Driver, error) {
	return r.scanOne(ctx, `
		SELECT id, user_id, license_number, license_type, license_expiry, cedula_id, emergency_phone, status, created_at, updated_at
		FROM drivers WHERE user_id = $1
	`, userID)
}

func (r *driverRepo) List(ctx context.Context, limit, offset int) ([]*domain.Driver, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM drivers`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count drivers: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, license_number, license_type, license_expiry, cedula_id, emergency_phone, status, created_at, updated_at
		FROM drivers ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query drivers: %w", err)
	}
	defer rows.Close()

	var drivers []*domain.Driver
	for rows.Next() {
		d, err := scanDriver(rows)
		if err != nil {
			return nil, 0, err
		}
		drivers = append(drivers, d)
	}
	return drivers, total, nil
}

// scanOne executes a query and scans a single Driver row.
func (r *driverRepo) scanOne(ctx context.Context, query string, arg interface{}) (*domain.Driver, error) {
	row := r.db.QueryRowContext(ctx, query, arg)
	d := &domain.Driver{}
	var emergencyPhone sql.NullString

	err := row.Scan(
		&d.ID, &d.UserID, &d.LicenseNumber, &d.LicenseType, &d.LicenseExpiry,
		&d.CedulaID, &emergencyPhone, &d.Status, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrDriverNotFound
		}
		return nil, fmt.Errorf("failed to scan driver: %w", err)
	}
	d.EmergencyPhone = emergencyPhone.String
	return d, nil
}

// scanDriver scans a *sql.Rows into a Driver.
func scanDriver(rows *sql.Rows) (*domain.Driver, error) {
	d := &domain.Driver{}
	var emergencyPhone sql.NullString
	err := rows.Scan(
		&d.ID, &d.UserID, &d.LicenseNumber, &d.LicenseType, &d.LicenseExpiry,
		&d.CedulaID, &emergencyPhone, &d.Status, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan driver row: %w", err)
	}
	d.EmergencyPhone = emergencyPhone.String
	return d, nil
}
