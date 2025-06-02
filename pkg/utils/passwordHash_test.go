package utils

import (
	"strings"
	"testing"
)

func TestHashPassword_Success(t *testing.T) {
	password := "testpassword123"

	hashedPassword, err := HashPassword(password)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if hashedPassword == "" {
		t.Error("Expected non-empty hashed password, got empty string")
	}

	if hashedPassword == password {
		t.Error("Hashed password should not be the same as original password")
	}

	// Check if it looks like a bcrypt hash
	if !strings.HasPrefix(hashedPassword, "$2a$") && !strings.HasPrefix(hashedPassword, "$2b$") && !strings.HasPrefix(hashedPassword, "$2y$") {
		t.Error("Hashed password doesn't look like a bcrypt hash")
	}
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	password := ""

	hashedPassword, err := HashPassword(password)

	if err != nil {
		t.Errorf("Expected no error for empty password, got %v", err)
	}

	if hashedPassword == "" {
		t.Error("Expected non-empty hashed password even for empty input")
	}
}

func TestHashPassword_LongPassword(t *testing.T) {
	// Create a password that's just under bcrypt's 72-byte limit
	password := strings.Repeat("a", 70)

	hashedPassword, err := HashPassword(password)

	if err != nil {
		t.Errorf("Expected no error for long password, got %v", err)
	}

	if hashedPassword == "" {
		t.Error("Expected non-empty hashed password, got empty string")
	}

	if hashedPassword == password {
		t.Error("Hashed password should not be the same as original password")
	}
}

func TestHashPassword_SpecialCharacters(t *testing.T) {
	password := "p@ssw0rd!@#$%^&*()_+-=[]{}|;':\",./<>?"

	hashedPassword, err := HashPassword(password)

	if err != nil {
		t.Errorf("Expected no error for password with special characters, got %v", err)
	}

	if hashedPassword == "" {
		t.Error("Expected non-empty hashed password, got empty string")
	}

	if hashedPassword == password {
		t.Error("Hashed password should not be the same as original password")
	}
}

func TestCheckPasswordHash_Success(t *testing.T) {
	password := "testpassword123"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password for test: %v", err)
	}

	isValid := CheckPasswordHash(password, hashedPassword)

	if !isValid {
		t.Error("Expected password check to return true for correct password")
	}
}

func TestCheckPasswordHash_WrongPassword(t *testing.T) {
	password := "testpassword123"
	wrongPassword := "wrongpassword456"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password for test: %v", err)
	}

	isValid := CheckPasswordHash(wrongPassword, hashedPassword)

	if isValid {
		t.Error("Expected password check to return false for wrong password")
	}
}

func TestCheckPasswordHash_EmptyPassword(t *testing.T) {
	password := "testpassword123"
	emptyPassword := ""

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password for test: %v", err)
	}

	isValid := CheckPasswordHash(emptyPassword, hashedPassword)

	if isValid {
		t.Error("Expected password check to return false for empty password")
	}
}

func TestCheckPasswordHash_EmptyHash(t *testing.T) {
	password := "testpassword123"
	emptyHash := ""

	isValid := CheckPasswordHash(password, emptyHash)

	if isValid {
		t.Error("Expected password check to return false for empty hash")
	}
}

func TestCheckPasswordHash_InvalidHash(t *testing.T) {
	password := "testpassword123"
	invalidHash := "invalid-hash-format"

	isValid := CheckPasswordHash(password, invalidHash)

	if isValid {
		t.Error("Expected password check to return false for invalid hash")
	}
}

func TestCheckPasswordHash_CaseSensitive(t *testing.T) {
	password := "TestPassword123"
	wrongCasePassword := "testpassword123"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password for test: %v", err)
	}

	isValid := CheckPasswordHash(wrongCasePassword, hashedPassword)

	if isValid {
		t.Error("Expected password check to be case sensitive and return false")
	}
}

func TestHashPassword_DifferentHashesForSamePassword(t *testing.T) {
	password := "testpassword123"

	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)

	if err1 != nil || err2 != nil {
		t.Fatalf("Failed to hash passwords: %v, %v", err1, err2)
	}

	// Hashes should be different due to salt
	if hash1 == hash2 {
		t.Error("Expected different hashes for same password due to random salt")
	}

	// But both should verify correctly
	if !CheckPasswordHash(password, hash1) {
		t.Error("First hash should verify correctly")
	}

	if !CheckPasswordHash(password, hash2) {
		t.Error("Second hash should verify correctly")
	}
}

func TestPasswordWorkflow_EndToEnd(t *testing.T) {
	// Test scenarios for a complete password workflow
	testCases := []struct {
		name     string
		password string
	}{
		{"Simple password", "password123"},
		{"Complex password", "P@ssw0rd!@#$%^&*()"},
		{"Long password", strings.Repeat("a", 70)}, // Just under bcrypt's limit
		{"Unicode password", "пароль123密码"},
		{"Spaces in password", "password with spaces"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Hash the password
			hashedPassword, err := HashPassword(tc.password)
			if err != nil {
				t.Errorf("Failed to hash password: %v", err)
				return
			}

			// Verify correct password
			if !CheckPasswordHash(tc.password, hashedPassword) {
				t.Error("Correct password should verify successfully")
			}

			// Verify wrong password fails
			wrongPassword := tc.password + "wrong"
			if CheckPasswordHash(wrongPassword, hashedPassword) {
				t.Error("Wrong password should fail verification")
			}
		})
	}
}

// Benchmark tests
func BenchmarkHashPassword(b *testing.B) {
	password := "testpassword123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HashPassword(password)
	}
}

func BenchmarkCheckPasswordHash(b *testing.B) {
	password := "testpassword123"

	// Pre-generate hash for benchmark
	hashedPassword, err := HashPassword(password)
	if err != nil {
		b.Fatalf("Failed to hash password for benchmark: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CheckPasswordHash(password, hashedPassword)
	}
}

func BenchmarkHashPassword_Different_Lengths(b *testing.B) {
	testCases := []struct {
		name     string
		password string
	}{
		{"Short", "pass"},
		{"Medium", "testpassword123"},
		{"Long", strings.Repeat("testpassword", 10)},
		{"Very Long", strings.Repeat("testpassword", 50)},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				HashPassword(tc.password)
			}
		})
	}
}
