package auth

import (
	"Altheia-Backend/internal/mail"
	"Altheia-Backend/internal/users/patient"
	"Altheia-Backend/pkg/utils"
	"errors"
	"fmt"
	"github.com/matoous/go-nanoid"
)

type Service interface {
	RegisterPatient(user *User) error
	Login(email, password string) (UserInfo, string, string, error)
	GetProfile(id string) (*User, error)
	verifyToken(token string) (UserInfo, string, error)
}

type service struct {
	repo Repository
}

type UserInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
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

func (s *service) Login(email, password string) (UserInfo, string, string, error) {
	user, err := s.repo.FindByEmail(email)

	if err != nil || !utils.CheckPasswordHash(password, user.Password) {
		return UserInfo{}, "", "", errors.New("invalid credentials")
	}

	accessToken, tokenError := utils.GenerateJWT(user.ID, 1)
	if tokenError != nil {
		return UserInfo{}, "", "", tokenError
	}

	refreshToken, refreshError := utils.GenerateJWT(user.ID, 72)
	if refreshError != nil {
		return UserInfo{}, "", "", refreshError
	}

	userInfo := UserInfo{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Rol,
	}

	return userInfo, accessToken, refreshToken, nil
}

func (s *service) GetProfile(id string) (*User, error) {
	return s.repo.FindByID(id)
}

func (s *service) Logout() error {
	return nil
}

func (s *service) verifyToken(token string) (UserInfo, string, error) {
	userId, _ := utils.ValidateJWT(token)

	userData, _ := s.repo.FindByID(userId)

	fmt.Print(userData)

	userInfo := UserInfo{
		ID:    userData.ID,
		Name:  userData.Name,
		Email: userData.Email,
		Role:  userData.Rol,
	}

	return userInfo, token, nil
}
