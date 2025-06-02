package physician

import (
	"testing"
	"time"
)

func TestPhysician_Validation(t *testing.T) {
	tests := []struct {
		name      string
		physician Physician
		wantErr   bool
	}{
		{
			name: "Valid Physician",
			physician: Physician{
				ID:            "physician-123",
				UserID:        "user-123",
				Specialty:     "Cardiology",
				LicenseNumber: "MD123456",
				Status:        true,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Empty UserID",
			physician: Physician{
				ID:            "physician-123",
				Specialty:     "Cardiology",
				LicenseNumber: "MD123456",
				Status:        true,
			},
			wantErr: true,
		},
		{
			name: "Empty Specialty",
			physician: Physician{
				ID:            "physician-123",
				UserID:        "user-123",
				LicenseNumber: "MD123456",
				Status:        true,
			},
			wantErr: true,
		},
		{
			name: "Empty License Number",
			physician: Physician{
				ID:        "physician-123",
				UserID:    "user-123",
				Specialty: "Cardiology",
				Status:    true,
			},
			wantErr: true,
		},
		{
			name: "Invalid License Number",
			physician: Physician{
				ID:            "physician-123",
				UserID:        "user-123",
				Specialty:     "Cardiology",
				LicenseNumber: "MD@123456",
				Status:        true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePhysician(&tt.physician)
			if (err != nil) != tt.wantErr {
				t.Errorf("Physician validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPhysician_TimeFields(t *testing.T) {
	now := time.Now()
	physician := Physician{
		ID:            "physician-123",
		UserID:        "user-123",
		Specialty:     "Cardiology",
		LicenseNumber: "MD123456",
		Status:        true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if physician.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}

	if physician.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}

	if physician.CreatedAt.After(physician.UpdatedAt) {
		t.Error("CreatedAt should not be after UpdatedAt")
	}
}

func TestPhysician_LicenseNumberValidation(t *testing.T) {
	validLicenseNumbers := []string{"MD123456", "DR789ABC", "123456", "ABC123"}
	invalidLicenseNumbers := []string{"", "MD@123", "DR#456", "ABC 123"}

	physician := Physician{
		ID:        "physician-123",
		UserID:    "user-123",
		Specialty: "Cardiology",
		Status:    true,
	}

	// Test valid license numbers
	for _, ln := range validLicenseNumbers {
		physician.LicenseNumber = ln
		if err := validatePhysician(&physician); err != nil {
			t.Errorf("License number %s should be valid, got error: %v", ln, err)
		}
	}

	// Test invalid license numbers
	for _, ln := range invalidLicenseNumbers {
		physician.LicenseNumber = ln
		if err := validatePhysician(&physician); err == nil {
			t.Errorf("License number %s should be invalid, got no error", ln)
		}
	}
}

// Benchmark tests
func BenchmarkPhysician_Creation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Physician{
			ID:            "physician-123",
			UserID:        "user-123",
			Specialty:     "Cardiology",
			LicenseNumber: "MD123456",
			Status:        true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
	}
}

func BenchmarkPhysician_Validation(b *testing.B) {
	physician := Physician{
		ID:            "physician-123",
		UserID:        "user-123",
		Specialty:     "Cardiology",
		LicenseNumber: "MD123456",
		Status:        true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validatePhysician(&physician)
	}
}
