package physician

import (
	"Altheia-Backend/internal/users"
	"Altheia-Backend/pkg/utils"
	"gorm.io/gorm"
	"time"

	gonanoid "github.com/matoous/go-nanoid"
)

type Service interface {
	RegisterPhysician(physician BasicPhysicianInfo) error
	UpdatePhysician(physician users.Physician) error
	DeletePhysician(physician users.Physician) error
	GetPhysicianByID(id string) (users.Physician, error)
	GetAllPhysicians() ([]users.Physician, error)
}

type BasicPhysicianInfo struct {
	Name                string `json:"name"`
	Email               string `json:"email"`
	Password            string `json:"password"`
	Gender              string `json:"gender"`
	Phone               string `json:"phone"`
	DocumentNumber      string `json:"document_number"`
	DateOfBirth         string `json:"date_of_birth"`
	PhysicianSpeciality string `json:"physician_specialty"`
	LicenseNumber       string `json:"license_number"`
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) RegisterPhysician(physician BasicPhysicianInfo) error {
	nanoid, _ := gonanoid.Nanoid()
	hashed, _ := utils.HashPassword(physician.Password)
	physicianNanoid, _ := gonanoid.Nanoid()

	newUser := users.User{
		ID:             nanoid,
		Name:           physician.Name,
		Email:          physician.Email,
		Password:       hashed,
		Rol:            "physician",
		Phone:          physician.Phone,
		DocumentNumber: physician.DocumentNumber,
		Status:         true,
		Gender:         physician.Gender,
		CreatedAt:      time.Time{},
		UpdatedAt:      time.Time{},
		DeletedAt:      gorm.DeletedAt{},
		LastLogin:      time.Time{},
		Physician: users.Physician{
			ID:                 physicianNanoid,
			UserID:             nanoid,
			PhysicianSpecialty: physician.PhysicianSpeciality,
			LicenseNumber:      physician.LicenseNumber,
			Status:             true,
			CreatedAt:          time.Time{},
			UpdatedAt:          time.Time{},
			DeletedAt:          gorm.DeletedAt{},
		},
	}

	newPhysician := s.repo.Create(&newUser)

	return newPhysician
}

func (s *service) UpdatePhysician(physician users.Physician) error {
	return s.repo.Update(physician)
}

func (s *service) DeletePhysician(physician users.Physician) error {
	return s.repo.Delete(physician)
}

func (s *service) GetPhysicianByID(id string) (users.Physician, error) {
	return s.repo.GetPhysicianByID(id)
}

func (s *service) GetAllPhysicians() ([]users.Physician, error) {
	return s.repo.GetAllPhysicians()
}
