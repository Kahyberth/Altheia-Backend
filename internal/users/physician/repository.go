package physician

import (
	"Altheia-Backend/internal/users"
	"gorm.io/gorm"
	"time"
)

type Repository interface {
	Create(user *users.User) error
	GetPhysicianByID(id string) (users.Physician, error)
	UpdateUserAndPhysician(UserId string, Info UpdatePhysicianInfo) error
	GetAllPhysiciansPaginated(page, limit int) (users.Pagination, error)
	SoftDelete(userId string) error
	GetUserByID(id string) (users.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository { return &repository{db} }

func (r *repository) Create(user *users.User) error {
	return r.db.Create(user).Error
}

func (r *repository) UpdateUserAndPhysician(UserId string, Info UpdatePhysicianInfo) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&users.User{}).Where("id = ?", UserId).
			Updates(map[string]interface{}{
				"name":     Info.Name,
				"password": Info.Password,
				"phone":    Info.Phone,
			}).Error; err != nil {
		}

		if err := tx.Model(&users.Physician{}).Where("id = ?", UserId).
			Updates(map[string]interface{}{
				"physician_specialty": Info.PhysicianSpecialty,
			}).Error; err != nil {
		}

		return nil
	})

}

func (r *repository) SoftDelete(userId string) error {
	physician := users.Physician{
		DeletedAt: gorm.DeletedAt{
			Time:  time.Now(),
			Valid: true,
		},
		Status: false,
	}
	return r.db.Model(&users.Physician{}).Where("id = ?", userId).Updates(physician).Error
}

func (r *repository) GetPhysicianByID(id string) (users.Physician, error) {
	var physician users.Physician
	err := r.db.Where("id = ?", id).First(&physician).Error
	if err != nil {
		return users.Physician{}, err
	}
	return physician, nil
}

func (r *repository) GetAllPhysiciansPaginated(page, limit int) (users.Pagination, error) {

	var physicians []users.Physician
	var totalRows int64

	pagination := users.Pagination{
		Limit: limit,
		Page:  page,
	}

	offset := (pagination.Page - 1) * pagination.Limit

	// Obtener total de registros
	r.db.Model(&users.Physician{}).Count(&totalRows)

	// Obtener registros paginados
	err := r.db.Limit(pagination.Limit).
		Offset(offset).
		Find(&physicians).Error

	if err != nil {
		return users.Pagination{}, err
	}

	pagination.Total = totalRows
	pagination.Result = physicians

	return pagination, nil
}

func (r *repository) GetUserByID(id string) (users.User, error) {
	var user users.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return users.User{}, err
	}
	return user, nil
}
