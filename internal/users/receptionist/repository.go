package receptionist

import (
	"Altheia-Backend/internal/users"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Create(user *users.User) error
	ValidateClinicExists(clinicID string) error
	UpdateUserAndReceptionist(UserId string, Info UpdateReceptionistInfo) error
	SoftDelete(userId string) error
	GetAllReceptionistPaginated(page, limit int) (users.Pagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(user *users.User) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Create(user).Error
}

func (r *repository) ValidateClinicExists(clinicID string) error {
	var count int64
	if err := r.db.Table("clinics").Where("id = ?", clinicID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *repository) UpdateUserAndReceptionist(UserId string, Info UpdateReceptionistInfo) error {
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
		return nil
	})
}

func (r *repository) SoftDelete(userId string) error {
	receptionist := users.Receptionist{
		DeletedAt: gorm.DeletedAt{
			Time:  time.Now(),
			Valid: true,
		},
		Status: false,
	}

	return r.db.Model(&users.Receptionist{}).Where("user_id = ?", userId).Updates(receptionist).Error
}

func (r *repository) GetAllReceptionistPaginated(page, limit int) (users.Pagination, error) {

	var receptionist []users.Receptionist
	var totalRows int64

	pagination := users.Pagination{
		Limit: limit,
		Page:  page,
	}

	offset := (pagination.Page - 1) * pagination.Limit

	r.db.Model(&users.Receptionist{}).Count(&totalRows)

	err := r.db.Limit(pagination.Limit).Offset(offset).Find(&receptionist).Error
	if err != nil {
		return users.Pagination{}, err
	}

	pagination.Total = totalRows
	pagination.Result = receptionist

	return pagination, nil
}
