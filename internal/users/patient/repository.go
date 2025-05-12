package patient

import (
	"Altheia-Backend/internal/users"
	"gorm.io/gorm"
	"time"
)

type Repository interface {
	Create(user *users.User) error
	UpdateUserAndPatient(UserId string, Info UpdatePatientInfo) error
	SoftDelete(userId string) error
	GetAllPatientsPaginated(page, limit int) (users.Pagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(user *users.User) error {
	r.db.Create(user)
	return nil
}

func (r *repository) UpdateUserAndPatient(UserId string, Info UpdatePatientInfo) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&users.User{}).Where("id = ?", UserId).
			Updates(map[string]interface{}{
				"name":     Info.Name,
				"password": Info.Password,
				"phone":    Info.Phone,
			}).Error; err != nil {
		}

		if err := tx.Model(&users.Patient{}).Where("id = ?", UserId).
			Updates(map[string]interface{}{
				"eps":     Info.Eps,
				"address": Info.Address,
			}).Error; err != nil {
		}

		return nil
	})
}

func (r *repository) SoftDelete(userId string) error {
	patient := users.Patient{
		DeletedAt: gorm.DeletedAt{
			Time:  time.Now(),
			Valid: true,
		},
		Status: false,
	}

	return r.db.Model(&users.Patient{}).Where("id = ?", userId).Updates(patient).Error
}

func (r *repository) GetAllPatientsPaginated(page, limit int) (users.Pagination, error) {
	var patients []users.Patient
	var totalRows int64

	Pagination := users.Pagination{
		Limit: limit,
		Page:  page,
	}

	offset := (Pagination.Page - 1) * Pagination.Limit

	r.db.Model(&users.Patient{}).Count(&totalRows)

	err := r.db.Limit(Pagination.Limit).Offset(offset).Find(&patients).Error
	if err != nil {
		return users.Pagination{}, err
	}
	Pagination.Total = totalRows
	Pagination.Result = patients

	return Pagination, nil
}
