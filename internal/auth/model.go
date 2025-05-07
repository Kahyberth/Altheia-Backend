package auth

import (
	"Altheia-Backend/internal/patient"
	"Altheia-Backend/internal/physician"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	Name      string         `json:"name"`
	Email     string         `gorm:"unique" json:"email"`
	Password  string         `json:"password"`
	Rol       string         `json:"rol"`
	Status    bool           `json:"status"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	LastLogin time.Time `json:"lastLogin"`

	Patient   *patient.Patient     `gorm:"foreignKey:UserId" json:"patient,omitempty"`
	Physician *physician.Physician `gorm:"foreignKey:UserId" json:"physician,omitempty"`
}
