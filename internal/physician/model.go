package physician

import (
	"gorm.io/gorm"
	"time"
)

type Physician struct {
	ID                 string         `gorm:"primaryKey" json:"id"`
	UserID             string         `json:"user_id"`
	Gender             string         `json:"gender"`
	PhysicianSpecialty string         `json:"physician_specialty"`
	LicenseNumber      string         `json:"license_number"`
	Phone              string         `json:"phone"`
	Status             bool           `json:"status"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}
