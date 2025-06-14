package physician

import (
	"Altheia-Backend/internal/users"
	"Altheia-Backend/pkg/utils"
	"fmt"
	"time"

	"gorm.io/gorm"

	gonanoid "github.com/matoous/go-nanoid"
)

type Service interface {
	RegisterPhysician(physician CreatePhysicianInfo) error
	UpdatePhysician(userId string, physician UpdatePhysicianInfo) error
	SoftDeletePhysician(userId string) error
	GetAllPhysiciansPaginated(page, limit int) (users.Pagination, error)
	GetPhysicianByID(id string) ([]ResultPhysicians, error)
	GetAllPhysicians() ([]ResultPhysicians, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) RegisterPhysician(physician CreatePhysicianInfo) error {
	// Validate clinic ID is provided
	if physician.ClinicID == "" {
		return fmt.Errorf("clinic_id is required")
	}

	// Check clinic exists
	exists, err := s.repo.ClinicExists(physician.ClinicID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("clinic not found")
	}

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
			ClinicID:           &physician.ClinicID,
			CreatedAt:          time.Time{},
			UpdatedAt:          time.Time{},
			DeletedAt:          gorm.DeletedAt{},
		},
	}

	newPhysician := s.repo.Create(&newUser)

	return newPhysician
}

func (s *service) UpdatePhysician(userId string, physicianData UpdatePhysicianInfo) error {
	updatedPhysician := UpdatePhysicianInfo{
		Name:               physicianData.Name,
		Phone:              physicianData.Phone,
		PhysicianSpecialty: physicianData.PhysicianSpecialty,
	}

	if physicianData.Password != "" {
		hashed, _ := utils.HashPassword(physicianData.Password)
		updatedPhysician.Password = hashed
	}

	physician := s.repo.UpdateUserAndPhysician(userId, updatedPhysician)

	return physician
}

func (s *service) SoftDeletePhysician(userId string) error {
	err := s.repo.SoftDelete(userId)
	if err != nil {
		return err
	}

	return nil

}

func (s *service) GetAllPhysicians() ([]ResultPhysicians, error) {
	var physicians []ResultPhysicians

	physicians, err := s.repo.GetAllPhysicians()
	if err != nil {
		return nil, err
	}
	return physicians, nil
}

func (s *service) GetPhysicianByID(id string) ([]ResultPhysicians, error) {
	fmt.Print("ID del usuario desde repository: ", id)
	results, err := s.repo.GetPhysicianByID(id)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *service) GetAllPhysiciansPaginated(page, limit int) (users.Pagination, error) {
	physicians, err := s.repo.GetAllPhysiciansPaginated(page, limit)
	if err != nil {
		return users.Pagination{}, err
	}
	return physicians, nil
}
