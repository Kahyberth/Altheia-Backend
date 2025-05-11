package auth

import (
	"Altheia-Backend/internal/users/patient"
	"Altheia-Backend/internal/users/physician"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	Name      string         `json:"name"`
	Email     string         `gorm:"unique" json:"email"`
	Password  string         `json:"password"`
	Rol       string         `json:"rol"`
	Phone     string         `json:"phone"`
	Status    bool           `json:"status"`
	Gender    string         `json:"gender"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	LastLogin time.Time      `json:"lastLogin"`

	Patient   *patient.Patient     `gorm:"foreignKey:UserID;references:ID" json:"patient,omitempty"`
	Physician *physician.Physician `gorm:"foreignKey:UserID;references:ID" json:"physician,omitempty"`
}
