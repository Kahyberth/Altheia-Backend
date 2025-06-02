package appointments

import (
	"testing"
	"time"
)

func TestMedicalAppointment_Validation(t *testing.T) {
	tests := []struct {
		name        string
		appointment MedicalAppointment
		wantErr     bool
	}{
		{
			name: "Valid Appointment",
			appointment: MedicalAppointment{
				ID:          "appt-123",
				PatientId:   "patient-123",
				PhysicianId: "physician-123",
				DateTime:    time.Now().Add(24 * time.Hour), // Tomorrow
				Status:      string(AppointmentStatusPending),
				Reason:      "Regular checkup",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Empty Patient ID",
			appointment: MedicalAppointment{
				ID:          "appt-123",
				PhysicianId: "physician-123",
				DateTime:    time.Now().Add(24 * time.Hour),
				Status:      string(AppointmentStatusPending),
				Reason:      "Regular checkup",
			},
			wantErr: true,
		},
		{
			name: "Empty Physician ID",
			appointment: MedicalAppointment{
				ID:        "appt-123",
				PatientId: "patient-123",
				DateTime:  time.Now().Add(24 * time.Hour),
				Status:    string(AppointmentStatusPending),
				Reason:    "Regular checkup",
			},
			wantErr: true,
		},
		{
			name: "Past DateTime",
			appointment: MedicalAppointment{
				ID:          "appt-123",
				PatientId:   "patient-123",
				PhysicianId: "physician-123",
				DateTime:    time.Now().Add(-24 * time.Hour), // Yesterday
				Status:      string(AppointmentStatusPending),
				Reason:      "Regular checkup",
			},
			wantErr: true,
		},
		{
			name: "Invalid Status",
			appointment: MedicalAppointment{
				ID:          "appt-123",
				PatientId:   "patient-123",
				PhysicianId: "physician-123",
				DateTime:    time.Now().Add(24 * time.Hour),
				Status:      "INVALID_STATUS",
				Reason:      "Regular checkup",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAppointment(&tt.appointment)
			if (err != nil) != tt.wantErr {
				t.Errorf("MedicalAppointment validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMedicalAppointment_TimeFields(t *testing.T) {
	now := time.Now()
	appointment := MedicalAppointment{
		ID:          "appt-123",
		PatientId:   "patient-123",
		PhysicianId: "physician-123",
		DateTime:    now.Add(24 * time.Hour),
		Status:      string(AppointmentStatusPending),
		Reason:      "Regular checkup",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if appointment.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}

	if appointment.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}

	if appointment.CreatedAt.After(appointment.UpdatedAt) {
		t.Error("CreatedAt should not be after UpdatedAt")
	}

	if appointment.DateTime.Before(now) {
		t.Error("DateTime should be in the future")
	}
}

func TestMedicalAppointment_StatusValidation(t *testing.T) {
	validStatuses := []string{
		string(AppointmentStatusPending),
		string(AppointmentStatusConfirmed),
		string(AppointmentStatusCancelled),
		string(AppointmentStatusCompleted),
	}
	invalidStatuses := []string{"", "INVALID", "PENDING", "CONFIRMED"}

	appointment := MedicalAppointment{
		ID:          "appt-123",
		PatientId:   "patient-123",
		PhysicianId: "physician-123",
		DateTime:    time.Now().Add(24 * time.Hour),
		Reason:      "Regular checkup",
	}

	// Test valid statuses
	for _, status := range validStatuses {
		appointment.Status = status
		if err := validateAppointment(&appointment); err != nil {
			t.Errorf("Status %s should be valid, got error: %v", status, err)
		}
	}

	// Test invalid statuses
	for _, status := range invalidStatuses {
		appointment.Status = status
		if err := validateAppointment(&appointment); err == nil {
			t.Errorf("Status %s should be invalid, got no error", status)
		}
	}
}

// Benchmark tests
func BenchmarkMedicalAppointment_Creation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MedicalAppointment{
			ID:          "appt-123",
			PatientId:   "patient-123",
			PhysicianId: "physician-123",
			DateTime:    time.Now().Add(24 * time.Hour),
			Status:      string(AppointmentStatusPending),
			Reason:      "Regular checkup",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
	}
}

func BenchmarkMedicalAppointment_Validation(b *testing.B) {
	appointment := MedicalAppointment{
		ID:          "appt-123",
		PatientId:   "patient-123",
		PhysicianId: "physician-123",
		DateTime:    time.Now().Add(24 * time.Hour),
		Status:      string(AppointmentStatusPending),
		Reason:      "Regular checkup",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateAppointment(&appointment)
	}
}
