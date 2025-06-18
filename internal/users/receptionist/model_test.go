package receptionist

import (
	"testing"
	"time"
)

func TestReceptionist_Validation(t *testing.T) {
	tests := []struct {
		name         string
		receptionist Receptionist
		wantErr      bool
	}{
		{
			name: "Valid Receptionist",
			receptionist: Receptionist{
				ID:        "receptionist-123",
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
			receptionist: Receptionist{
				ID:       "receptionist-123",
				ClinicID: "clinic-123",
				Status:   true,
			},
			wantErr: true,
		},
		{
			name: "Empty ClinicID",
			receptionist: Receptionist{
				ID:     "receptionist-123",
				UserID: "user-123",
				Status: true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateReceptionist(&tt.receptionist)
			if (err != nil) != tt.wantErr {
				t.Errorf("Receptionist validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReceptionist_TimeFields(t *testing.T) {
	now := time.Now()
	receptionist := Receptionist{
		ID:        "receptionist-123",
		UserID:    "user-123",
		ClinicID:  "clinic-123",
		Status:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if receptionist.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}

	if receptionist.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}

	if receptionist.CreatedAt.After(receptionist.UpdatedAt) {
		t.Error("CreatedAt should not be after UpdatedAt")
	}
}

func BenchmarkReceptionist_Creation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Receptionist{
			ID:        "receptionist-123",
			UserID:    "user-123",
			ClinicID:  "clinic-123",
			Status:    true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}
}

func BenchmarkReceptionist_Validation(b *testing.B) {
	receptionist := Receptionist{
		ID:        "receptionist-123",
		UserID:    "user-123",
		ClinicID:  "clinic-123",
		Status:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateReceptionist(&receptionist)
	}
}
