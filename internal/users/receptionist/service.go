package receptionist

import (
	"Altheia-Backend/internal/users"
	"Altheia-Backend/pkg/utils"
	"fmt"
	"time"

	gonanoid "github.com/matoous/go-nanoid"
	"gorm.io/gorm"
)

type Service interface {
	RegisterReceptionist(receptionist CreateReceptionistInfo) error
	UpdateReceptionist(userId string, receptionistData UpdateReceptionistInfo) error
	SoftDelete(userId string) error
	GetAllReceptionistPaginated(page, limit int) (users.Pagination, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) RegisterReceptionist(receptionist CreateReceptionistInfo) error {
	if receptionist.ClinicID != "" {
		if err := s.repo.ValidateClinicExists(receptionist.ClinicID); err != nil {
			return fmt.Errorf("clinic not found: %v", err)
		}
	}

	nanoid, _ := gonanoid.Nanoid()
	receptionistNanoid, _ := gonanoid.Nanoid()
	hashed, _ := utils.HashPassword(receptionist.Password)

	newUser := users.User{
		ID:             nanoid,
		Name:           receptionist.Name,
		Email:          receptionist.Email,
		Password:       hashed,
		Rol:            "receptionist",
		Phone:          receptionist.Phone,
		DocumentNumber: receptionist.DocumentNumber,
		Status:         true,
		Gender:         receptionist.Gender,
		CreatedAt:      time.Time{},
		UpdatedAt:      time.Time{},
		DeletedAt:      gorm.DeletedAt{},
		LastLogin:      time.Time{},
		Receptionist: users.Receptionist{
			ID:     receptionistNanoid,
			UserID: nanoid,
			ClinicID: func() *string {
				if receptionist.ClinicID != "" {
					return &receptionist.ClinicID
				}
				return nil
			}(),
			Status: true,
		},
	}

	newReceptionist := s.repo.Create(&newUser)

	return newReceptionist
}

func (s *service) UpdateReceptionist(userId string, receptionistData UpdateReceptionistInfo) error {
	hashed, _ := utils.HashPassword(receptionistData.Password)

	updatedReceptionist := UpdateReceptionistInfo{
		Name:     receptionistData.Name,
		Password: hashed,
		Phone:    receptionistData.Phone,
	}

	receptionist := s.repo.UpdateUserAndReceptionist(userId, updatedReceptionist)

	return receptionist
}

func (s *service) SoftDelete(userId string) error {
	err := s.repo.SoftDelete(userId)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetAllReceptionistPaginated(page, limit int) (users.Pagination, error) {
	receptionists, err := s.repo.GetAllReceptionistPaginated(page, limit)
	if err != nil {
		return users.Pagination{}, err
	}
	return receptionists, nil
}
