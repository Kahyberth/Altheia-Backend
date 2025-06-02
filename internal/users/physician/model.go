package physician

import (
	"errors"
	"time"
)

// Physician represents a physician in the system
type Physician struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Specialty     string    `json:"specialty"`
	LicenseNumber string    `json:"license_number"`
	Status        bool      `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// validatePhysician validates a physician's data
func validatePhysician(p *Physician) error {
	if p.UserID == "" {
		return errors.New("user ID is required")
	}

	if p.Specialty == "" {
		return errors.New("specialty is required")
	}

	if p.LicenseNumber == "" {
		return errors.New("license number is required")
	}

	// Validate license number format (alphanumeric)
	for _, c := range p.LicenseNumber {
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')) {
			return errors.New("license number must be alphanumeric")
		}
	}

	return nil
}
