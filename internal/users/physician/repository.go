package physician

import (
	"Altheia-Backend/internal/users"
	"gorm.io/gorm"
	"time"
)

type Repository interface {
	Create(user *users.User) error
	UpdateUserAndPhysician(UserId string, Info UpdatePhysicianInfo) error
	GetAllPhysiciansPaginated(page, limit int) (users.Pagination, error)
	SoftDelete(userId string) error
	GetPhysicianByID(id string) ([]ResultPhysicians, error)
	GetAllPhysicians() ([]users.Physician, error)
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

func (r *repository) GetAllPhysicians() ([]users.Physician, error) {
	var physicians []users.Physician

	err := r.db.Find(&physicians).Error
	if err != nil {
		return nil, err
	}

	return physicians, nil
}

func (r *repository) GetAllPhysiciansPaginated(page, limit int) (users.Pagination, error) {

	var physicians []users.Physician
	var totalRows int64

	pagination := users.Pagination{
		Limit: limit,
		Page:  page,
	}

	offset := (pagination.Page - 1) * pagination.Limit

	r.db.Model(&users.Physician{}).Count(&totalRows)

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

func (r *repository) GetPhysicianByID(id string) ([]ResultPhysicians, error) {
	var result []ResultPhysicians

	err := r.db.Table("users").
		Select("users.id, users.name, users.email, users.rol, users.status, users.gender, users.last_login, physicians.physician_specialty").
		Joins("INNER JOIN physicians ON physicians.user_id = users.id").
		Where("users.id = ?", id).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}
