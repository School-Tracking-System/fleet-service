package domain

import "time"

// VehicleCreatedEvent is published after a vehicle is successfully persisted.
type VehicleCreatedEvent struct {
	ID        string    `json:"id"`
	Plate     string    `json:"plate"`
	Brand     string    `json:"brand"`
	Model     string    `json:"model"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// VehicleUpdatedEvent is published after a vehicle is successfully updated.
type VehicleUpdatedEvent struct {
	ID        string    `json:"id"`
	Plate     string    `json:"plate"`
	Brand     string    `json:"brand"`
	Model     string    `json:"model"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

// StudentAssignedEvent is published after a student is registered and assigned to a school.
type StudentAssignedEvent struct {
	StudentID string    `json:"student_id"`
	SchoolID  string    `json:"school_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
}
