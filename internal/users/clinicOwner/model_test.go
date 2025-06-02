package clinicOwner

import (
	"Altheia-Backend/internal/users"
	"errors"
	"testing"
	"time"
)

// validateClinicOwner validates a clinic owner's data
func validateClinicOwner(owner *users.ClinicOwner) error {
	if owner.UserID == "" {
		return errors.New("clinic owner user ID is required")
	}
	if owner.ClinicID == "" {
		return errors.New("clinic ID is required")
	}
	return nil
}

func TestClinicOwner_Validation(t *testing.T) {
	tests := []struct {
		name        string
		clinicOwner users.ClinicOwner
		wantErr     bool
	}{
		{
			name: "Valid Clinic Owner",
			clinicOwner: users.ClinicOwner{
				ID:        "owner-123",
				UserID:    "user-123",
				ClinicID:  "clinic-123",
				Status:    true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Empty UserID",
			clinicOwner: users.ClinicOwner{
				ID:       "owner-123",
				ClinicID: "clinic-123",
				Status:   true,
			},
			wantErr: true,
		},
		{
			name: "Empty ClinicID",
			clinicOwner: users.ClinicOwner{
				ID:     "owner-123",
				UserID: "user-123",
				Status: true,
			},
			wantErr: true,
		},
		{
			name: "Invalid Status",
			clinicOwner: users.ClinicOwner{
				ID:       "owner-123",
				UserID:   "user-123",
				ClinicID: "clinic-123",
				Status:   false,
			},
			wantErr: false, // Status can be false
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateClinicOwner(&tt.clinicOwner)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicOwner validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClinicOwner_TimeFields(t *testing.T) {
	now := time.Now()
	owner := users.ClinicOwner{
		ID:        "owner-123",
		UserID:    "user-123",
		ClinicID:  "clinic-123",
		Status:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if owner.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}

	if owner.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}

	if owner.CreatedAt.After(owner.UpdatedAt) {
		t.Error("CreatedAt should not be after UpdatedAt")
	}
}

func TestClinicOwner_StatusTransition(t *testing.T) {
	owner := users.ClinicOwner{
		ID:       "owner-123",
		UserID:   "user-123",
		ClinicID: "clinic-123",
		Status:   true,
	}

	// Test status change
	owner.Status = false
	if owner.Status {
		t.Error("Status should be false after change")
	}

	owner.Status = true
	if !owner.Status {
		t.Error("Status should be true after change")
	}
}

// Benchmark tests
func BenchmarkClinicOwner_Creation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = users.ClinicOwner{
			ID:        "owner-123",
			UserID:    "user-123",
			ClinicID:  "clinic-123",
			Status:    true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}
}

func BenchmarkClinicOwner_Validation(b *testing.B) {
	owner := users.ClinicOwner{
		ID:        "owner-123",
		UserID:    "user-123",
		ClinicID:  "clinic-123",
		Status:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateClinicOwner(&owner)
	}
}
