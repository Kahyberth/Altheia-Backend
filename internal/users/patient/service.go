package patient

import (
	"Altheia-Backend/internal/users"
	"Altheia-Backend/pkg/utils"
	"time"

	gonanoid "github.com/matoous/go-nanoid"
	"gorm.io/gorm"
)

type Service interface {
	RegisterPatient(patient CreatePatientInfo) error
	UpdatePatient(userId string, patientData UpdatePatientInfo) error
	SoftDeletePatient(userId string) error
	GetAllPatientsPaginated(page, limit int) (users.Pagination, error)
	GetAllPatients() ([]users.Patient, error)
	GetPatientByClinicId(clinicId string) ([]users.Patient, error)
	GetPatientByClinicIdPaginated(clinicId string, page, limit int) (users.Pagination, error)
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
			ID:     patientNanoid,
			UserID: nanoid,
			ClinicID: func() *string {
				if patient.ClinicID != "" {
					return &patient.ClinicID
				}
				return nil
			}(),
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
	updatedPatient := UpdatePatientInfo{
		Name:    patientData.Name,
		Phone:   patientData.Phone,
		Address: patientData.Address,
		Eps:     patientData.Eps,
	}

	if patientData.Password != "" {
		hashed, _ := utils.HashPassword(patientData.Password)
		updatedPatient.Password = hashed
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

func (s *service) GetAllPatients() ([]users.Patient, error) {
	patients, err := s.repository.GetAllPatients()
	if err != nil {
		return nil, err
	}
	return patients, nil
}

func (s *service) GetPatientByClinicId(clinicId string) ([]users.Patient, error) {
	patients, err := s.repository.GetPatientByClinicId(clinicId)
	if err != nil {
		return nil, err
	}
	return patients, nil
}

func (s *service) GetPatientByClinicIdPaginated(clinicId string, page, limit int) (users.Pagination, error) {
	return s.repository.GetPatientByClinicIdPaginated(clinicId, page, limit)
}
