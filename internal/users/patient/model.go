package patient

import (
	"errors"
	"time"
)

type Patient struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	DateOfBirth string    `json:"date_of_birth"`
	Address     string    `json:"address"`
	Eps         string    `json:"eps"`
	BloodType   string    `json:"blood_type"`
	Status      bool      `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func validatePatient(p *Patient) error {
	if p.UserID == "" {
		return errors.New("user ID is required")
	}

	if p.DateOfBirth == "" {
		return errors.New("date of birth is required")
	}

	_, err := time.Parse("2006-01-02", p.DateOfBirth)
	if err != nil {
		return errors.New("invalid date of birth format, expected YYYY-MM-DD")
	}

	if p.BloodType == "" {
		return errors.New("blood type is required")
	}

	validBloodTypes := map[string]bool{
		"A+": true, "A-": true,
		"B+": true, "B-": true,
		"AB+": true, "AB-": true,
		"O+": true, "O-": true,
	}

	if !validBloodTypes[p.BloodType] {
		return errors.New("invalid blood type")
	}

	return nil
}
