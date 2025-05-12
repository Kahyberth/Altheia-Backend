package receptionist

import (
	"Altheia-Backend/internal/users"
	"gorm.io/gorm"
	"time"
)

type Repository interface {
	Create(user *users.User) error
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
	return r.db.Create(user).Error
}

func (r *repository) UpdateUserAndReceptionist(UserId string, Info UpdateReceptionistInfo) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&users.User{}).Where("id = ?", UserId).
			Updates(map[string]interface{}{
				"name":     Info.Name,
				"password": Info.Password,
				"phone":    Info.Phone,
			}).Error; err != nil {
		}

		// No hay necesidad de actualizar la tabla receptionist, ya que no hay campos adicionales

		// if err := tx.Model(&users.Receptionist{}).Where("id = ?", UserId).
		//	Updates(map[string]interface{}{}).Error; err != nil {
		// }

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

	return r.db.Model(&users.Receptionist{}).Where("id = ?", userId).Updates(receptionist).Error
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
