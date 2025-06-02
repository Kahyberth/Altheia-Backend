package receptionist

import (
	"errors"
	"time"
)

// Receptionist represents a receptionist in the system
type Receptionist struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ClinicID  string    `json:"clinic_id"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// validateReceptionist validates a receptionist's data
func validateReceptionist(r *Receptionist) error {
	if r.UserID == "" {
		return errors.New("user ID is required")
	}

	if r.ClinicID == "" {
		return errors.New("clinic ID is required")
	}

	return nil
}
