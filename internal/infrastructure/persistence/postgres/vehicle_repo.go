package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/repositories"
	"github.com/fercho/school-tracking/services/fleet/pkg/env"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type vehicleRepo struct {
	db  *sql.DB
	log *zap.Logger
}

// NewDatabase opens a connection to PostgreSQL and verifies it is reachable.
// Schema creation and table migrations are handled externally by Flyway.
func NewDatabase(cfg *env.Config, log *zap.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("error opening db connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	log.Info("Successfully connected to PostgreSQL database")
	return db, nil
}

func NewVehicleRepository(db *sql.DB, log *zap.Logger) repositories.VehicleRepository {
	return &vehicleRepo{
		db:  db,
		log: log,
	}
}

func (r *vehicleRepo) Create(ctx context.Context, v *domain.Vehicle) error {
	query := `
		INSERT INTO vehicles (id, plate, brand, model, year, capacity, status, color, vehicle_type, chassis_num, insurance_exp, tech_review_exp, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	_, err := r.db.ExecContext(ctx, query,
		v.ID, v.Plate, v.Brand, v.Model, v.Year, v.Capacity, v.Status,
		nullableString(v.Color), nullableString(string(v.VehicleType)), nullableString(v.ChassisNum),
		v.InsuranceExp, v.TechReviewExp, v.CreatedAt, v.UpdatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "unique constraint") {
			return domain.ErrDuplicateVehicle
		}
		return fmt.Errorf("failed to insert vehicle: %w", err)
	}
	return nil
}

func (r *vehicleRepo) Update(ctx context.Context, v *domain.Vehicle) error {
	query := `
		UPDATE vehicles
		SET brand = $1, model = $2, year = $3, capacity = $4, status = $5,
		    color = $6, vehicle_type = $7, chassis_num = $8,
		    insurance_exp = $9, tech_review_exp = $10, updated_at = NOW()
		WHERE id = $11
	`
	result, err := r.db.ExecContext(ctx, query,
		v.Brand, v.Model, v.Year, v.Capacity, v.Status,
		nullableString(v.Color), nullableString(string(v.VehicleType)), nullableString(v.ChassisNum),
		v.InsuranceExp, v.TechReviewExp, v.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update vehicle: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return domain.ErrVehicleNotFound
	}
	return nil
}

func (r *vehicleRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Vehicle, error) {
	query := `
		SELECT id, plate, brand, model, year, capacity, status, color, vehicle_type, chassis_num, insurance_exp, tech_review_exp, created_at, updated_at
		FROM vehicles
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var v domain.Vehicle
	var color, vehicleType, chassisNum sql.NullString
	var insuranceExp, techReviewExp sql.NullTime
	if err := row.Scan(
		&v.ID, &v.Plate, &v.Brand, &v.Model, &v.Year, &v.Capacity, &v.Status,
		&color, &vehicleType, &chassisNum, &insuranceExp, &techReviewExp,
		&v.CreatedAt, &v.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrVehicleNotFound
		}
		return nil, fmt.Errorf("failed to scan vehicle row: %w", err)
	}

	v.Color = color.String
	v.VehicleType = domain.VehicleType(vehicleType.String)
	v.ChassisNum = chassisNum.String
	if insuranceExp.Valid {
		t := insuranceExp.Time
		v.InsuranceExp = &t
	}
	if techReviewExp.Valid {
		t := techReviewExp.Time
		v.TechReviewExp = &t
	}

	return &v, nil
}

func (r *vehicleRepo) List(ctx context.Context, limit, offset int) ([]*domain.Vehicle, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM vehicles`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count vehicles: %w", err)
	}

	query := `
		SELECT id, plate, brand, model, year, capacity, status, color, vehicle_type, chassis_num, insurance_exp, tech_review_exp, created_at, updated_at
		FROM vehicles
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query vehicles list: %w", err)
	}
	defer rows.Close()

	var vehicles []*domain.Vehicle
	for rows.Next() {
		var v domain.Vehicle
		var color, vehicleType, chassisNum sql.NullString
		var insuranceExp, techReviewExp sql.NullTime
		if err := rows.Scan(
			&v.ID, &v.Plate, &v.Brand, &v.Model, &v.Year, &v.Capacity, &v.Status,
			&color, &vehicleType, &chassisNum, &insuranceExp, &techReviewExp,
			&v.CreatedAt, &v.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan vehicle row: %w", err)
		}
		v.Color = color.String
		v.VehicleType = domain.VehicleType(vehicleType.String)
		v.ChassisNum = chassisNum.String
		if insuranceExp.Valid {
			t := insuranceExp.Time
			v.InsuranceExp = &t
		}
		if techReviewExp.Valid {
			t := techReviewExp.Time
			v.TechReviewExp = &t
		}
		vehicles = append(vehicles, &v)
	}

	return vehicles, total, nil
}

// nullableString converts an empty string to a nil value for optional DB columns.
func nullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Module provides the PostgreSQL infrastructure dependencies to the application graph.
var Module = fx.Options(
	fx.Provide(NewDatabase),
	fx.Provide(NewVehicleRepository),
	fx.Provide(NewRouteRepository),
	fx.Provide(NewDriverRepository),
	fx.Provide(NewStudentRepository),
	fx.Provide(NewGuardianRepository),
	fx.Provide(NewSchoolRepository),
	fx.Provide(NewSchoolContactRepository),
)
