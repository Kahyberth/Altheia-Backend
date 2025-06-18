package patient

import (
	"testing"
	"time"
)

func TestPatient_Validation(t *testing.T) {
	tests := []struct {
		name    string
		patient Patient
		wantErr bool
	}{
		{
			name: "Valid Patient",
			patient: Patient{
				ID:          "patient-123",
				UserID:      "user-123",
				DateOfBirth: "1990-01-01",
				Address:     "123 Main St",
				Eps:         "Sura",
				BloodType:   "O+",
				Status:      true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Empty UserID",
			patient: Patient{
				ID:          "patient-123",
				DateOfBirth: "1990-01-01",
				Address:     "123 Main St",
				Eps:         "Sura",
				BloodType:   "O+",
				Status:      true,
			},
			wantErr: true,
		},
		{
			name: "Invalid Blood Type",
			patient: Patient{
				ID:          "patient-123",
				UserID:      "user-123",
				DateOfBirth: "1990-01-01",
				Address:     "123 Main St",
				Eps:         "Sura",
				BloodType:   "Invalid",
				Status:      true,
			},
			wantErr: true,
		},
		{
			name: "Invalid Date of Birth",
			patient: Patient{
				ID:          "patient-123",
				UserID:      "user-123",
				DateOfBirth: "invalid-date",
				Address:     "123 Main St",
				Eps:         "Sura",
				BloodType:   "O+",
				Status:      true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePatient(&tt.patient)
			if (err != nil) != tt.wantErr {
				t.Errorf("Patient validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPatient_TimeFields(t *testing.T) {
	now := time.Now()
	patient := Patient{
		ID:          "patient-123",
		UserID:      "user-123",
		DateOfBirth: "1990-01-01",
		Address:     "123 Main St",
		Eps:         "Sura",
		BloodType:   "O+",
		Status:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if patient.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}

	if patient.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}

	if patient.CreatedAt.After(patient.UpdatedAt) {
		t.Error("CreatedAt should not be after UpdatedAt")
	}
}

func TestPatient_BloodTypeValidation(t *testing.T) {
	validBloodTypes := []string{"A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-"}
	invalidBloodTypes := []string{"", "Invalid", "A", "B", "O", "AB"}

	patient := Patient{
		ID:          "patient-123",
		UserID:      "user-123",
		DateOfBirth: "1990-01-01",
		Address:     "123 Main St",
		Eps:         "Sura",
		Status:      true,
	}

	for _, bt := range validBloodTypes {
		patient.BloodType = bt
		if err := validatePatient(&patient); err != nil {
			t.Errorf("Blood type %s should be valid, got error: %v", bt, err)
		}
	}

	for _, bt := range invalidBloodTypes {
		patient.BloodType = bt
		if err := validatePatient(&patient); err == nil {
			t.Errorf("Blood type %s should be invalid, got no error", bt)
		}
	}
}

func BenchmarkPatient_Creation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Patient{
			ID:          "patient-123",
			UserID:      "user-123",
			DateOfBirth: "1990-01-01",
			Address:     "123 Main St",
			Eps:         "Sura",
			BloodType:   "O+",
			Status:      true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
	}
}

func BenchmarkPatient_Validation(b *testing.B) {
	patient := Patient{
		ID:          "patient-123",
		UserID:      "user-123",
		DateOfBirth: "1990-01-01",
		Address:     "123 Main St",
		Eps:         "Sura",
		BloodType:   "O+",
		Status:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validatePatient(&patient)
	}
}
