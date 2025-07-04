package clinical

type Service interface {
	CreateClinical(createClinicDto CreateClinicDTO) error
	CreateEps(epsDto CreateEpsDto) error
	GetAllEps(page int, pagSize int) ([]EPS, error)
	CreateServicesOffered(servicesOffered CreateServicesDto) error
	GetAllServicesOffered(page int, pagSize int) ([]ServicesOffered, error)
	GetClinicByOwnerID(ownerID string) (*ClinicCompleteInfoResponse, error)
	GetClinicByID(clinicID string) (*ClinicCompleteInfoResponse, error)
	AssignServicesToClinic(dto AssignServicesClinicDTO) error
	GetClinicsByEps(epsID string, page int, pageSize int) ([]Clinic, error)
	GetClinicPersonnel(clinicID string) (ClinicPersonnelResponse, error)

	GetMedicalHistoryByPatientID(patientID string) (*MedicalHistoryResponseDTO, error)
	GetMedicalHistoryComprehensive(patientID string) (*ComprehensiveMedicalRecordsResponse, error)
	CreateMedicalHistory(dto CreateMedicalHistoryDTO) error
	CreateMedicalHistoryComprehensive(dto CreateMedicalHistoryDTO) (*ComprehensiveMedicalRecordsResponse, error)
	CreateConsultation(dto CreateConsultationDTO) error
	UpdateMedicalHistory(historyID string, dto UpdateMedicalHistoryDTO) error
	GetClinicMedicalHistoriesPaginated(clinicID string, page int, pageSize int) (*PaginatedMedicalHistoriesResponse, error)

	// Document methods
	AddDocumentsToMedicalHistory(dto AddDocumentsToMedicalHistoryDTO) (*AddDocumentsResponseDTO, error)
	AddDocumentsToConsultation(dto AddDocumentsToConsultationDTO) (*AddDocumentsResponseDTO, error)
	GetDocumentsByMedicalHistory(medicalHistoryId string) ([]DocumentResponseDTO, error)
	GetDocumentsByConsultation(consultationId string) ([]DocumentResponseDTO, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) CreateClinical(createClinicDto CreateClinicDTO) error {
	err := s.repo.CreateClinic(createClinicDto)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) CreateEps(epsDto CreateEpsDto) error {
	err := s.repo.CreateEps(epsDto)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) CreateServicesOffered(servicesOffered CreateServicesDto) error {
	err := s.repo.CreateServices(servicesOffered)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetAllServicesOffered(page int, pagSize int) ([]ServicesOffered, error) {
	var servicesOffered []ServicesOffered
	servicesOffered, serviceError := s.repo.GetAllServices(page, pagSize)
	if serviceError != nil {
		return nil, serviceError
	}
	return servicesOffered, nil
}

func (s *service) GetAllEps(page int, pagSize int) ([]EPS, error) {
	var eps []EPS
	eps, epsError := s.repo.GetAllEps(page, pagSize)
	if epsError != nil {
		return eps, epsError
	}
	return eps, nil
}

func (s *service) GetClinicByOwnerID(ownerID string) (*ClinicCompleteInfoResponse, error) {
	return s.repo.GetClinicByOwnerID(ownerID)
}

func (s *service) GetClinicByID(clinicID string) (*ClinicCompleteInfoResponse, error) {
	return s.repo.GetClinicByID(clinicID)
}

func (s *service) AssignServicesToClinic(dto AssignServicesClinicDTO) error {
	err := s.repo.AssignServicesToClinic(dto)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetClinicsByEps(epsID string, page int, pageSize int) ([]Clinic, error) {
	return s.repo.GetClinicsByEps(epsID, page, pageSize)
}

func (s *service) GetClinicPersonnel(clinicID string) (ClinicPersonnelResponse, error) {
	users, err := s.repo.GetClinicPersonnel(clinicID)
	if err != nil {
		return ClinicPersonnelResponse{}, err
	}

	var personnel []PersonnelResponse
	for _, user := range users {
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
		case "lab_technician":
			roleDetails = map[string]interface{}{
				"lab_technician_id": user.LabTechnician.ID,
				"clinic_id":         user.LabTechnician.ClinicID,
			}
		}

		personnelItem := PersonnelResponse{
			ID:             user.ID,
			Name:           user.Name,
			Email:          user.Email,
			Role:           user.Rol,
			Phone:          user.Phone,
			DocumentNumber: user.DocumentNumber,
			Status:         user.Status,
			Gender:         user.Gender,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
			LastLogin:      user.LastLogin,
			RoleDetails:    roleDetails,
		}

		personnel = append(personnel, personnelItem)
	}

	response := ClinicPersonnelResponse{
		Personnel: personnel,
		Count:     len(personnel),
	}

	return response, nil
}

func (s *service) GetMedicalHistoryByPatientID(patientID string) (*MedicalHistoryResponseDTO, error) {
	return s.repo.GetMedicalHistoryByPatientID(patientID)
}

func (s *service) GetMedicalHistoryComprehensive(patientID string) (*ComprehensiveMedicalRecordsResponse, error) {
	return s.repo.GetMedicalHistoryComprehensive(patientID)
}

func (s *service) CreateMedicalHistory(dto CreateMedicalHistoryDTO) error {
	err := s.repo.CreateMedicalHistory(dto)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) CreateMedicalHistoryComprehensive(dto CreateMedicalHistoryDTO) (*ComprehensiveMedicalRecordsResponse, error) {
	response, err := s.repo.CreateMedicalHistoryComprehensive(dto)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *service) CreateConsultation(dto CreateConsultationDTO) error {
	err := s.repo.CreateConsultation(dto)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateMedicalHistory(historyID string, dto UpdateMedicalHistoryDTO) error {
	err := s.repo.UpdateMedicalHistory(historyID, dto)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetClinicMedicalHistoriesPaginated(clinicID string, page int, pageSize int) (*PaginatedMedicalHistoriesResponse, error) {
	return s.repo.GetClinicMedicalHistoriesPaginated(clinicID, page, pageSize)
}

func (s *service) AddDocumentsToMedicalHistory(dto AddDocumentsToMedicalHistoryDTO) (*AddDocumentsResponseDTO, error) {
	return s.repo.AddDocumentsToMedicalHistory(dto)
}

func (s *service) AddDocumentsToConsultation(dto AddDocumentsToConsultationDTO) (*AddDocumentsResponseDTO, error) {
	return s.repo.AddDocumentsToConsultation(dto)
}

func (s *service) GetDocumentsByMedicalHistory(medicalHistoryId string) ([]DocumentResponseDTO, error) {
	return s.repo.GetDocumentsByMedicalHistory(medicalHistoryId)
}

func (s *service) GetDocumentsByConsultation(consultationId string) ([]DocumentResponseDTO, error) {
	return s.repo.GetDocumentsByConsultation(consultationId)
}
