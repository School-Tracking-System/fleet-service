package domain

import "errors"

// Common domain errors that can be used across the application.
// These errors should not contain infrastructure-specific details (like SQL or gRPC codes).
var (
	// Vehicle errors
	ErrVehicleNotFound  = errors.New("vehicle not found")
	ErrDuplicateVehicle = errors.New("a vehicle with this plate already exists")
	ErrInvalidVehicle   = errors.New("invalid vehicle data")

	// Driver errors
	ErrDriverNotFound      = errors.New("driver not found")
	ErrDuplicateDriver     = errors.New("a driver with this license or cedula already exists")
	ErrInvalidDriver       = errors.New("invalid driver data")
	ErrDriverAlreadyLinked = errors.New("this user is already registered as a driver")

	// School errors
	ErrSchoolNotFound   = errors.New("school not found")
	ErrDuplicateSchool  = errors.New("school already exists")
	ErrInvalidSchool    = errors.New("invalid school data")
	ErrContactNotFound  = errors.New("school contact not found")
	ErrDuplicateContact = errors.New("contact already linked to school")
	ErrInvalidContact   = errors.New("invalid contact data")

	// Student errors
	ErrStudentNotFound  = errors.New("student not found")
	ErrDuplicateStudent = errors.New("a student with this identification already exists")
	ErrInvalidStudent   = errors.New("invalid student data")

	ErrGuardianNotFound  = errors.New("guardian not found")
	ErrDuplicateGuardian = errors.New("this user is already a guardian for this student")
	ErrInvalidGuardian   = errors.New("invalid guardian data")

	// Route errors
	ErrRouteNotFound = errors.New("route not found")
	ErrInvalidRoute  = errors.New("invalid route data")

	// Stop errors
	ErrStopNotFound = errors.New("stop not found")
	ErrInvalidStop  = errors.New("invalid stop data")
)
