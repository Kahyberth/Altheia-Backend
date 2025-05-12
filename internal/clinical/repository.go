package clinical

import (
	"Altheia-Backend/internal/users"
	"Altheia-Backend/pkg/utils"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid"
	"gorm.io/gorm"
	"time"
)

type Repository interface {
	CreateClinic(createClinicDto CreateClinicDTO) error
	CreateEps(epsDto CreateEpsDto) error
	GetAllEps(page int, pagSize int) ([]EPS, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateClinic(createClinicDto CreateClinicDTO) error {

	tempUserPassword, _ := utils.GeneratePassword(7)
	nanoid, _ := gonanoid.Nanoid()
	hashed, _ := utils.HashPassword(tempUserPassword)

	err := r.db.Transaction(func(tx *gorm.DB) error {

		newUser := &users.User{
			ID:             nanoid,
			Name:           createClinicDto.OwnerName,
			Email:          createClinicDto.OwnerEmail,
			Password:       hashed,
			Rol:            "Admin",
			Phone:          createClinicDto.OwnerPhone,
			DocumentNumber: createClinicDto.OwenerDocumentNumber,
			Status:         false,
			Gender:         createClinicDto.OwenerGender,
			CreatedAt:      time.Time{},
			UpdatedAt:      time.Time{},
			DeletedAt:      gorm.DeletedAt{},
			LastLogin:      time.Time{},
			Patient:        users.Patient{},
			Physician:      users.Physician{},
			Receptionist:   users.Receptionist{},
		}

		userError := tx.Create(newUser).Error

		if userError != nil {
			return userError
		}

		var services []Services
		if createClinicDto.ServicesOffered != nil {
			for _, name := range createClinicDto.ServicesOffered {
				id, _ := gonanoid.Nanoid()
				services = append(services, Services{
					ID:   id,
					Name: name,
				})
			}
			if err := tx.Create(&services).Error; err != nil {
				fmt.Println("Error creating clinic service")
				return err
			}
		}

		clinicNanoId, _ := gonanoid.Nanoid()
		clinicInformationNanoId, _ := gonanoid.Nanoid()
		newClinic := &Clinic{
			ID:            clinicNanoId,
			Status:        false,
			CreatedAt:     time.Time{},
			UserID:        newUser.ID,
			UpdatedAt:     time.Time{},
			DeletedAt:     gorm.DeletedAt{},
			Physicians:    nil,
			Receptionists: nil,
			ClinicInformation: ClinicInformation{
				ClinicID:          clinicInformationNanoId,
				ClinicEmail:       createClinicDto.Email,
				ClinicName:        createClinicDto.Name,
				ClinicPhone:       createClinicDto.Phone,
				ClinicDescription: createClinicDto.Description,
				ClinicWebsite:     createClinicDto.Website,
				EmployeeCount:     createClinicDto.MemberCount,
				Services:          services,
				Photos:            nil,
				Address:           createClinicDto.Address,
				City:              createClinicDto.City,
				State:             createClinicDto.State,
				PostalCode:        createClinicDto.PostalCode,
				Country:           createClinicDto.Country,
			},
		}
		err := tx.Create(newClinic).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) CreateEps(epsDto CreateEpsDto) error {
	var eps []EPS

	for _, name := range epsDto.Eps {
		id, _ := gonanoid.Nanoid()
		eps = append(eps, EPS{
			ID:   id,
			Name: name,
		})
	}
	EpsError := r.db.Create(&eps).Error
	if EpsError != nil {
		return EpsError
	}
	return nil
}

func (r *repository) GetAllEps(page int, pagSize int) ([]EPS, error) {
	var eps []EPS
	offset := (page - 1) * pagSize
	result := r.db.Limit(pagSize).Offset(offset).Find(&eps)
	return eps, result.Error
}
