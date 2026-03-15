package domain

import "errors"

// Common domain errors that can be used across the application.
// These errors should not contain infrastructure-specific details (like SQL or gRPC codes).
var (
	ErrVehicleNotFound  = errors.New("vehicle not found")
	ErrDuplicateVehicle = errors.New("a vehicle with this plate already exists")
	ErrInvalidVehicle   = errors.New("invalid vehicle data")
)
