package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type RouteDirection string

const (
	RouteDirectionToSchool   RouteDirection = "to_school"
	RouteDirectionFromSchool RouteDirection = "from_school"
)

// Route represents a scheduled transportation path.
type Route struct {
	ID           uuid.UUID
	Name         string
	Description  string
	VehicleID    *uuid.UUID // Optional: can be nil if no vehicle assigned
	DriverID     *uuid.UUID // Optional: can be nil if no driver assigned
	SchoolID     uuid.UUID
	Direction    RouteDirection
	ScheduleTime string // HH:MM format
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Stops        []*RouteStop // Loaded optionally
}

// RouteStop represents a specific pickup/dropoff point in a route.
type RouteStop struct {
	ID        uuid.UUID
	RouteID   uuid.UUID
	StudentID uuid.UUID
	Order     int
	Location  Location
	Address   string
	EstTime   string // HH:MM format
	CreatedAt time.Time
}

// NewRouteParams holds data for creating a Route.
type NewRouteParams struct {
	Name         string
	Description  string
	VehicleID    *uuid.UUID
	DriverID     *uuid.UUID
	SchoolID     uuid.UUID
	Direction    RouteDirection
	ScheduleTime string
}

// NewRoute creates a valid Route instance.
func NewRoute(p NewRouteParams) (*Route, error) {
	if p.Name == "" {
		return nil, errors.New("route name is required")
	}
	if p.SchoolID == uuid.Nil {
		return nil, errors.New("school_id is required")
	}
	if p.ScheduleTime == "" {
		return nil, errors.New("schedule time is required")
	}

	now := time.Now().UTC()
	return &Route{
		ID:           uuid.New(),
		Name:         p.Name,
		Description:  p.Description,
		VehicleID:    p.VehicleID,
		DriverID:     p.DriverID,
		SchoolID:     p.SchoolID,
		Direction:    p.Direction,
		ScheduleTime: p.ScheduleTime,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// RoutePatch for partial updates.
type RoutePatch struct {
	Name         *string
	Description  *string
	VehicleID    **uuid.UUID // double pointer to allow setting to nil
	DriverID     **uuid.UUID // double pointer to allow setting to nil
	Direction    *RouteDirection
	ScheduleTime *string
	IsActive     *bool
}

// Apply updates the route.
func (r *Route) Apply(patch RoutePatch) {
	if patch.Name != nil {
		r.Name = *patch.Name
	}
	if patch.Description != nil {
		r.Description = *patch.Description
	}
	if patch.VehicleID != nil {
		r.VehicleID = *patch.VehicleID
	}
	if patch.DriverID != nil {
		r.DriverID = *patch.DriverID
	}
	if patch.Direction != nil {
		r.Direction = *patch.Direction
	}
	if patch.ScheduleTime != nil {
		r.ScheduleTime = *patch.ScheduleTime
	}
	if patch.IsActive != nil {
		r.IsActive = *patch.IsActive
	}
	r.UpdatedAt = time.Now().UTC()
}

// NewStopParams holds data for creating a RouteStop.
type NewStopParams struct {
	RouteID   uuid.UUID
	StudentID uuid.UUID
	Order     int
	Location  Location
	Address   string
	EstTime   string
}

// NewRouteStop creates a valid RouteStop instance.
func NewRouteStop(p NewStopParams) (*RouteStop, error) {
	if p.RouteID == uuid.Nil {
		return nil, errors.New("route_id is required")
	}
	if p.StudentID == uuid.Nil {
		return nil, errors.New("student_id is required")
	}
	return &RouteStop{
		ID:        uuid.New(),
		RouteID:   p.RouteID,
		StudentID: p.StudentID,
		Order:     p.Order,
		Location:  p.Location,
		Address:   p.Address,
		EstTime:   p.EstTime,
		CreatedAt: time.Now().UTC(),
	}, nil
}
