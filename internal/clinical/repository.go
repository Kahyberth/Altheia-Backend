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
	GetClinicByID(clinicID string) (*ClinicCompleteInfoResponse, error)
	AssignServicesToClinic(dto AssignServicesClinicDTO) error
	GetClinicsByEps(epsID string, page int, pageSize int) ([]Clinic, error)
	GetClinicPersonnel(clinicID string) ([]users.User, error)

	GetMedicalHistoryByPatientID(patientID string) (*MedicalHistoryResponseDTO, error)
	CreateMedicalHistory(dto CreateMedicalHistoryDTO) error
	CreateConsultation(dto CreateConsultationDTO) error
	GetOrCreateMedicalHistory(patientID string) (*MedicalHistory, error)
	UpdateMedicalHistory(historyID string, dto UpdateMedicalHistoryDTO) error
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

func (r *repository) GetClinicByID(clinicID string) (*ClinicCompleteInfoResponse, error) {
	var clinic Clinic
	var owner users.User

	if err := r.db.
		Preload("ClinicInformation.ServicesOffered").
		Preload("ClinicInformation.EpsOffered").
		Preload("ClinicInformation.Photos").
		Preload("Physicians.User").
		Preload("Receptionists.User").
		Preload("LabTechnicians.User").
		Preload("Patients.User").
		Where("id = ?", clinicID).
		First(&clinic).Error; err != nil {
		return nil, fmt.Errorf("clinic not found: %v", err)
	}

	// Get the owner information
	var clinicOwner users.ClinicOwner
	if err := r.db.Where("clinic_id = ?", clinicID).First(&clinicOwner).Error; err != nil {
		return nil, fmt.Errorf("clinic owner not found: %v", err)
	}

	if err := r.db.Where("id = ?", clinicOwner.UserID).First(&owner).Error; err != nil {
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

func (r *repository) GetClinicsByEps(epsID string, page int, pageSize int) ([]Clinic, error) {
	var clinics []Clinic

	offset := (page - 1) * pageSize

	subQuery := r.db.
		Table("clinic_eps").
		Select("clinic_information_clinic_id").
		Where("eps_id = ?", epsID)

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

func (r *repository) GetClinicPersonnel(clinicID string) ([]users.User, error) {
	var clinic Clinic

	if err := r.db.
		Preload("Physicians.User").
		Preload("Physicians").
		Preload("Receptionists.User").
		Preload("Receptionists").
		Preload("Patients.User").
		Preload("Patients").
		Preload("LabTechnicians.User").
		Preload("LabTechnicians").
		Where("id = ?", clinicID).
		First(&clinic).Error; err != nil {
		return nil, err
	}

	var personnel []users.User

	for _, p := range clinic.Physicians {
		if p.User != nil {
			user := *p.User
			user.Physician = p
			personnel = append(personnel, user)
		}
	}

	for _, rcp := range clinic.Receptionists {
		if rcp.User != nil {
			user := *rcp.User
			user.Receptionist = rcp
			personnel = append(personnel, user)
		}
	}

	for _, patient := range clinic.Patients {
		if patient.User != nil {
			user := *patient.User
			user.Patient = patient
			personnel = append(personnel, user)
		}
	}

	for _, lab := range clinic.LabTechnicians {
		if lab.User != nil {
			user := *lab.User
			user.LabTechnician = lab
			personnel = append(personnel, user)
		}
	}

	return personnel, nil
}

func (r *repository) GetMedicalHistoryByPatientID(patientID string) (*MedicalHistoryResponseDTO, error) {
	var medicalHistory MedicalHistory
	var patient users.Patient

	if err := r.db.Preload("User").Where("id = ?", patientID).First(&patient).Error; err != nil {
		return nil, fmt.Errorf("patient not found: %v", err)
	}

	err := r.db.
		Preload("Consultations.Physician.User").
		Preload("Consultations.Prescriptions").
		Where("patient_id = ?", patientID).
		First(&medicalHistory).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &MedicalHistoryResponseDTO{
				ID:                 patientID + "_empty",
				PatientId:          patientID,
				PatientName:        patient.User.Name,
				ConsultReason:      "",
				PersonalInfo:       "",
				FamilyInfo:         "",
				Allergies:          "",
				Observations:       "",
				LastUpdate:         time.Now(),
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
				Consultations:      []EnhancedConsultationResponseDTO{},
				TotalConsultations: 0,
			}, nil
		}
		return nil, fmt.Errorf("error fetching medical history: %v", err)
	}

	response := &MedicalHistoryResponseDTO{
		ID:                 medicalHistory.ID,
		PatientId:          medicalHistory.PatientId,
		PatientName:        patient.User.Name,
		ConsultReason:      medicalHistory.ConsultReason,
		PersonalInfo:       medicalHistory.PersonalInfo,
		FamilyInfo:         medicalHistory.FamilyInfo,
		Allergies:          medicalHistory.Allergies,
		Observations:       medicalHistory.Observations,
		LastUpdate:         medicalHistory.LastUpdate,
		CreatedAt:          medicalHistory.CreatedAt,
		UpdatedAt:          medicalHistory.UpdatedAt,
		Consultations:      []EnhancedConsultationResponseDTO{},
		TotalConsultations: len(medicalHistory.Consultations),
	}

	for _, consultation := range medicalHistory.Consultations {
		consultationDTO := EnhancedConsultationResponseDTO{
			ID:               consultation.ID,
			MedicalHistoryId: consultation.MedicalHistoryId,
			Symptoms:         consultation.Symptoms,
			Diagnosis:        consultation.Diagnosis,
			Treatment:        consultation.Treatment,
			Notes:            consultation.Notes,
			PhysicianInfo:    PhysicianInfoDTO{},
			Metadata: ConsultationMetadata{
				CreatedAt:   consultation.CreatedAt,
				UpdatedAt:   consultation.UpdatedAt,
				ConsultDate: consultation.ConsultDate,
			},
			Prescriptions: []PrescriptionResponseDTO{},
		}

		if consultation.Physician.User != nil {
			consultationDTO.PhysicianInfo = PhysicianInfoDTO{
				ID:                 consultation.Physician.User.ID,
				Name:               consultation.Physician.User.Name,
				Email:              consultation.Physician.User.Email,
				Phone:              consultation.Physician.User.Phone,
				DocumentNumber:     consultation.Physician.User.DocumentNumber,
				PhysicianSpecialty: consultation.Physician.PhysicianSpecialty,
				LicenseNumber:      consultation.Physician.LicenseNumber,
				Gender:             consultation.Physician.User.Gender,
			}
		}

		for _, prescription := range consultation.Prescriptions {
			consultationDTO.Prescriptions = append(consultationDTO.Prescriptions, PrescriptionResponseDTO{
				ID:             prescription.ID,
				ConsultationId: prescription.ConsultationId,
				Medicine:       prescription.Medicine,
				Dosage:         prescription.Dosage,
				Frequency:      prescription.Frequency,
				Duration:       prescription.Duration,
				Instructions:   prescription.Instructions,
				IssuedAt:       prescription.IssuedAt,
				CreatedAt:      prescription.CreatedAt,
				UpdatedAt:      prescription.UpdatedAt,
			})
		}

		response.Consultations = append(response.Consultations, consultationDTO)
	}

	return response, nil
}

func (r *repository) CreateMedicalHistory(dto CreateMedicalHistoryDTO) error {

	var patient users.Patient
	if err := r.db.Where("id = ?", dto.PatientId).First(&patient).Error; err != nil {
		return fmt.Errorf("patient not found: %v", err)
	}

	var existingHistory MedicalHistory
	if err := r.db.Where("patient_id = ?", dto.PatientId).First(&existingHistory).Error; err == nil {
		return fmt.Errorf("medical history already exists for this patient")
	}

	nanoid, _ := gonanoid.Nanoid()
	medicalHistory := MedicalHistory{
		ID:            nanoid,
		PatientId:     dto.PatientId,
		ConsultReason: dto.ConsultReason,
		PersonalInfo:  dto.PersonalInfo,
		FamilyInfo:    dto.FamilyInfo,
		Allergies:     dto.Allergies,
		Observations:  dto.Observations,
		LastUpdate:    time.Now(),
	}

	return r.db.Create(&medicalHistory).Error
}

func (r *repository) CreateConsultation(dto CreateConsultationDTO) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		var patient users.Patient
		if err := tx.Where("id = ?", dto.PatientId).First(&patient).Error; err != nil {
			return fmt.Errorf("patient not found: %v", err)
		}

		var physician users.Physician
		if err := tx.Where("id = ?", dto.PhysicianId).First(&physician).Error; err != nil {
			return fmt.Errorf("physician not found: %v", err)
		}

		medicalHistory, err := r.getOrCreateMedicalHistoryTx(tx, dto.PatientId)
		if err != nil {
			return err
		}

		if dto.UpdateMedicalHistory {
			updates := map[string]interface{}{}
			if dto.ConsultReason != "" {
				updates["consult_reason"] = dto.ConsultReason
			}
			if dto.PersonalInfo != "" {
				updates["personal_info"] = dto.PersonalInfo
			}
			if dto.FamilyInfo != "" {
				updates["family_info"] = dto.FamilyInfo
			}
			if dto.Allergies != "" {
				updates["allergies"] = dto.Allergies
			}
			if dto.Observations != "" {
				updates["observations"] = dto.Observations
			}

			if len(updates) > 0 {
				updates["last_update"] = time.Now()
				if err := tx.Model(&medicalHistory).Updates(updates).Error; err != nil {
					return fmt.Errorf("error updating medical history: %v", err)
				}
			}
		}

		consultationID, _ := gonanoid.Nanoid()
		consultation := MedicalConsultation{
			ID:               consultationID,
			MedicalHistoryId: medicalHistory.ID,
			PhysicianId:      dto.PhysicianId,
			ConsultDate:      time.Now(),
			Symptoms:         dto.Symptoms,
			Diagnosis:        dto.Diagnosis,
			Treatment:        dto.Treatment,
			Notes:            dto.Notes,
		}

		if err := tx.Create(&consultation).Error; err != nil {
			return fmt.Errorf("error creating consultation: %v", err)
		}

		for _, prescDto := range dto.Prescriptions {
			prescriptionID, _ := gonanoid.Nanoid()
			prescription := MedicalPrescription{
				ID:             prescriptionID,
				ConsultationId: consultationID,
				Medicine:       prescDto.Medicine,
				Dosage:         prescDto.Dosage,
				Frequency:      prescDto.Frequency,
				Duration:       prescDto.Duration,
				Instructions:   prescDto.Instructions,
				IssuedAt:       time.Now(),
			}

			if err := tx.Create(&prescription).Error; err != nil {
				return fmt.Errorf("error creating prescription: %v", err)
			}
		}

		if err := tx.Model(&medicalHistory).Update("last_update", time.Now()).Error; err != nil {
			return fmt.Errorf("error updating medical history: %v", err)
		}

		return nil
	})
}

func (r *repository) GetOrCreateMedicalHistory(patientID string) (*MedicalHistory, error) {
	return r.getOrCreateMedicalHistoryTx(r.db, patientID)
}

func (r *repository) getOrCreateMedicalHistoryTx(tx *gorm.DB, patientID string) (*MedicalHistory, error) {
	var medicalHistory MedicalHistory

	err := tx.Where("patient_id = ?", patientID).First(&medicalHistory).Error
	if err == nil {
		return &medicalHistory, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("error searching medical history: %v", err)
	}

	nanoid, _ := gonanoid.Nanoid()
	medicalHistory = MedicalHistory{
		ID:         nanoid,
		PatientId:  patientID,
		LastUpdate: time.Now(),
	}

	if err := tx.Create(&medicalHistory).Error; err != nil {
		return nil, fmt.Errorf("error creating medical history: %v", err)
	}

	return &medicalHistory, nil
}

func (r *repository) UpdateMedicalHistory(historyID string, dto UpdateMedicalHistoryDTO) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		var medicalHistory MedicalHistory
		if err := tx.Where("id = ?", historyID).First(&medicalHistory).Error; err != nil {
			return fmt.Errorf("medical history not found: %v", err)
		}

		updates := map[string]interface{}{
			"last_update": time.Now(),
		}

		updates["consult_reason"] = dto.ConsultReason
		updates["personal_info"] = dto.PersonalInfo
		updates["family_info"] = dto.FamilyInfo
		updates["allergies"] = dto.Allergies
		updates["observations"] = dto.Observations

		if err := tx.Model(&medicalHistory).Updates(updates).Error; err != nil {
			return fmt.Errorf("error updating medical history: %v", err)
		}

		return nil
	})
}
