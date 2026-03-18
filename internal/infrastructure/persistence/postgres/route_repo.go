package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"


	"github.com/fercho/school-tracking/services/fleet/internal/core/domain"
	"github.com/fercho/school-tracking/services/fleet/internal/core/ports/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type routeRepo struct {
	db  *sql.DB
	log *zap.Logger
}

func NewRouteRepository(db *sql.DB, log *zap.Logger) repositories.RouteRepository {
	return &routeRepo{db: db, log: log}
}

func (r *routeRepo) Create(ctx context.Context, rt *domain.Route) error {
	query := `
		INSERT INTO routes (id, name, description, vehicle_id, driver_id, school_id, direction, schedule_time, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, query,
		rt.ID, rt.Name, nullableString(rt.Description), rt.VehicleID, rt.DriverID,
		rt.SchoolID, string(rt.Direction), rt.ScheduleTime, rt.IsActive, rt.CreatedAt, rt.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert route: %w", err)
	}
	return nil
}

func (r *routeRepo) Update(ctx context.Context, rt *domain.Route) error {
	query := `
		UPDATE routes
		SET name = $1, description = $2, vehicle_id = $3, driver_id = $4, direction = $5, schedule_time = $6, is_active = $7, updated_at = NOW()
		WHERE id = $8
	`
	result, err := r.db.ExecContext(ctx, query,
		rt.Name, nullableString(rt.Description), rt.VehicleID, rt.DriverID,
		string(rt.Direction), rt.ScheduleTime, rt.IsActive, rt.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update route: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrRouteNotFound
	}
	return nil
}

func (r *routeRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Route, error) {
	rt := &domain.Route{}
	var desc sql.NullString
	query := `SELECT id, name, description, vehicle_id, driver_id, school_id, direction, schedule_time, is_active, created_at, updated_at FROM routes WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rt.ID, &rt.Name, &desc, &rt.VehicleID, &rt.DriverID,
		&rt.SchoolID, &rt.Direction, &rt.ScheduleTime, &rt.IsActive, &rt.CreatedAt, &rt.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRouteNotFound
		}
		return nil, fmt.Errorf("failed to get route: %w", err)
	}
	rt.Description = desc.String
	return rt, nil
}

func (r *routeRepo) List(ctx context.Context, limit, offset int) ([]*domain.Route, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM routes`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, description, vehicle_id, driver_id, school_id, direction, schedule_time, is_active, created_at, updated_at FROM routes ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	return scanRouteRows(rows, total)
}

func (r *routeRepo) ListBySchool(ctx context.Context, schoolID uuid.UUID, limit, offset int) ([]*domain.Route, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM routes WHERE school_id = $1`, schoolID).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, description, vehicle_id, driver_id, school_id, direction, schedule_time, is_active, created_at, updated_at FROM routes WHERE school_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, schoolID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	return scanRouteRows(rows, total)
}

func (r *routeRepo) CreateStop(ctx context.Context, s *domain.RouteStop) error {
	query := `INSERT INTO route_stops (id, route_id, student_id, stop_order, location, address, est_time, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, s.ID, s.RouteID, s.StudentID, s.Order, locationToWKT(&s.Location), nullableString(s.Address), s.EstTime, s.CreatedAt)
	return err
}

func (r *routeRepo) DeleteStop(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM route_stops WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domain.ErrStopNotFound
	}
	return nil
}

func (r *routeRepo) UpdateStopOrder(ctx context.Context, id uuid.UUID, order int) error {
	res, err := r.db.ExecContext(ctx, `UPDATE route_stops SET stop_order = $1 WHERE id = $2`, order, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domain.ErrStopNotFound
	}
	return nil
}

func (r *routeRepo) GetStopsByRouteID(ctx context.Context, routeID uuid.UUID) ([]*domain.RouteStop, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, route_id, student_id, stop_order, ST_AsText(location), address, est_time, created_at FROM route_stops WHERE route_id = $1 ORDER BY stop_order ASC`, routeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var stops []*domain.RouteStop
	for rows.Next() {
		s := &domain.RouteStop{}
		var locationWKT, address sql.NullString
		if err := rows.Scan(&s.ID, &s.RouteID, &s.StudentID, &s.Order, &locationWKT, &address, &s.EstTime, &s.CreatedAt); err != nil {
			return nil, err
		}
		s.Address = address.String
		loc := parseWKT(locationWKT)
		if loc != nil {
			s.Location = *loc
		}
		stops = append(stops, s)
	}
	return stops, nil
}

func (r *routeRepo) GetStopByID(ctx context.Context, id uuid.UUID) (*domain.RouteStop, error) {
	s := &domain.RouteStop{}
	var locationWKT, address sql.NullString
	err := r.db.QueryRowContext(ctx, `SELECT id, route_id, student_id, stop_order, ST_AsText(location), address, est_time, created_at FROM route_stops WHERE id = $1`, id).Scan(
		&s.ID, &s.RouteID, &s.StudentID, &s.Order, &locationWKT, &address, &s.EstTime, &s.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrStopNotFound
		}
		return nil, err
	}
	s.Address = address.String
	loc := parseWKT(locationWKT)
	if loc != nil {
		s.Location = *loc
	}
	return s, nil
}

func scanRouteRows(rows *sql.Rows, total int) ([]*domain.Route, int, error) {
	var routes []*domain.Route
	for rows.Next() {
		rt := &domain.Route{}
		var desc sql.NullString
		err := rows.Scan(&rt.ID, &rt.Name, &desc, &rt.VehicleID, &rt.DriverID, &rt.SchoolID, &rt.Direction, &rt.ScheduleTime, &rt.IsActive, &rt.CreatedAt, &rt.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		rt.Description = desc.String
		routes = append(routes, rt)
	}
	return routes, total, nil
}
