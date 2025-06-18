package superAdmin

import (
	"Altheia-Backend/internal/users"
	"Altheia-Backend/pkg/utils"
	"fmt"
	"time"

	gonanoid "github.com/matoous/go-nanoid"
)

type Service interface {
	RegisterSuperAdmin(superAdmin CreateSuperAdminInfo) error
	UpdateSuperAdmin(userID string, superAdminData UpdateSuperAdminInfo) error
	GetSuperAdminByID(userID string) (SuperAdminResponse, error)
	GetAllSuperAdminsPaginated(page, limit int) (users.Pagination, error)
	SoftDelete(userID string) error
	GetDeactivatedUsersPaginated(page, limit int) (users.Pagination, error)
	GetClinicOwnersPaginated(page, limit int) (users.Pagination, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) RegisterSuperAdmin(superAdmin CreateSuperAdminInfo) error {
	nanoid, _ := gonanoid.Nanoid()
	superAdminNanoid, _ := gonanoid.Nanoid()
	hashed, _ := utils.HashPassword(superAdmin.Password)

	permissions := superAdmin.Permissions
	if permissions == "" {
		permissions = `{
			"users": {"create": true, "read": true, "update": true, "delete": true},
			"clinics": {"create": true, "read": true, "update": true, "delete": true},
			"appointments": {"create": true, "read": true, "update": true, "delete": true},
			"medical_records": {"create": true, "read": true, "update": true, "delete": true},
			"reports": {"create": true, "read": true, "update": true, "delete": true},
			"system": {"manage": true, "backup": true, "restore": true}
		}`
	}

	newUser := users.User{
		ID:             nanoid,
		Name:           superAdmin.Name,
		Email:          superAdmin.Email,
		Password:       hashed,
		Rol:            "super-admin",
		Phone:          superAdmin.Phone,
		DocumentNumber: superAdmin.DocumentNumber,
		Status:         true,
		Gender:         superAdmin.Gender,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		LastLogin:      time.Time{},
		SuperAdmin: users.SuperAdmin{
			ID:          superAdminNanoid,
			UserID:      nanoid,
			Permissions: permissions,
			Status:      true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	return s.repo.Create(&newUser)
}

func (s *service) UpdateSuperAdmin(userID string, superAdminData UpdateSuperAdminInfo) error {
	if err := s.repo.ValidateUserExists(userID); err != nil {
		return fmt.Errorf("super admin not found: %v", err)
	}

	return s.repo.Update(userID, superAdminData)
}

func (s *service) GetSuperAdminByID(userID string) (SuperAdminResponse, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return SuperAdminResponse{}, fmt.Errorf("super admin not found: %v", err)
	}

	if user.Rol != "super-admin" {
		return SuperAdminResponse{}, fmt.Errorf("user is not a super admin")
	}

	response := SuperAdminResponse{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		Phone:          user.Phone,
		DocumentNumber: user.DocumentNumber,
		Gender:         user.Gender,
		Status:         user.Status,
		Permissions:    user.SuperAdmin.Permissions,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
		LastLogin:      user.LastLogin,
	}

	return response, nil
}

func (s *service) GetAllSuperAdminsPaginated(page, limit int) (users.Pagination, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	superAdmins, total, err := s.repo.GetAll(page, limit)
	if err != nil {
		return users.Pagination{}, err
	}

	var responses []SuperAdminResponse
	for _, user := range superAdmins {
		response := SuperAdminResponse{
			ID:             user.ID,
			Name:           user.Name,
			Email:          user.Email,
			Phone:          user.Phone,
			DocumentNumber: user.DocumentNumber,
			Gender:         user.Gender,
			Status:         user.Status,
			Permissions:    user.SuperAdmin.Permissions,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
			LastLogin:      user.LastLogin,
		}
		responses = append(responses, response)
	}

	pagination := users.Pagination{
		Limit:  limit,
		Page:   page,
		Total:  total,
		Result: responses,
	}

	return pagination, nil
}

func (s *service) SoftDelete(userID string) error {
	if err := s.repo.ValidateUserExists(userID); err != nil {
		return fmt.Errorf("super admin not found: %v", err)
	}

	return s.repo.SoftDelete(userID)
}

func (s *service) GetDeactivatedUsersPaginated(page, limit int) (users.Pagination, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	deactivatedUsers, total, err := s.repo.GetDeactivatedUsers(page, limit)
	if err != nil {
		return users.Pagination{}, fmt.Errorf("failed to retrieve deactivated users: %v", err)
	}

	var responses []DeactivatedUserResponse
	for _, user := range deactivatedUsers {
		roleDetails := make(map[string]interface{})

		switch user.Rol {
		case "patient":
			roleDetails = map[string]interface{}{
				"patient_id":    user.Patient.ID,
				"date_of_birth": user.Patient.DateOfBirth,
				"address":       user.Patient.Address,
				"eps":           user.Patient.Eps,
				"blood_type":    user.Patient.BloodType,
				"clinic_id":     user.Patient.ClinicID,
			}
		case "physician":
			roleDetails = map[string]interface{}{
				"physician_id":        user.Physician.ID,
				"physician_specialty": user.Physician.PhysicianSpecialty,
				"license_number":      user.Physician.LicenseNumber,
				"clinic_id":           user.Physician.ClinicID,
			}
		case "receptionist":
			roleDetails = map[string]interface{}{
				"receptionist_id": user.Receptionist.ID,
				"clinic_id":       user.Receptionist.ClinicID,
			}
		case "owner":
			roleDetails = map[string]interface{}{
				"clinic_owner_id": user.ClinicOwner.ID,
				"clinic_id":       user.ClinicOwner.ClinicID,
			}
		case "lab_technician":
			roleDetails = map[string]interface{}{
				"lab_technician_id": user.LabTechnician.ID,
				"clinic_id":         user.LabTechnician.ClinicID,
			}
		case "super-admin":
			roleDetails = map[string]interface{}{
				"super_admin_id": user.SuperAdmin.ID,
				"permissions":    user.SuperAdmin.Permissions,
				"access_level":   "system_admin",
			}
		}

		response := DeactivatedUserResponse{
			ID:             user.ID,
			Name:           user.Name,
			Email:          user.Email,
			Role:           user.Rol,
			Phone:          user.Phone,
			DocumentNumber: user.DocumentNumber,
			Gender:         user.Gender,
			Status:         user.Status,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
			LastLogin:      user.LastLogin,
			RoleDetails:    roleDetails,
		}
		responses = append(responses, response)
	}

	pagination := users.Pagination{
		Limit:  limit,
		Page:   page,
		Total:  total,
		Result: responses,
	}

	return pagination, nil
}

func (s *service) GetClinicOwnersPaginated(page, limit int) (users.Pagination, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	clinicOwners, total, err := s.repo.GetClinicOwners(page, limit)
	if err != nil {
		return users.Pagination{}, fmt.Errorf("failed to retrieve clinic owners: %v", err)
	}

	var responses []ClinicOwnerResponse
	for _, user := range clinicOwners {
		response := ClinicOwnerResponse{
			ID:             user.ID,
			Name:           user.Name,
			Email:          user.Email,
			Phone:          user.Phone,
			DocumentNumber: user.DocumentNumber,
			Gender:         user.Gender,
			Status:         user.Status,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
			LastLogin:      user.LastLogin,
			ClinicOwnerID:  user.ClinicOwner.ID,
			ClinicID:       user.ClinicOwner.ClinicID,
			OwnerStatus:    user.ClinicOwner.Status,
			OwnerCreatedAt: user.ClinicOwner.CreatedAt,
		}
		responses = append(responses, response)
	}

	pagination := users.Pagination{
		Limit:  limit,
		Page:   page,
		Total:  total,
		Result: responses,
	}

	return pagination, nil
}
