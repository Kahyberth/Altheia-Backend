package patient

import (
	"Altheia-Backend/internal/users"
	"Altheia-Backend/pkg/utils"
	gonanoid "github.com/matoous/go-nanoid"
	"gorm.io/gorm"
	"time"
)

type Service interface {
	RegisterPatient(patient CreatePatientInfo) error
	UpdatePatient(userId string, patientData UpdatePatientInfo) error
	SoftDeletePatient(userId string) error
	GetAllPatientsPaginated(page, limit int) (users.Pagination, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) RegisterPatient(patient CreatePatientInfo) error {
	nanoid, _ := gonanoid.Nanoid()
	patientNanoid, _ := gonanoid.Nanoid()
	hashed, _ := utils.HashPassword(patient.Password)

	newUser := users.User{
		ID:             nanoid,
		Name:           patient.Name,
		Email:          patient.Email,
		Password:       hashed,
		Rol:            "patient",
		Phone:          patient.Phone,
		DocumentNumber: patient.DocumentNumber,
		Status:         true,
		Gender:         patient.Gender,
		CreatedAt:      time.Time{},
		UpdatedAt:      time.Time{},
		DeletedAt:      gorm.DeletedAt{},
		LastLogin:      time.Time{},
		Patient: users.Patient{
			ID:          patientNanoid,
			UserID:      nanoid,
			DateOfBirth: patient.DateOfBirth,
			Address:     patient.Address,
			Eps:         patient.Eps,
			BloodType:   patient.BloodType,
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
			DeletedAt:   gorm.DeletedAt{},
		},
	}

	newPatient := s.repository.Create(&newUser)

	return newPatient

}

func (s *service) UpdatePatient(userId string, patientData UpdatePatientInfo) error {
	hashed, _ := utils.HashPassword(patientData.Password)

	updatedPatient := UpdatePatientInfo{
		Name:     patientData.Name,
		Password: hashed,
		Phone:    patientData.Phone,
		Address:  patientData.Address,
		Eps:      patientData.Eps,
	}
	patient := s.repository.UpdateUserAndPatient(userId, updatedPatient)

	return patient
}

func (s *service) SoftDeletePatient(userId string) error {
	err := s.repository.SoftDelete(userId)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetAllPatientsPaginated(limit, page int) (users.Pagination, error) {
	patients, err := s.repository.GetAllPatientsPaginated(limit, page)
	if err != nil {
		return users.Pagination{}, err
	}
	return patients, nil
}
