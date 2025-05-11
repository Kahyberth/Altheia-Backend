package auth

import (
	"Altheia-Backend/internal/physician"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id string) (*User, error)
}

type UpdateInfo struct {
	Name               string
	Email              string
	Password           string
	Gender             string
	PhysicianSpecialty string
	Phone              string
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *repository) FindByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *repository) FindByID(id string) (*User, error) {
	var user User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *repository) Update(user *User) error {
	return r.db.Save(user).Error
}

func (r *repository) FilterByIdPhysician(physicianID string) (physician.Physician, error) {
	var physicianInfo physician.Physician
	r.db.Where("id = ?", physicianID).First(&physicianInfo)
	return physicianInfo, nil
}

func (r *repository) UpdatePhysician(user *physician.Physician) error {
	return r.db.Save(user).Error
}

func (r *repository) UpdateUserAndPhysician(UserId string, Info UpdateInfo) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&User{}).Where("id = ?", UserId).
			Updates(map[string]interface{}{
				"name":     Info.Name,
				"email":    Info.Email,
				"password": Info.Password,
				"gender":   Info.Gender,
			}).Error; err != nil {
		}

		if err := tx.Model(&physician.Physician{}).Where("id = ?", UserId).
			Updates(map[string]interface{}{
				"physicianSpecialty": Info.PhysicianSpecialty,
				"phone":              Info.Phone,
			}).Error; err != nil {
		}

		return nil
	})

}
