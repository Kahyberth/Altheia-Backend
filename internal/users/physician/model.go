package physician

import (
	"gorm.io/gorm"
	"time"
)

type Physician struct {
	ID                 string         `gorm:"primaryKey" json:"id"`
	UserID             string         `gorm:"not null;index" json:"user_id"`
	PhysicianSpecialty string         `json:"physician_specialty"`
	LicenseNumber      string         `json:"license_number"`
	Status             bool           `json:"status"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}
