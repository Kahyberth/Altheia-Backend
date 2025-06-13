package clinical

import (
	"Altheia-Backend/internal/users"
	"Altheia-Backend/pkg/utils"
	"fmt"
	"time"

	gonanoid "github.com/matoous/go-nanoid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateClinic(createClinicDto CreateClinicDTO) error
	CreateEps(epsDto CreateEpsDto) error
	GetAllEps(page int, pagSize int) ([]EPS, error)
	GetAllServices(page int, pagSize int) ([]ServicesOffered, error)
	CreateServices(dto CreateServicesDto) error
	GetClinicByOwnerID(ownerID string) (*ClinicCompleteInfoResponse, error)
	AssignServicesToClinic(dto AssignServicesClinicDTO) error
	GetClinicsByEps(epsID string, page int, pageSize int) ([]Clinic, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateClinic(createClinicDto CreateClinicDTO) error {

	tempUserPassword, _ := utils.GeneratePassword(10)
	nanoid, _ := gonanoid.Nanoid()
	fmt.Println(tempUserPassword)
	hashed, _ := utils.HashPassword(tempUserPassword)

	err := r.db.Transaction(func(tx *gorm.DB) error {

		newUser := &users.User{
			ID:             nanoid,
			Name:           createClinicDto.OwnerName,
			Email:          createClinicDto.OwnerEmail,
			Password:       hashed,
			Rol:            "owner",
			Phone:          createClinicDto.OwnerPhone,
			DocumentNumber: createClinicDto.OwnerDocumentNumber,
			Status:         false,
			Gender:         createClinicDto.OwnerGender,
			CreatedAt:      time.Time{},
			UpdatedAt:      time.Time{},
			DeletedAt:      gorm.DeletedAt{},
			LastLogin:      time.Time{},
			Patient:        users.Patient{},
			Physician:      users.Physician{},
			Receptionist:   users.Receptionist{},
			ClinicOwner:    users.ClinicOwner{},
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
				ServicesOffered:   nil,
				EpsOffered:        nil,
				Photos:            nil,
				Address:           createClinicDto.Address,
				City:              createClinicDto.City,
				State:             createClinicDto.State,
				PostalCode:        createClinicDto.PostalCode,
				Country:           createClinicDto.Country,
			},
		}

		clinicOwnerId, _ := gonanoid.Nanoid()
		newUser.ClinicOwner = users.ClinicOwner{
			ID:        clinicOwnerId,
			UserID:    nanoid,
			ClinicID:  newClinic.ID,
			Status:    false,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		}

		userError := tx.Create(newUser).Error
		if userError != nil {
			return userError
		}

		err := tx.Create(newClinic).Error
		if err != nil {
			return err
		}

		var services []ServicesOffered
		if len(createClinicDto.ServicesOffered) > 0 {
			if servicesError := tx.
				Where("id IN ?", createClinicDto.ServicesOffered).
				Find(&services).Error; servicesError != nil {
				return servicesError
			}

			if len(services) != len(createClinicDto.ServicesOffered) {
				return fmt.Errorf("one or more service IDs are invalid")
			}
		}

		if len(services) > 0 {
			if clinicError := tx.Model(&newClinic.ClinicInformation).
				Association("ServicesOffered").
				Append(services); clinicError != nil {
				return clinicError
			}
		}

		var acceptedEps []EPS
		if len(createClinicDto.AcceptedEPS) > 0 {
			if epsError := tx.
				Where("id IN ?", createClinicDto.AcceptedEPS).
				Find(&acceptedEps).Error; epsError != nil {
				return epsError
			}

			if len(acceptedEps) != len(createClinicDto.AcceptedEPS) {
				return fmt.Errorf("one or more EPS IDs are invalid")
			}
		}

		if len(acceptedEps) > 0 {
			if clinicError := tx.Model(&newClinic.ClinicInformation).
				Association("EpsOffered").
				Append(acceptedEps); clinicError != nil {
				return clinicError
			}
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

func (r *repository) CreateServices(dto CreateServicesDto) error {
	var services []ServicesOffered
	for _, name := range dto.ServicesOffered {
		id, _ := gonanoid.Nanoid()
		services = append(services, ServicesOffered{
			ID:   id,
			Name: name,
		})
	}
	if err := r.db.Create(&services).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) GetAllEps(page int, pagSize int) ([]EPS, error) {
	var eps []EPS
	offset := (page - 1) * pagSize
	result := r.db.Limit(pagSize).Offset(offset).Find(&eps)
	return eps, result.Error
}

func (r *repository) GetAllServices(page int, pagSize int) ([]ServicesOffered, error) {
	var servicesOffered []ServicesOffered
	offset := (page - 1) * pagSize
	result := r.db.Limit(pagSize).Offset(offset).Find(&servicesOffered)
	return servicesOffered, result.Error
}

func (r *repository) GetClinicByOwnerID(ownerID string) (*ClinicCompleteInfoResponse, error) {
	var clinic Clinic
	var owner users.User

	var clinicOwner users.ClinicOwner
	if err := r.db.Where("user_id = ?", ownerID).First(&clinicOwner).Error; err != nil {
		return nil, fmt.Errorf("clinic owner not found: %v", err)
	}

	if err := r.db.
		Preload("ClinicInformation.ServicesOffered").
		Preload("ClinicInformation.EpsOffered").
		Preload("ClinicInformation.Photos").
		Preload("Physicians.User").
		Preload("Receptionists.User").
		Preload("LabTechnicians.User").
		Preload("Patients.User").
		Where("id = ?", clinicOwner.ClinicID).
		First(&clinic).Error; err != nil {
		return nil, fmt.Errorf("clinic not found: %v", err)
	}

	if err := r.db.Where("id = ?", ownerID).First(&owner).Error; err != nil {
		return nil, fmt.Errorf("owner user not found: %v", err)
	}

	response := &ClinicCompleteInfoResponse{
		Clinic:      clinic,
		Owner:       owner,
		Information: clinic.ClinicInformation,
	}

	return response, nil
}

func (r *repository) AssignServicesToClinic(dto AssignServicesClinicDTO) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		var clinicInfo ClinicInformation
		if err := tx.Where("clinic_id = ?", dto.ClinicID).
			First(&clinicInfo).Error; err != nil {
			return fmt.Errorf("clinic information not found: %v", err)
		}

		var services []ServicesOffered
		if len(dto.Services) > 0 {
			if err := tx.Where("id IN ?", dto.Services).Find(&services).Error; err != nil {
				return fmt.Errorf("services not found: %v", err)
			}
		}

		if err := tx.Model(&clinicInfo).Association("ServicesOffered").Replace(services); err != nil {
			return fmt.Errorf("failed to associate services: %v", err)
		}
		return nil
	})
}

// GetClinicsByEps returns a paginated list of clinics that accept the given EPS
func (r *repository) GetClinicsByEps(epsID string, page int, pageSize int) ([]Clinic, error) {
	var clinics []Clinic

	// Calculate pagination offset
	offset := (page - 1) * pageSize

	// Sub-query that returns the clinic IDs associated to the EPS
	subQuery := r.db.
		Table("clinic_eps").
		Select("clinic_information_clinic_id").
		Where("eps_id = ?", epsID)

	// Main query fetching the clinics together with their information and relations
	err := r.db.
		Preload("ClinicInformation.ServicesOffered").
		Preload("ClinicInformation.EpsOffered").
		Preload("ClinicInformation.Photos").
		Where("id IN (?)", subQuery).
		Limit(pageSize).
		Offset(offset).
		Find(&clinics).Error

	if err != nil {
		return nil, err
	}

	return clinics, nil
}
