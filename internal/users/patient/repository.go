package patient

import (
	"Altheia-Backend/internal/users"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Create(user *users.User) error
	UpdateUserAndPatient(UserId string, Info UpdatePatientInfo) error
	SoftDelete(userId string) error
	GetAllPatientsPaginated(page, limit int) (users.Pagination, error)
	GetAllPatients() ([]users.Patient, error)
	GetPatientByClinicId(clinicId string) ([]users.Patient, error)
	GetPatientByClinicIdPaginated(clinicId string, page, limit int) (users.Pagination, error)
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

func (r *repository) GetAllPatients() ([]users.Patient, error) {
	type Result struct {
		users.Patient
		Name string `json:"name"`
	}
	var results []Result

	err := r.db.Model(&users.Patient{}).
		Select("patients.*, users.name as name").
		Joins("JOIN users ON patients.user_id = users.id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	patients := make([]users.Patient, len(results))
	for i, result := range results {
		patients[i] = result.Patient
		patients[i].Name = result.Name
	}

	return patients, nil
}

func (r *repository) GetPatientByClinicId(clinicId string) ([]users.Patient, error) {
	var patients []users.Patient

	err := r.db.Preload("User").Where("clinic_id = ?", clinicId).Find(&patients).Error
	if err != nil {
		return nil, err
	}

	return patients, nil
}

func (r *repository) GetPatientByClinicIdPaginated(clinicId string, page, limit int) (users.Pagination, error) {
	var patients []users.Patient
	var totalRows int64

	pagination := users.Pagination{
		Limit: limit,
		Page:  page,
	}

	offset := (pagination.Page - 1) * pagination.Limit

	r.db.Model(&users.Patient{}).Where("clinic_id = ?", clinicId).Count(&totalRows)

	if err := r.db.Preload("User").Where("clinic_id = ?", clinicId).
		Limit(pagination.Limit).
		Offset(offset).
		Find(&patients).Error; err != nil {
		return users.Pagination{}, err
	}

	pagination.Total = totalRows
	pagination.Result = patients
	return pagination, nil
}
