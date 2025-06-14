package physician

import (
	"Altheia-Backend/internal/clinical"
	"Altheia-Backend/internal/users"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Create(user *users.User) error
	UpdateUserAndPhysician(UserId string, Info UpdatePhysicianInfo) error
	GetAllPhysiciansPaginated(page, limit int) (users.Pagination, error)
	SoftDelete(userId string) error
	GetPhysicianByID(id string) ([]ResultPhysicians, error)
	GetAllPhysicians() ([]ResultPhysicians, error)
	ClinicExists(id string) (bool, error)
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

		userUpdates := map[string]interface{}{
			"name":  Info.Name,
			"phone": Info.Phone,
		}

		if Info.Password != "" {
			userUpdates["password"] = Info.Password
		}

		if err := tx.Model(&users.User{}).Where("id = ?", UserId).
			Updates(userUpdates).Error; err != nil {
			return err
		}

		if err := tx.Model(&users.Physician{}).Where("user_id = ?", UserId).
			Updates(map[string]interface{}{
				"physician_specialty": Info.PhysicianSpecialty,
			}).Error; err != nil {
			return err
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
	return r.db.Model(&users.Physician{}).Where("user_id = ?", userId).Updates(physician).Error
}

func (r *repository) GetAllPhysicians() ([]ResultPhysicians, error) {
	var physicians []ResultPhysicians

	err := r.db.Table("physicians").
		Select("users.id as user_id, physicians.id as physician_id, users.name, users.email, users.rol, users.status, users.gender, users.last_login, physicians.physician_specialty").
		Joins("INNER JOIN users ON users.id = physicians.user_id").
		Scan(&physicians).Error
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

func (r *repository) ClinicExists(id string) (bool, error) {
	var count int64
	if err := r.db.Model(&clinical.Clinic{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
