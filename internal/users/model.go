package users

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID             string         `gorm:"primaryKey" json:"id"`
	Name           string         `json:"name"`
	Email          string         `gorm:"unique" json:"email"`
	Password       string         `json:"password"`
	Rol            string         `json:"rol"`
	Phone          string         `json:"phone"`
	DocumentNumber string         `json:"document_number"`
	Status         bool           `json:"status"`
	Gender         string         `json:"gender"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	LastLogin      time.Time      `json:"lastLogin"`

	Patient      Patient      `gorm:"foreignKey:UserID;references:ID" json:"patient,omitempty"`
	Physician    Physician    `gorm:"foreignKey:UserID;references:ID" json:"physician,omitempty"`
	Receptionist Receptionist `gorm:"foreignKey:UserID;references:ID" json:"receptionist,omitempty"`
	ClinicOwner  ClinicOwner  `gorm:"foreignKey:UserID;references:ID" json:"clinic_owner,omitempty"`
}

type Patient struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	UserID      string         `gorm:"not null;index" json:"user_id"`
	DateOfBirth string         `json:"date_of_birth"`
	Address     string         `json:"address"`
	Eps         string         `json:"eps"`
	BloodType   string         `json:"blood_type"`
	Status      bool           `json:"status"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Physician struct {
	ID                 string         `gorm:"primaryKey" json:"id"`
	UserID             string         `gorm:"not null;index" json:"user_id"`
	PhysicianSpecialty string         `json:"physician_specialty"`
	LicenseNumber      string         `json:"license_number"`
	Status             bool           `json:"status"`
	ClinicID           *string        `json:"clinic_id"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

type Receptionist struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	UserID    string         `gorm:"not null;index" json:"user_id"`
	ClinicID  *string        `json:"clinic_id"`
	Status    bool           `json:"status"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Pagination struct {
	Limit  int         `json:"limit"`
	Page   int         `json:"page"`
	Sort   string      `json:"sort"`
	Total  int64       `json:"total"`
	Result interface{} `json:"result"`
}

type ClinicOwner struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	UserID    string         `gorm:"not null;index" json:"user_id"`
	ClinicID  string         `gorm:"not null;index" json:"clinic_id"`
	Status    bool           `json:"status"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
