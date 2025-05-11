package physician

import (
	"Altheia-Backend/internal/users"
	"gorm.io/gorm"
	"time"
)

type Repository interface {
	Create(user *users.User) error
	Update(physician users.Physician) error
	GetPhysicianByID(id string) (users.Physician, error)
	UpdatePhysicianWithUser(physician users.Physician, userUpdates map[string]interface{}) error
	Delete(physician users.Physician) error
	SoftDelete(physician users.Physician) error
	GetAllPhysicians() ([]users.Physician, error)
	GetUserByID(id string) (users.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository { return &repository{db} }

func (r *repository) Create(user *users.User) error {
	return r.db.Create(user).Error
}

func (r *repository) Update(physician users.Physician) error {
	return r.db.Save(physician).Error
}

func (r *repository) Delete(physician users.Physician) error {
	return r.db.Delete(physician).Error
}

func (r *repository) SoftDelete(physician users.Physician) error {
	return r.db.Delete(physician).Update("deleted_at", time.Now()).Error
}

func (r *repository) GetPhysicianByID(id string) (users.Physician, error) {
	var physician users.Physician
	err := r.db.Where("id = ?", id).First(&physician).Error
	if err != nil {
		return users.Physician{}, err
	}
	return physician, nil
}

func (r *repository) GetAllPhysicians() ([]users.Physician, error) {
	var physicians []users.Physician
	err := r.db.Find(&physicians).Error
	if err != nil {
		return nil, err
	}
	return physicians, nil
}

func (r *repository) UpdatePhysicianWithUser(physician users.Physician, userUpdates map[string]interface{}) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Save(physician).Error; err != nil {
			return err
		}

		// If there are user updates, update the user table
		if len(userUpdates) > 0 {
			if err := tx.Table("users").Where("id = ?", physician.UserID).Updates(userUpdates).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *repository) GetUserByID(id string) (users.User, error) {
	var user users.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return users.User{}, err
	}
	return user, nil
}
