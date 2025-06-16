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
	GetMedicalHistoryComprehensive(patientID string) (*ComprehensiveMedicalRecordsResponse, error)
	CreateMedicalHistory(dto CreateMedicalHistoryDTO) error
	CreateMedicalHistoryComprehensive(dto CreateMedicalHistoryDTO) (*ComprehensiveMedicalRecordsResponse, error)
	CreateConsultation(dto CreateConsultationDTO) error
	GetOrCreateMedicalHistory(patientID string) (*MedicalHistory, error)
	UpdateMedicalHistory(historyID string, dto UpdateMedicalHistoryDTO) error
	GetClinicMedicalHistoriesPaginated(clinicID string, page int, pageSize int) (*PaginatedMedicalHistoriesResponse, error)

	// Document methods
	AddDocumentsToMedicalHistory(dto AddDocumentsToMedicalHistoryDTO) (*AddDocumentsResponseDTO, error)
	AddDocumentsToConsultation(dto AddDocumentsToConsultationDTO) (*AddDocumentsResponseDTO, error)
	GetDocumentsByMedicalHistory(medicalHistoryId string) ([]DocumentResponseDTO, error)
	GetDocumentsByConsultation(consultationId string) ([]DocumentResponseDTO, error)
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

func (r *repository) GetMedicalHistoryComprehensive(patientID string) (*ComprehensiveMedicalRecordsResponse, error) {

	var patient users.Patient
	if err := r.db.Preload("User").Where("id = ?", patientID).First(&patient).Error; err != nil {
		return nil, fmt.Errorf("patient not found: %v", err)
	}

	response := r.buildComprehensiveResponseWithRealData(patientID)
	return response, nil
}

func (r *repository) CreateMedicalHistory(dto CreateMedicalHistoryDTO) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var patient users.Patient
		if err := tx.Where("id = ?", dto.PatientId).First(&patient).Error; err != nil {
			return fmt.Errorf("patient not found: %v", err)
		}

		var existingHistory MedicalHistory
		if err := tx.Where("patient_id = ?", dto.PatientId).First(&existingHistory).Error; err == nil {
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

		if err := tx.Create(&medicalHistory).Error; err != nil {
			return fmt.Errorf("error creating medical history: %v", err)
		}

		if len(dto.Prescriptions) > 0 || dto.PhysicianId != "" {
			if len(dto.Prescriptions) > 0 && dto.PhysicianId == "" {
				return fmt.Errorf("physician_id is required when prescriptions are provided")
			}

			if dto.PhysicianId != "" {
				var physician users.Physician
				if err := tx.Where("id = ?", dto.PhysicianId).First(&physician).Error; err != nil {
					return fmt.Errorf("physician not found: %v", err)
				}
			}

			physicianId := dto.PhysicianId
			if physicianId == "" {
				physicianId = "SYSTEM"
			}

			consultationID, _ := gonanoid.Nanoid()
			consultation := MedicalConsultation{
				ID:               consultationID,
				MedicalHistoryId: medicalHistory.ID,
				PhysicianId:      physicianId,
				ConsultDate:      time.Now(),
				Symptoms:         "",
				Diagnosis:        "",
				Treatment:        "Initial medical history",
				Notes:            fmt.Sprintf("Medical history created with %d prescriptions", len(dto.Prescriptions)),
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
		}

		if len(dto.Documents) > 0 {
			uploadedBy := "SYSTEM"
			for _, doc := range dto.Documents {
				_, err := r.saveDocument(tx, doc, &medicalHistory.ID, nil, uploadedBy)
				if err != nil {
					return fmt.Errorf("error saving document %s: %v", doc.Name, err)
				}
			}
		}

		return nil
	})
}

func (r *repository) CreateMedicalHistoryComprehensive(dto CreateMedicalHistoryDTO) (*ComprehensiveMedicalRecordsResponse, error) {
	err := r.CreateMedicalHistory(dto)
	if err != nil {
		return nil, err
	}

	response := r.buildComprehensiveResponseWithRealData(dto.PatientId)
	return response, nil
}

func (r *repository) buildComprehensiveResponseWithRealData(patientID string) *ComprehensiveMedicalRecordsResponse {
	patients := r.getRealPatients(patientID)

	medicalRecords := r.getRealMedicalRecords(patientID)

	metadata := r.getRealMetadata(patientID)

	return &ComprehensiveMedicalRecordsResponse{
		Success: true,
		Data: MedicalRecordsData{
			Patients:       patients,
			MedicalRecords: medicalRecords,
			Metadata:       metadata,
		},
	}
}

func (r *repository) getRealPatients(patientID string) []PatientBasicInfo {
	var patient users.Patient
	if err := r.db.Preload("User").Where("id = ?", patientID).First(&patient).Error; err != nil {
		return []PatientBasicInfo{}
	}

	age := 0
	if patient.DateOfBirth != "" {
		if dob, err := time.Parse("2006-01-02", patient.DateOfBirth); err == nil {
			age = int(time.Since(dob).Hours() / 24 / 365)
		}
	}

	mrn := "MRN-" + patient.ID
	if len(patient.ID) > 5 {
		mrn = "MRN-" + patient.ID[len(patient.ID)-5:]
	}

	avatarText := "P"
	if patient.User != nil && len(patient.User.Name) > 0 {
		avatarText = string(patient.User.Name[0])
	}

	return []PatientBasicInfo{
		{
			ID:     patient.ID,
			Name:   patient.User.Name,
			Age:    age,
			Gender: patient.User.Gender,
			DOB:    patient.DateOfBirth,
			MRN:    mrn,
			Avatar: "/placeholder.svg?height=128&width=128&text=" + avatarText,
		},
	}
}

func (r *repository) getRealMedicalRecords(patientID string) []MedicalRecord {
	var medicalHistory MedicalHistory
	var records []MedicalRecord

	err := r.db.
		Preload("Consultations.Physician.User").
		Preload("Consultations.Prescriptions").
		Where("patient_id = ?", patientID).
		First(&medicalHistory).Error

	if err != nil {
		return []MedicalRecord{}
	}

	if medicalHistory.ConsultReason != "" || medicalHistory.Observations != "" || medicalHistory.PersonalInfo != "" || medicalHistory.FamilyInfo != "" {
		historyRecord := MedicalRecord{
			ID:        "REC-HISTORY-" + medicalHistory.ID,
			PatientID: patientID,
			Type:      "diagnoses",
			Title:     "Historia Clínica - " + medicalHistory.CreatedAt.Format("2006-01-02"),
			Date:      medicalHistory.CreatedAt.Format("2006-01-02"),
			Provider:  "Sistema",
			Status:    "active",
			Content: MedicalRecordContent{
				"consult_reason": medicalHistory.ConsultReason,
				"description":    medicalHistory.ConsultReason,
				"personal_info":  medicalHistory.PersonalInfo,
				"family_info":    medicalHistory.FamilyInfo,
				"allergies":      medicalHistory.Allergies,
				"observations":   medicalHistory.Observations,
				"notes":          fmt.Sprintf("Historia clínica creada. Motivo: %s. Observaciones: %s", medicalHistory.ConsultReason, medicalHistory.Observations),
			},
			Documents: []Document{},
		}
		records = append(records, historyRecord)
	}

	for _, consultation := range medicalHistory.Consultations {

		if len(consultation.Prescriptions) > 0 {
			medications := make([]map[string]interface{}, 0)
			for _, prescription := range consultation.Prescriptions {
				medications = append(medications, map[string]interface{}{
					"name":         prescription.Medicine,
					"dosage":       prescription.Dosage,
					"frequency":    prescription.Frequency,
					"duration":     prescription.Duration,
					"startDate":    prescription.IssuedAt.Format("2006-01-02"),
					"prescriber":   consultation.Physician.User.Name,
					"instructions": prescription.Instructions,
				})
			}

			record := MedicalRecord{
				ID:        "REC-" + consultation.ID,
				PatientID: patientID,
				Type:      "medications",
				Title:     "Prescripciones - " + consultation.ConsultDate.Format("2006-01-02"),
				Date:      consultation.ConsultDate.Format("2006-01-02"),
				Provider:  consultation.Physician.User.Name,
				Status:    "active",
				Content: MedicalRecordContent{
					"medications":    medications,
					"notes":          consultation.Notes,
					"description":    fmt.Sprintf("Prescripciones médicas. Tratamiento: %s", consultation.Treatment),
					"observations":   consultation.Notes,
					"consult_reason": consultation.Treatment,
				},
				Documents: []Document{},
			}
			records = append(records, record)
		}

		if consultation.Diagnosis != "" || consultation.Symptoms != "" || consultation.Treatment != "" {
			record := MedicalRecord{
				ID:        "REC-DIAG-" + consultation.ID,
				PatientID: patientID,
				Type:      "diagnoses",
				Title:     "Consulta Médica - " + consultation.ConsultDate.Format("2006-01-02"),
				Date:      consultation.ConsultDate.Format("2006-01-02"),
				Provider:  consultation.Physician.User.Name,
				Status:    "completed",
				Content: MedicalRecordContent{
					"symptoms":       consultation.Symptoms,
					"diagnosis":      consultation.Diagnosis,
					"treatment":      consultation.Treatment,
					"notes":          consultation.Notes,
					"description":    fmt.Sprintf("Consulta médica. Síntomas: %s. Diagnóstico: %s", consultation.Symptoms, consultation.Diagnosis),
					"observations":   consultation.Notes,
					"consult_reason": consultation.Treatment,
				},
				Documents: []Document{},
			}
			records = append(records, record)
		}
	}

	return records
}

func (r *repository) getRealMetadata(patientID string) Metadata {
	var totalRecords int64
	var totalPatients int64

	r.db.Model(&MedicalConsultation{}).
		Joins("JOIN medical_histories ON medical_consultations.medical_history_id = medical_histories.id").
		Where("medical_histories.patient_id = ?", patientID).
		Count(&totalRecords)

	r.db.Model(&users.Patient{}).Count(&totalPatients)

	return Metadata{
		TotalRecords:  int(totalRecords),
		TotalPatients: int(totalPatients),
		LastUpdated:   time.Now().Format("2006-01-02T15:04:05"),
		Version:       "1.0",
	}
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

		if len(dto.Documents) > 0 {
			uploadedBy := dto.PhysicianId
			for _, doc := range dto.Documents {
				_, err := r.saveDocument(tx, doc, &medicalHistory.ID, &consultationID, uploadedBy)
				if err != nil {
					return fmt.Errorf("error saving consultation document %s: %v", doc.Name, err)
				}
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

func (r *repository) GetClinicMedicalHistoriesPaginated(clinicID string, page int, pageSize int) (*PaginatedMedicalHistoriesResponse, error) {

	var patients []users.Patient
	var totalPatients int64

	err := r.db.Model(&users.Patient{}).
		Where("clinic_id = ?", clinicID).
		Count(&totalPatients).Error
	if err != nil {
		return nil, fmt.Errorf("error counting patients in clinic: %v", err)
	}

	offset := (page - 1) * pageSize
	err = r.db.
		Preload("User").
		Where("clinic_id = ?", clinicID).
		Limit(pageSize).
		Offset(offset).
		Find(&patients).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching patients from clinic: %v", err)
	}

	var comprehensiveRecords []ComprehensiveMedicalRecord
	totalMedicalRecords := 0

	for _, patient := range patients {
		if patient.User == nil {
			continue
		}

		age := 0
		if patient.DateOfBirth != "" {
			if dob, err := time.Parse("2006-01-02", patient.DateOfBirth); err == nil {
				age = int(time.Since(dob).Hours() / 24 / 365)
			}
		}

		mrn := "MRN-" + patient.ID
		if len(patient.ID) > 5 {
			mrn = "MRN-" + patient.ID[len(patient.ID)-5:]
		}

		avatarText := "P"
		if patient.User != nil && len(patient.User.Name) > 0 {
			avatarText = string(patient.User.Name[0])
		}

		patientInfo := PatientBasicInfo{
			ID:     patient.ID,
			Name:   patient.User.Name,
			Age:    age,
			Gender: patient.User.Gender,
			DOB:    patient.DateOfBirth,
			MRN:    mrn,
			Avatar: "/placeholder.svg?height=128&width=128&text=" + avatarText,
		}

		medicalRecords := r.getRealMedicalRecords(patient.ID)

		lastUpdate := time.Now().Format("2006-01-02")
		var medicalHistory MedicalHistory
		if err := r.db.Where("patient_id = ?", patient.ID).Order("updated_at DESC").First(&medicalHistory).Error; err == nil {
			lastUpdate = medicalHistory.UpdatedAt.Format("2006-01-02")
		}

		comprehensiveRecord := ComprehensiveMedicalRecord{
			Patient:        patientInfo,
			MedicalRecords: medicalRecords,
			LastUpdate:     lastUpdate,
			RecordCount:    len(medicalRecords),
		}

		comprehensiveRecords = append(comprehensiveRecords, comprehensiveRecord)
		totalMedicalRecords += len(medicalRecords)
	}

	totalPages := int((totalPatients + int64(pageSize) - 1) / int64(pageSize))

	paginationInfo := PaginationInfo{
		CurrentPage:  page,
		PageSize:     pageSize,
		TotalPages:   totalPages,
		TotalRecords: totalPatients,
		HasNext:      page < totalPages,
		HasPrevious:  page > 1,
	}

	var mostActivePatient string
	if len(comprehensiveRecords) > 0 {
		maxRecords := 0
		for _, record := range comprehensiveRecords {
			if record.RecordCount > maxRecords {
				maxRecords = record.RecordCount
				mostActivePatient = record.Patient.Name
			}
		}
	}

	summary := ClinicMedicalSummary{
		TotalPatients:       int(totalPatients),
		TotalMedicalRecords: totalMedicalRecords,
		RecentActivity:      time.Now().Format("2006-01-02"),
		MostActivePatient:   mostActivePatient,
		LastUpdated:         time.Now().Format("2006-01-02 15:04:05"),
	}

	response := &PaginatedMedicalHistoriesResponse{
		Success:    true,
		Data:       comprehensiveRecords,
		Pagination: paginationInfo,
		Summary:    summary,
	}

	return response, nil
}

func (r *repository) saveDocument(tx *gorm.DB, doc CreateDocumentDTO, medicalHistoryId *string, consultationId *string, uploadedBy string) (*MedicalDocument, error) {
	docID, _ := gonanoid.Nanoid()

	var filePath, fileURL string
	var fileSize int64
	var mimeType string

	if doc.Base64Data != "" {
		savedPath, savedURL, size, mime, err := r.saveBase64File(doc.Base64Data, doc.Name, doc.Type)
		if err != nil {
			return nil, fmt.Errorf("error saving base64 file: %v", err)
		}
		filePath = savedPath
		fileURL = savedURL
		fileSize = size
		mimeType = mime
	} else if doc.URL != "" {

		fileURL = doc.URL
		filePath = ""
		mimeType = doc.Type

		if doc.Size != "" {
			fileSize = 0
		}
	} else {
		return nil, fmt.Errorf("either base64_data or url must be provided")
	}

	medicalDoc := MedicalDocument{
		ID:               docID,
		MedicalHistoryId: medicalHistoryId,
		ConsultationId:   consultationId,
		Name:             r.generateSafeFileName(doc.Name),
		OriginalName:     doc.Name,
		Type:             doc.Type,
		Size:             fileSize,
		MimeType:         mimeType,
		FilePath:         filePath,
		URL:              fileURL,
		Description:      doc.Description,
		UploadedBy:       uploadedBy,
		UploadedAt:       time.Now(),
		IsPublic:         false,
	}

	if err := tx.Create(&medicalDoc).Error; err != nil {
		return nil, fmt.Errorf("error saving document metadata: %v", err)
	}

	return &medicalDoc, nil
}

func (r *repository) saveBase64File(base64Data, fileName, fileType string) (filePath, fileURL string, fileSize int64, mimeType string, err error) {

	if len(base64Data) == 0 {
		return "", "", 0, "", fmt.Errorf("empty base64 data")
	}

	uploadsDir := "./uploads/medical-documents"
	if err := r.ensureDirectoryExists(uploadsDir); err != nil {
		return "", "", 0, "", err
	}

	docID, _ := gonanoid.Nanoid()
	safeFileName := r.generateSafeFileName(fileName)
	fullPath := fmt.Sprintf("%s/%s_%s", uploadsDir, docID, safeFileName)

	fileSize = int64(len(base64Data) * 3 / 4)
	mimeType = r.getMimeTypeFromExtension(fileType)
	filePath = fullPath
	fileURL = fmt.Sprintf("/api/documents/%s_%s", docID, safeFileName)

	fmt.Printf("Would save document: %s (Size: %d bytes, Type: %s)\n", fullPath, fileSize, mimeType)

	return filePath, fileURL, fileSize, mimeType, nil
}

func (r *repository) ensureDirectoryExists(dir string) error {

	fmt.Printf("Ensuring directory exists: %s\n", dir)
	return nil
}

func (r *repository) generateSafeFileName(originalName string) string {
	return fmt.Sprintf("%d_%s", time.Now().Unix(), originalName)
}

func (r *repository) getMimeTypeFromExtension(fileType string) string {
	mimeTypes := map[string]string{
		"pdf":  "application/pdf",
		"jpg":  "image/jpeg",
		"jpeg": "image/jpeg",
		"png":  "image/png",
		"gif":  "image/gif",
		"doc":  "application/msword",
		"docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"txt":  "text/plain",
	}

	if mime, exists := mimeTypes[fileType]; exists {
		return mime
	}
	return "application/octet-stream"
}

func (r *repository) AddDocumentsToMedicalHistory(dto AddDocumentsToMedicalHistoryDTO) (*AddDocumentsResponseDTO, error) {
	var savedDocuments []DocumentResponseDTO

	err := r.db.Transaction(func(tx *gorm.DB) error {

		var medicalHistory MedicalHistory
		if err := tx.Where("id = ?", dto.MedicalHistoryId).First(&medicalHistory).Error; err != nil {
			return fmt.Errorf("medical history not found: %v", err)
		}

		for _, doc := range dto.Documents {
			uploadedBy := dto.UploadedBy
			if uploadedBy == "" {
				uploadedBy = "SYSTEM"
			}

			savedDoc, err := r.saveDocument(tx, doc, &dto.MedicalHistoryId, nil, uploadedBy)
			if err != nil {
				return fmt.Errorf("error saving document %s: %v", doc.Name, err)
			}

			savedDocuments = append(savedDocuments, DocumentResponseDTO{
				ID:           savedDoc.ID,
				Name:         savedDoc.Name,
				OriginalName: savedDoc.OriginalName,
				Type:         savedDoc.Type,
				Size:         savedDoc.Size,
				MimeType:     savedDoc.MimeType,
				URL:          savedDoc.URL,
				Description:  savedDoc.Description,
				UploadedBy:   savedDoc.UploadedBy,
				UploadedAt:   savedDoc.UploadedAt.Format("2006-01-02T15:04:05"),
				IsPublic:     savedDoc.IsPublic,
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &AddDocumentsResponseDTO{
		Success:   true,
		Message:   fmt.Sprintf("Successfully uploaded %d document(s)", len(savedDocuments)),
		Documents: savedDocuments,
	}, nil
}

func (r *repository) AddDocumentsToConsultation(dto AddDocumentsToConsultationDTO) (*AddDocumentsResponseDTO, error) {
	var savedDocuments []DocumentResponseDTO

	err := r.db.Transaction(func(tx *gorm.DB) error {

		var consultation MedicalConsultation
		if err := tx.Where("id = ?", dto.ConsultationId).First(&consultation).Error; err != nil {
			return fmt.Errorf("consultation not found: %v", err)
		}

		for _, doc := range dto.Documents {
			uploadedBy := dto.UploadedBy
			if uploadedBy == "" {
				uploadedBy = consultation.PhysicianId
			}

			savedDoc, err := r.saveDocument(tx, doc, &consultation.MedicalHistoryId, &dto.ConsultationId, uploadedBy)
			if err != nil {
				return fmt.Errorf("error saving document %s: %v", doc.Name, err)
			}

			savedDocuments = append(savedDocuments, DocumentResponseDTO{
				ID:           savedDoc.ID,
				Name:         savedDoc.Name,
				OriginalName: savedDoc.OriginalName,
				Type:         savedDoc.Type,
				Size:         savedDoc.Size,
				MimeType:     savedDoc.MimeType,
				URL:          savedDoc.URL,
				Description:  savedDoc.Description,
				UploadedBy:   savedDoc.UploadedBy,
				UploadedAt:   savedDoc.UploadedAt.Format("2006-01-02T15:04:05"),
				IsPublic:     savedDoc.IsPublic,
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &AddDocumentsResponseDTO{
		Success:   true,
		Message:   fmt.Sprintf("Successfully uploaded %d document(s) to consultation", len(savedDocuments)),
		Documents: savedDocuments,
	}, nil
}

func (r *repository) GetDocumentsByMedicalHistory(medicalHistoryId string) ([]DocumentResponseDTO, error) {
	var documents []MedicalDocument

	err := r.db.Where("medical_history_id = ?", medicalHistoryId).Find(&documents).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching documents: %v", err)
	}

	var documentDTOs []DocumentResponseDTO
	for _, doc := range documents {
		documentDTOs = append(documentDTOs, DocumentResponseDTO{
			ID:           doc.ID,
			Name:         doc.Name,
			OriginalName: doc.OriginalName,
			Type:         doc.Type,
			Size:         doc.Size,
			MimeType:     doc.MimeType,
			URL:          doc.URL,
			Description:  doc.Description,
			UploadedBy:   doc.UploadedBy,
			UploadedAt:   doc.UploadedAt.Format("2006-01-02T15:04:05"),
			IsPublic:     doc.IsPublic,
		})
	}

	return documentDTOs, nil
}

func (r *repository) GetDocumentsByConsultation(consultationId string) ([]DocumentResponseDTO, error) {
	var documents []MedicalDocument

	err := r.db.Where("consultation_id = ?", consultationId).Find(&documents).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching documents: %v", err)
	}

	var documentDTOs []DocumentResponseDTO
	for _, doc := range documents {
		documentDTOs = append(documentDTOs, DocumentResponseDTO{
			ID:           doc.ID,
			Name:         doc.Name,
			OriginalName: doc.OriginalName,
			Type:         doc.Type,
			Size:         doc.Size,
			MimeType:     doc.MimeType,
			URL:          doc.URL,
			Description:  doc.Description,
			UploadedBy:   doc.UploadedBy,
			UploadedAt:   doc.UploadedAt.Format("2006-01-02T15:04:05"),
			IsPublic:     doc.IsPublic,
		})
	}

	return documentDTOs, nil
}
