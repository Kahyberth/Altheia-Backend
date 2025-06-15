package auth

import (
	"Altheia-Backend/internal/users"
	"Altheia-Backend/pkg/utils"
	"errors"
	"fmt"
	"time"
)

type Service interface {
	Login(email, password string) (UserInfo, string, string, error)
	LoginWithActivity(email, password, userAgent, ipAddress string) (UserInfo, string, string, error)
	GetProfile(id string) (*users.User, error)
	GetUserDetails(id string) (UserDetailsResponse, error)
	GetUserLoginActivities(userID string, limit int) ([]users.LoginActivity, error)
	ChangePassword(id string, request ChangePasswordRequest) error
	DeleteUserCompletely(userID string) error
	DeactivateUser(userID string) error
	ReactivateUser(userID string) error
	verifyToken(token string) (UserInfo, string, error)
}

type service struct {
	repo Repository
}

type UserInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	ClinicId string `json:"clinic_id"`
}

type UserDetailsResponse struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Email          string                 `json:"email"`
	Role           string                 `json:"role"`
	Phone          string                 `json:"phone"`
	DocumentNumber string                 `json:"document_number"`
	Status         bool                   `json:"status"`
	Gender         string                 `json:"gender"`
	ClinicId       string                 `json:"clinic_id"`
	LastLogin      string                 `json:"last_login"`
	RoleDetails    map[string]interface{} `json:"role_details"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) getClinicIDFromUser(user *users.User) string {
	switch user.Rol {
	case "patient":
		if user.Patient.ClinicID != nil {
			return *user.Patient.ClinicID
		}
	case "physician":
		if user.Physician.ClinicID != nil {
			return *user.Physician.ClinicID
		}
	case "receptionist":
		if user.Receptionist.ClinicID != nil {
			return *user.Receptionist.ClinicID
		}
	case "owner":
		return user.ClinicOwner.ClinicID
	}
	return ""
}

func (s *service) Login(email, password string) (UserInfo, string, string, error) {

	return s.LoginWithActivity(email, password, "", "")
}

func (s *service) LoginWithActivity(email, password, userAgent, ipAddress string) (UserInfo, string, string, error) {
	user, err := s.repo.FindByEmail(email)

	if err != nil || !utils.CheckPasswordHash(password, user.Password) {
		return UserInfo{}, "", "", errors.New("invalid credentials")
	}

	if !user.Status {
		return UserInfo{}, "", "", errors.New("account is deactivated. Please contact support")
	}

	accessToken, tokenError := utils.GenerateJWT(user.ID, 1)
	if tokenError != nil {
		return UserInfo{}, "", "", tokenError
	}

	refreshToken, refreshError := utils.GenerateJWT(user.ID, 72)
	if refreshError != nil {
		return UserInfo{}, "", "", refreshError
	}

	s.repo.UpdateLastLogin(user.ID)

	s.repo.MarkAllSessionsAsInactive(user.ID)

	if userAgent != "" || ipAddress != "" {
		loginActivity := &users.LoginActivity{
			ID:               utils.GenerateNanoID(),
			UserID:           user.ID,
			DeviceType:       utils.GetDeviceTypeFromUserAgent(userAgent),
			IPAddress:        ipAddress,
			Location:         utils.GetLocationFromIP(ipAddress),
			LoginTime:        time.Now(),
			IsCurrentSession: true,
			UserAgent:        userAgent,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		s.repo.CreateLoginActivity(loginActivity)
	}

	userInfo := UserInfo{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Role:     user.Rol,
		ClinicId: s.getClinicIDFromUser(user),
	}

	return userInfo, accessToken, refreshToken, nil
}

func (s *service) GetProfile(id string) (*users.User, error) {
	return s.repo.FindByID(id)
}

func (s *service) GetUserDetails(id string) (UserDetailsResponse, error) {
	user, err := s.repo.GetUserWithAllDetails(id)
	if err != nil {
		return UserDetailsResponse{}, err
	}

	roleDetails := make(map[string]interface{})
	clinicID := s.getClinicIDFromUser(user)

	switch user.Rol {
	case "patient":
		roleDetails = map[string]interface{}{
			"patient_id":    user.Patient.ID,
			"date_of_birth": user.Patient.DateOfBirth,
			"address":       user.Patient.Address,
			"eps":           user.Patient.Eps,
			"blood_type":    user.Patient.BloodType,
		}
	case "physician":
		roleDetails = map[string]interface{}{
			"physician_id":        user.Physician.ID,
			"physician_specialty": user.Physician.PhysicianSpecialty,
			"license_number":      user.Physician.LicenseNumber,
		}
	case "receptionist":
		roleDetails = map[string]interface{}{
			"receptionist_id": user.Receptionist.ID,
		}
	case "owner":
		roleDetails = map[string]interface{}{
			"clinic_owner_id": user.ClinicOwner.ID,
		}
	}

	response := UserDetailsResponse{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		Role:           user.Rol,
		Phone:          user.Phone,
		DocumentNumber: user.DocumentNumber,
		Status:         user.Status,
		Gender:         user.Gender,
		ClinicId:       clinicID,
		LastLogin:      user.LastLogin.Format("2006-01-02 15:04:05"),
		RoleDetails:    roleDetails,
	}

	return response, nil
}

func (s *service) Logout() error {
	return nil
}

func (s *service) verifyToken(token string) (UserInfo, string, error) {
	userId, _ := utils.ValidateJWT(token)

	userData, _ := s.repo.FindByID(userId)

	fmt.Print(userData)

	userInfo := UserInfo{
		ID:       userData.ID,
		Name:     userData.Name,
		Email:    userData.Email,
		Role:     userData.Rol,
		ClinicId: s.getClinicIDFromUser(userData),
	}

	return userInfo, token, nil
}

func (s *service) ChangePassword(id string, request ChangePasswordRequest) error {

	if request.CurrentPassword == "" || request.NewPassword == "" {
		return errors.New("current password and new password are required")
	}

	user, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	if !utils.CheckPasswordHash(request.CurrentPassword, user.Password) {
		return errors.New("current password is incorrect")
	}

	hashedNewPassword, err := utils.HashPassword(request.NewPassword)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	return s.repo.ChangePassword(id, hashedNewPassword)
}

func (s *service) GetUserLoginActivities(userID string, limit int) ([]users.LoginActivity, error) {
	return s.repo.GetUserLoginActivities(userID, limit)
}

func (s *service) DeleteUserCompletely(userID string) error {
	return s.repo.DeleteUserCompletely(userID)
}

func (s *service) DeactivateUser(userID string) error {
	return s.repo.DeactivateUser(userID)
}

func (s *service) ReactivateUser(userID string) error {
	return s.repo.ReactivateUser(userID)
}
