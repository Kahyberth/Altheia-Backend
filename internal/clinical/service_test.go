package clinical

import (
	"errors"
	"testing"
)

// Mock Repository for testing
type mockClinicalRepository struct {
	clinics  []CreateClinicDTO
	eps      []EPS
	services []ServicesOffered
	err      error
}

func (m *mockClinicalRepository) CreateClinic(clinic CreateClinicDTO) error {
	if m.err != nil {
		return m.err
	}
	m.clinics = append(m.clinics, clinic)
	return nil
}

func (m *mockClinicalRepository) CreateEps(epsDto CreateEpsDto) error {
	if m.err != nil {
		return m.err
	}
	for _, name := range epsDto.Eps {
		m.eps = append(m.eps, EPS{ID: "eps-" + name, Name: name})
	}
	return nil
}

func (m *mockClinicalRepository) CreateServices(servicesDto CreateServicesDto) error {
	if m.err != nil {
		return m.err
	}
	for _, name := range servicesDto.ServicesOffered {
		m.services = append(m.services, ServicesOffered{ID: "service-" + name, Name: name})
	}
	return nil
}

func (m *mockClinicalRepository) GetAllEps(page int, pageSize int) ([]EPS, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.eps, nil
}

func (m *mockClinicalRepository) GetAllServices(page int, pageSize int) ([]ServicesOffered, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.services, nil
}

func setupClinicalTestService() (Service, *mockClinicalRepository) {
	mockRepo := &mockClinicalRepository{
		clinics:  []CreateClinicDTO{},
		eps:      []EPS{},
		services: []ServicesOffered{},
	}
	service := NewService(mockRepo)
	return service, mockRepo
}

func TestClinicalService_CreateClinical_Success(t *testing.T) {
	service, mockRepo := setupClinicalTestService()

	clinicDto := CreateClinicDTO{
		OwnerName:       "Dr. John Doe",
		OwnerEmail:      "john@example.com",
		OwnerPhone:      "1234567890",
		Name:            "Central Clinic",
		Email:           "clinic@example.com",
		Description:     "A comprehensive healthcare facility",
		Phone:           "0987654321",
		Website:         "https://centralclinic.com",
		Address:         "123 Main St",
		Country:         "Colombia",
		City:            "Bogotá",
		State:           "Cundinamarca",
		MemberCount:     50,
		ServicesOffered: []string{"Cardiology", "Neurology"},
		AcceptedEPS:     []string{"Sura", "Sanitas"},
	}

	err := service.CreateClinical(clinicDto)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(mockRepo.clinics) != 1 {
		t.Errorf("Expected 1 clinic to be created, got %d", len(mockRepo.clinics))
	}

	if mockRepo.clinics[0].Name != "Central Clinic" {
		t.Errorf("Expected clinic name 'Central Clinic', got %s", mockRepo.clinics[0].Name)
	}
}

func TestClinicalService_CreateClinical_Error(t *testing.T) {
	service, mockRepo := setupClinicalTestService()
	mockRepo.err = errors.New("database error")

	clinicDto := CreateClinicDTO{
		Name: "Test Clinic",
	}

	err := service.CreateClinical(clinicDto)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "database error" {
		t.Errorf("Expected 'database error', got %s", err.Error())
	}
}

func TestClinicalService_CreateEps_Success(t *testing.T) {
	service, mockRepo := setupClinicalTestService()

	epsDto := CreateEpsDto{
		Eps: []string{"Sura", "Sanitas", "Compensar"},
	}

	err := service.CreateEps(epsDto)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(mockRepo.eps) != 3 {
		t.Errorf("Expected 3 EPS to be created, got %d", len(mockRepo.eps))
	}

	// Check if all EPS were created
	expectedNames := []string{"Sura", "Sanitas", "Compensar"}
	for i, eps := range mockRepo.eps {
		if eps.Name != expectedNames[i] {
			t.Errorf("Expected EPS name '%s', got '%s'", expectedNames[i], eps.Name)
		}
	}
}

func TestClinicalService_CreateEps_Error(t *testing.T) {
	service, mockRepo := setupClinicalTestService()
	mockRepo.err = errors.New("database connection failed")

	epsDto := CreateEpsDto{
		Eps: []string{"Sura"},
	}

	err := service.CreateEps(epsDto)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "database connection failed" {
		t.Errorf("Expected 'database connection failed', got %s", err.Error())
	}
}

func TestClinicalService_CreateServicesOffered_Success(t *testing.T) {
	service, mockRepo := setupClinicalTestService()

	servicesDto := CreateServicesDto{
		ServicesOffered: []string{"Cardiology", "Neurology", "Pediatrics"},
	}

	err := service.CreateServicesOffered(servicesDto)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(mockRepo.services) != 3 {
		t.Errorf("Expected 3 services to be created, got %d", len(mockRepo.services))
	}

	expectedServices := []string{"Cardiology", "Neurology", "Pediatrics"}
	for i, service := range mockRepo.services {
		if service.Name != expectedServices[i] {
			t.Errorf("Expected service name '%s', got '%s'", expectedServices[i], service.Name)
		}
	}
}

func TestClinicalService_CreateServicesOffered_Error(t *testing.T) {
	service, mockRepo := setupClinicalTestService()
	mockRepo.err = errors.New("service creation failed")

	servicesDto := CreateServicesDto{
		ServicesOffered: []string{"Cardiology"},
	}

	err := service.CreateServicesOffered(servicesDto)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "service creation failed" {
		t.Errorf("Expected 'service creation failed', got %s", err.Error())
	}
}

func TestClinicalService_GetAllEps_Success(t *testing.T) {
	service, mockRepo := setupClinicalTestService()

	// Prepare mock data
	mockRepo.eps = []EPS{
		{ID: "eps-1", Name: "Sura"},
		{ID: "eps-2", Name: "Sanitas"},
		{ID: "eps-3", Name: "Compensar"},
	}

	eps, err := service.GetAllEps(1, 10)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(eps) != 3 {
		t.Errorf("Expected 3 EPS, got %d", len(eps))
	}

	if eps[0].Name != "Sura" {
		t.Errorf("Expected first EPS name 'Sura', got '%s'", eps[0].Name)
	}
}

func TestClinicalService_GetAllEps_Error(t *testing.T) {
	service, mockRepo := setupClinicalTestService()
	mockRepo.err = errors.New("failed to fetch EPS")

	eps, err := service.GetAllEps(1, 10)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "failed to fetch EPS" {
		t.Errorf("Expected 'failed to fetch EPS', got %s", err.Error())
	}

	if len(eps) != 0 {
		t.Errorf("Expected empty EPS list on error, got %d items", len(eps))
	}
}

func TestClinicalService_GetAllServicesOffered_Success(t *testing.T) {
	service, mockRepo := setupClinicalTestService()

	// Prepare mock data
	mockRepo.services = []ServicesOffered{
		{ID: "service-1", Name: "Cardiology"},
		{ID: "service-2", Name: "Neurology"},
		{ID: "service-3", Name: "Pediatrics"},
	}

	services, err := service.GetAllServicesOffered(1, 10)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(services) != 3 {
		t.Errorf("Expected 3 services, got %d", len(services))
	}

	if services[0].Name != "Cardiology" {
		t.Errorf("Expected first service name 'Cardiology', got '%s'", services[0].Name)
	}
}

func TestClinicalService_GetAllServicesOffered_Error(t *testing.T) {
	service, mockRepo := setupClinicalTestService()
	mockRepo.err = errors.New("failed to fetch services")

	services, err := service.GetAllServicesOffered(1, 10)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "failed to fetch services" {
		t.Errorf("Expected 'failed to fetch services', got %s", err.Error())
	}

	if services != nil {
		t.Errorf("Expected nil services list on error, got %v", services)
	}
}

func TestCreateClinicDTO_FieldValidation(t *testing.T) {
	dto := CreateClinicDTO{
		OwnerName:           "Dr. Jane Smith",
		OwnerEmail:          "jane@example.com",
		OwnerPhone:          "1234567890",
		OwnerPosition:       "Director",
		OwnerDocumentNumber: "12345678",
		OwnerGender:         "F",
		Name:                "Advanced Medical Center",
		Email:               "info@amc.com",
		Description:         "State-of-the-art medical facility",
		Phone:               "0987654321",
		Website:             "https://amc.com",
		Address:             "456 Health Ave",
		Country:             "Colombia",
		City:                "Medellín",
		State:               "Antioquia",
		PostalCode:          "050001",
		MemberCount:         100,
		ServicesOffered:     []string{"Surgery", "Emergency"},
		AcceptedEPS:         []string{"Nueva EPS", "Famisanar"},
	}

	if dto.OwnerName != "Dr. Jane Smith" {
		t.Errorf("Expected owner name 'Dr. Jane Smith', got %s", dto.OwnerName)
	}

	if dto.Name != "Advanced Medical Center" {
		t.Errorf("Expected clinic name 'Advanced Medical Center', got %s", dto.Name)
	}

	if dto.MemberCount != 100 {
		t.Errorf("Expected member count 100, got %d", dto.MemberCount)
	}

	if len(dto.ServicesOffered) != 2 {
		t.Errorf("Expected 2 services offered, got %d", len(dto.ServicesOffered))
	}

	if len(dto.AcceptedEPS) != 2 {
		t.Errorf("Expected 2 accepted EPS, got %d", len(dto.AcceptedEPS))
	}
}

// Benchmark tests
func BenchmarkClinicalService_CreateClinical(b *testing.B) {
	service, _ := setupClinicalTestService()

	clinicDto := CreateClinicDTO{
		Name:        "Benchmark Clinic",
		Email:       "bench@example.com",
		Description: "Test clinic for benchmarking",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.CreateClinical(clinicDto)
	}
}

func BenchmarkClinicalService_CreateEps(b *testing.B) {
	service, _ := setupClinicalTestService()

	epsDto := CreateEpsDto{
		Eps: []string{"Benchmark EPS"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.CreateEps(epsDto)
	}
}

func BenchmarkClinicalService_GetAllEps(b *testing.B) {
	service, mockRepo := setupClinicalTestService()

	// Prepare benchmark data
	mockRepo.eps = []EPS{
		{ID: "eps-1", Name: "EPS 1"},
		{ID: "eps-2", Name: "EPS 2"},
		{ID: "eps-3", Name: "EPS 3"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GetAllEps(1, 10)
	}
}
