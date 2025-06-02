package appointments

import (
	"Altheia-Backend/internal/users"
	"errors"
	"time"

	"gorm.io/gorm"
)

type MedicalAppointment struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	PatientId   string    `json:"patient_id"`
	PhysicianId string    `json:"physician_id"`
	DateTime    time.Time `json:"date_time"`
	Status      string    `json:"status"`
	Reason      string    `json:"reason"`

	Patient   users.Patient   `gorm:"foreignKey:PatientId"`
	Physician users.Physician `gorm:"foreignKey:PhysicianId"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// validateAppointment validates a medical appointment's data
func validateAppointment(a *MedicalAppointment) error {
	if a.PatientId == "" {
		return errors.New("patient ID is required")
	}

	if a.PhysicianId == "" {
		return errors.New("physician ID is required")
	}

	if a.DateTime.IsZero() {
		return errors.New("date and time are required")
	}

	if a.DateTime.Before(time.Now()) {
		return errors.New("appointment date and time must be in the future")
	}

	if a.Status == "" {
		return errors.New("status is required")
	}

	// Validate status
	validStatuses := map[string]bool{
		string(AppointmentStatusPending):   true,
		string(AppointmentStatusConfirmed): true,
		string(AppointmentStatusCancelled): true,
		string(AppointmentStatusCompleted): true,
		string(AppointmentStatusNoShow):    true,
	}

	if !validStatuses[a.Status] {
		return errors.New("invalid appointment status")
	}

	return nil
}
