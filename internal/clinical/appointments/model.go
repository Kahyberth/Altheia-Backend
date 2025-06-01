package appointments

import (
	"Altheia-Backend/internal/users"
	"gorm.io/gorm"
	"time"
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
