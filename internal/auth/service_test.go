package auth

import (
	"Altheia-Backend/internal/users"
	"Altheia-Backend/pkg/utils"
	"errors"
	"testing"
)

// Mock Repository for testing
type mockRepository struct {
	users map[string]*users.User
	err   error
}

func (m *mockRepository) FindByEmail(email string) (*users.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockRepository) FindByID(id string) (*users.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	if user, ok := m.users[id]; ok {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func setupTestService() (Service, *mockRepository) {
	mockRepo := &mockRepository{
		users: make(map[string]*users.User),
	}
	service := NewService(mockRepo)
	return service, mockRepo
}

func TestService_Login_Success(t *testing.T) {
	service, mockRepo := setupTestService()

	// Create a test user with hashed password
	hashedPassword, _ := utils.HashPassword("testpassword123")
	testUser := &users.User{
		ID:       "test-user-1",
		Name:     "Test User",
		Email:    "test@example.com",
		Password: hashedPassword,
		Rol:      "patient",
	}
	mockRepo.users["test-user-1"] = testUser

	// Test successful login
	userInfo, accessToken, refreshToken, err := service.Login("test@example.com", "testpassword123")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if userInfo.ID != "test-user-1" {
		t.Errorf("Expected user ID 'test-user-1', got %s", userInfo.ID)
	}

	if userInfo.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got %s", userInfo.Email)
	}

	if userInfo.Role != "patient" {
		t.Errorf("Expected role 'patient', got %s", userInfo.Role)
	}

	if accessToken == "" {
		t.Error("Expected access token, got empty string")
	}

	if refreshToken == "" {
		t.Error("Expected refresh token, got empty string")
	}
}

func TestService_Login_InvalidEmail(t *testing.T) {
	service, _ := setupTestService()

	// Test login with non-existent email
	_, _, _, err := service.Login("nonexistent@example.com", "password")

	if err == nil {
		t.Error("Expected error for invalid email, got nil")
	}

	if err.Error() != "invalid credentials" {
		t.Errorf("Expected 'invalid credentials' error, got %s", err.Error())
	}
}

func TestService_Login_InvalidPassword(t *testing.T) {
	service, mockRepo := setupTestService()

	// Create a test user
	hashedPassword, _ := utils.HashPassword("correctpassword")
	testUser := &users.User{
		ID:       "test-user-1",
		Email:    "test@example.com",
		Password: hashedPassword,
	}
	mockRepo.users["test-user-1"] = testUser

	// Test login with wrong password
	_, _, _, err := service.Login("test@example.com", "wrongpassword")

	if err == nil {
		t.Error("Expected error for invalid password, got nil")
	}

	if err.Error() != "invalid credentials" {
		t.Errorf("Expected 'invalid credentials' error, got %s", err.Error())
	}
}

func TestService_Login_RepositoryError(t *testing.T) {
	service, mockRepo := setupTestService()

	// Set repository to return error
	mockRepo.err = errors.New("database error")

	// Test login with repository error
	_, _, _, err := service.Login("test@example.com", "password")

	if err == nil {
		t.Error("Expected error from repository, got nil")
	}

	if err.Error() != "invalid credentials" {
		t.Errorf("Expected 'invalid credentials' error, got %s", err.Error())
	}
}

func TestService_GetProfile_Success(t *testing.T) {
	service, mockRepo := setupTestService()

	// Create a test user
	testUser := &users.User{
		ID:    "test-user-1",
		Name:  "Test User",
		Email: "test@example.com",
		Rol:   "physician",
	}
	mockRepo.users["test-user-1"] = testUser

	// Test getting profile
	user, err := service.GetProfile("test-user-1")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("Expected user, got nil")
	}

	if user.ID != "test-user-1" {
		t.Errorf("Expected user ID 'test-user-1', got %s", user.ID)
	}

	if user.Name != "Test User" {
		t.Errorf("Expected name 'Test User', got %s", user.Name)
	}
}

func TestService_GetProfile_UserNotFound(t *testing.T) {
	service, _ := setupTestService()

	// Test getting non-existent profile
	user, err := service.GetProfile("non-existent-id")

	if err == nil {
		t.Error("Expected error for non-existent user, got nil")
	}

	if user != nil {
		t.Error("Expected nil user, got user object")
	}
}

func TestService_GetProfile_RepositoryError(t *testing.T) {
	service, mockRepo := setupTestService()

	// Set repository to return error
	mockRepo.err = errors.New("database error")

	// Test getting profile with repository error
	user, err := service.GetProfile("test-user-1")

	if err == nil {
		t.Error("Expected error from repository, got nil")
	}

	if user != nil {
		t.Error("Expected nil user, got user object")
	}
}

func TestUserInfo_Creation(t *testing.T) {
	userInfo := UserInfo{
		ID:    "test-123",
		Name:  "John Doe",
		Email: "john@example.com",
		Role:  "admin",
	}

	if userInfo.ID != "test-123" {
		t.Errorf("Expected ID 'test-123', got %s", userInfo.ID)
	}

	if userInfo.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got %s", userInfo.Name)
	}

	if userInfo.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got %s", userInfo.Email)
	}

	if userInfo.Role != "admin" {
		t.Errorf("Expected role 'admin', got %s", userInfo.Role)
	}
}

// Benchmark tests
func BenchmarkService_Login(b *testing.B) {
	service, mockRepo := setupTestService()

	hashedPassword, _ := utils.HashPassword("testpassword123")
	testUser := &users.User{
		ID:       "test-user-1",
		Email:    "test@example.com",
		Password: hashedPassword,
		Rol:      "patient",
	}
	mockRepo.users["test-user-1"] = testUser

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.Login("test@example.com", "testpassword123")
	}
}
