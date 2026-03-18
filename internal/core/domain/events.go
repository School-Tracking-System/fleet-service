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
