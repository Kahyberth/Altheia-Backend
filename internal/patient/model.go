package patient

import (
	"gorm.io/gorm"
	"time"
)

type Patient struct {
	ID             string         `gorm:"primaryKey" json:"id"`
	UserID         string         `json:"user_id"`
	DocumentNumber string         `json:"document_number"`
	DateOfBirth    string         `json:"date_of_birth"`
	Gender         string         `json:"gender"`
	Address        string         `json:"address"`
	Phone          string         `json:"phone"`
	Eps            string         `json:"eps"`
	BloodType      string         `json:"blood_type"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
