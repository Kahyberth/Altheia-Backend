package auth

import (
	"Altheia-Backend/internal/mail"
	"Altheia-Backend/pkg/utils"
	"errors"
	"github.com/matoous/go-nanoid"
)

type Service interface {
	Register(user *User) error
	Login(email, password string) (string, error)
	GetProfile(id string) (*User, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) Register(user *User) error {
	hashed, _ := utils.HashPassword(user.Password)
	nanoid, err := gonanoid.Nanoid()
	if err != nil {
		return err
	}
	user.ID = nanoid
	user.Password = hashed

	newUser := s.repo.Create(user)

	emailError := mail.SendWelcomeMessage(user.Name, []string{user.Email})
	if emailError != nil {
		return err
	}

	return newUser
}

func (s *service) Login(email, password string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil || !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}
	return utils.GenerateJWT(user.ID, 0)
}

func (s *service) GetProfile(id string) (*User, error) {
	return s.repo.FindByID(id)
}
