package physician

import (
	"Altheia-Backend/internal/mail"
	"Altheia-Backend/internal/users"
	"Altheia-Backend/pkg/utils"
	gonanoid "github.com/matoous/go-nanoid"
)

type Service interface {
	RegisterPhysician(user *users.User) error
	UpdatePhysician(physician users.Physician) error
	DeletePhysician(physician users.Physician) error
	GetPhysicianByID(id string) (users.Physician, error)
	GetAllPhysicians() ([]users.Physician, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) RegisterPhysician(user *users.User) error {
	hashed, _ := utils.HashPassword(user.Password)
	nanoid, err := gonanoid.Nanoid()
	if err != nil {
		return err
	}
	physicianNanoid, err := gonanoid.Nanoid()
	if err != nil {
		return err
	}
	user.ID = nanoid
	user.Password = hashed

	newPhysician := users.Physician{
		ID:                 physicianNanoid,
		UserID:             user.ID,
		PhysicianSpecialty: "",
		LicenseNumber:      "",
	}

	user.Physician = newPhysician

	newUser := s.repo.Create(user)

	emailError := mail.SendWelcomeMessage(user.Name, []string{user.Email})
	if emailError != nil {
		return err
	}

	return newUser
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
