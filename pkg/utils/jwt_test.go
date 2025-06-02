package utils

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Test setup
func TestMain(m *testing.M) {
	// Set up test environment variables
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing-purposes")
	code := m.Run()
	os.Exit(code)
}

func TestGenerateJWT_Success(t *testing.T) {
	userID := "test-user-123"
	validTime := int16(1)

	token, err := GenerateJWT(userID, validTime)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if token == "" {
		t.Error("Expected non-empty token, got empty string")
	}

	// Verify token structure
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		t.Errorf("Failed to parse generated token: %v", err)
	}

	if !parsedToken.Valid {
		t.Error("Generated token is not valid")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Error("Failed to extract claims from token")
	}

	if claims["sub"] != userID {
		t.Errorf("Expected sub claim %s, got %v", userID, claims["sub"])
	}
}

func TestGenerateJWT_ZeroTime(t *testing.T) {
	userID := "test-user-123"
	zeroTime := int16(0)

	token, err := GenerateJWT(userID, zeroTime)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if token == "" {
		t.Error("Expected non-empty token, got empty string")
	}

	// Verify that zero time defaults to 1 hour
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		t.Errorf("Failed to parse generated token: %v", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Error("Failed to extract claims from token")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		t.Error("Failed to extract exp claim")
	}

	// Check if expiration is approximately 1 hour from now
	expectedExp := time.Now().Add(time.Hour).Unix()
	if abs(int64(exp)-expectedExp) > 60 { // Allow 60 seconds tolerance
		t.Errorf("Expected exp around %d, got %d", expectedExp, int64(exp))
	}
}

func TestGenerateJWT_NegativeTime(t *testing.T) {
	userID := "test-user-123"
	negativeTime := int16(-5)

	token, err := GenerateJWT(userID, negativeTime)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if token == "" {
		t.Error("Expected non-empty token, got empty string")
	}

	// Verify that negative time defaults to 1 hour
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		t.Errorf("Failed to parse generated token: %v", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Error("Failed to extract claims from token")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		t.Error("Failed to extract exp claim")
	}

	// Check if expiration is approximately 1 hour from now
	expectedExp := time.Now().Add(time.Hour).Unix()
	if abs(int64(exp)-expectedExp) > 60 { // Allow 60 seconds tolerance
		t.Errorf("Expected exp around %d, got %d", expectedExp, int64(exp))
	}
}

func TestValidateJWT_Success(t *testing.T) {
	userID := "test-user-123"
	validTime := int16(1)

	// Generate a valid token
	token, err := GenerateJWT(userID, validTime)
	if err != nil {
		t.Fatalf("Failed to generate token for test: %v", err)
	}

	// Validate the token
	extractedUserID, err := ValidateJWT(token)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if extractedUserID != userID {
		t.Errorf("Expected user ID %s, got %s", userID, extractedUserID)
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	invalidToken := "invalid.token.string"

	_, err := ValidateJWT(invalidToken)

	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
}

func TestValidateJWT_EmptyToken(t *testing.T) {
	emptyToken := ""

	_, err := ValidateJWT(emptyToken)

	if err == nil {
		t.Error("Expected error for empty token, got nil")
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	// Create an expired token manually
	claims := jwt.MapClaims{
		"sub": "test-user-123",
		"exp": time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		t.Fatalf("Failed to create expired token for test: %v", err)
	}

	_, err = ValidateJWT(signedToken)

	if err == nil {
		t.Error("Expected error for expired token, got nil")
	}
}

func TestValidateJWT_MissingSub(t *testing.T) {
	// Create a token without sub claim
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		t.Fatalf("Failed to create token without sub for test: %v", err)
	}

	_, err = ValidateJWT(signedToken)

	if err == nil {
		t.Error("Expected error for token without sub claim, got nil")
	}

	if err.Error() != "missing or invalid 'sub' claim" {
		t.Errorf("Expected 'missing or invalid 'sub' claim' error, got %s", err.Error())
	}
}

func TestRefreshToken_Success(t *testing.T) {
	userID := "test-user-123"
	validTime := int16(1)

	// Generate a valid token
	originalToken, err := GenerateJWT(userID, validTime)
	if err != nil {
		t.Fatalf("Failed to generate token for test: %v", err)
	}

	// Refresh the token
	newAccessToken, newRefreshToken, err := RefreshToken(originalToken)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if newAccessToken == "" {
		t.Error("Expected non-empty access token, got empty string")
	}

	if newRefreshToken == "" {
		t.Error("Expected non-empty refresh token, got empty string")
	}

	// Verify the new access token
	extractedUserID, err := ValidateJWT(newAccessToken)
	if err != nil {
		t.Errorf("Failed to validate new access token: %v", err)
	}

	if extractedUserID != userID {
		t.Errorf("Expected user ID %s, got %s", userID, extractedUserID)
	}

	// Verify the new refresh token
	extractedUserID, err = ValidateJWT(newRefreshToken)
	if err != nil {
		t.Errorf("Failed to validate new refresh token: %v", err)
	}

	if extractedUserID != userID {
		t.Errorf("Expected user ID %s, got %s", userID, extractedUserID)
	}
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	invalidToken := "invalid.token.string"

	_, _, err := RefreshToken(invalidToken)

	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
}

// Helper function for absolute value
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// Benchmark tests
func BenchmarkGenerateJWT(b *testing.B) {
	userID := "test-user-123"
	validTime := int16(1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateJWT(userID, validTime)
	}
}

func BenchmarkValidateJWT(b *testing.B) {
	userID := "test-user-123"
	validTime := int16(1)

	// Pre-generate token for benchmark
	token, err := GenerateJWT(userID, validTime)
	if err != nil {
		b.Fatalf("Failed to generate token for benchmark: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateJWT(token)
	}
}

func BenchmarkRefreshToken(b *testing.B) {
	userID := "test-user-123"
	validTime := int16(1)

	// Pre-generate token for benchmark
	token, err := GenerateJWT(userID, validTime)
	if err != nil {
		b.Fatalf("Failed to generate token for benchmark: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RefreshToken(token)
	}
}
