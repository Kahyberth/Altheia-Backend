package auth

import (
	"Altheia-Backend/internal/mail"
	"Altheia-Backend/internal/patient"
	"Altheia-Backend/pkg/utils"
	"errors"
	"github.com/matoous/go-nanoid"
)

type Service interface {
	RegisterPatient(user *User) error
	Login(email, password string) (string, string, error)
	GetProfile(id string) (*User, error)
	verifyToken(token string) (string, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) RegisterPatient(user *User) error {
	hashed, _ := utils.HashPassword(user.Password)
	nanoid, err := gonanoid.Nanoid()
	patientNanoid, err := gonanoid.Nanoid()
	if err != nil {
		return err
	}
	user.ID = nanoid
	user.Password = hashed

	newPatient := &patient.Patient{
		ID:             patientNanoid,
		UserID:         user.ID,
		DocumentNumber: "",
		DateOfBirth:    "",
		Address:        "",
		Eps:            "",
		BloodType:      "",
	}

	user.Patient = newPatient

	newUser := s.repo.Create(user)

	emailError := mail.SendWelcomeMessage(user.Name, []string{user.Email})
	if emailError != nil {
		return err
	}

	return newUser
}

func (s *service) Login(email, password string) (string, string, error) {
	user, err := s.repo.FindByEmail(email)

	if err != nil || !utils.CheckPasswordHash(password, user.Password) {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, tokenError := utils.GenerateJWT(user.ID, 1)

	if tokenError != nil {
	}

	refreshToken, refreshError := utils.GenerateJWT(user.ID, 72)

	if refreshError != nil {
	}

	return accessToken, refreshToken, nil
}

func (s *service) GetProfile(id string) (*User, error) {
	return s.repo.FindByID(id)
}

func (s *service) Logout() error {
	return nil
}

func (s *service) verifyToken(token string) (string, error) {
	return utils.ValidateJWT(token)
}
