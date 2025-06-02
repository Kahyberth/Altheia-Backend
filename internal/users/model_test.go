package users

import (
	"errors"
	"testing"
	"time"
)

func TestUser_Validation(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name: "Valid User",
			user: User{
				ID:             "test-123",
				Name:           "John Doe",
				Email:          "john@example.com",
				Password:       "hashedpassword123",
				Rol:            "patient",
				Phone:          "1234567890",
				DocumentNumber: "123456789",
				Status:         true,
				Gender:         "male",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
				LastLogin:      time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Empty ID",
			user: User{
				Name:           "John Doe",
				Email:          "john@example.com",
				Password:       "hashedpassword123",
				Rol:            "patient",
				Phone:          "1234567890",
				DocumentNumber: "123456789",
				Status:         true,
				Gender:         "male",
			},
			wantErr: true,
		},
		{
			name: "Empty Email",
			user: User{
				ID:             "test-123",
				Name:           "John Doe",
				Password:       "hashedpassword123",
				Rol:            "patient",
				Phone:          "1234567890",
				DocumentNumber: "123456789",
				Status:         true,
				Gender:         "male",
			},
			wantErr: true,
		},
		{
			name: "Invalid Role",
			user: User{
				ID:             "test-123",
				Name:           "John Doe",
				Email:          "john@example.com",
				Password:       "hashedpassword123",
				Rol:            "invalid_role",
				Phone:          "1234567890",
				DocumentNumber: "123456789",
				Status:         true,
				Gender:         "male",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUser(&tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("User validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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

func TestPhysician_Validation(t *testing.T) {
	tests := []struct {
		name      string
		physician Physician
		wantErr   bool
	}{
		{
			name: "Valid Physician",
			physician: Physician{
				ID:                 "physician-123",
				UserID:             "user-123",
				PhysicianSpecialty: "Cardiology",
				LicenseNumber:      "LIC123456",
				Status:             true,
				ClinicID:           stringPtr("clinic-123"),
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Empty UserID",
			physician: Physician{
				ID:                 "physician-123",
				PhysicianSpecialty: "Cardiology",
				LicenseNumber:      "LIC123456",
				Status:             true,
			},
			wantErr: true,
		},
		{
			name: "Empty License Number",
			physician: Physician{
				ID:                 "physician-123",
				UserID:             "user-123",
				PhysicianSpecialty: "Cardiology",
				Status:             true,
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
				ClinicID:  stringPtr("clinic-123"),
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
				ClinicID: stringPtr("clinic-123"),
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

func TestClinicOwner_Validation(t *testing.T) {
	tests := []struct {
		name        string
		clinicOwner ClinicOwner
		wantErr     bool
	}{
		{
			name: "Valid Clinic Owner",
			clinicOwner: ClinicOwner{
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
			clinicOwner: ClinicOwner{
				ID:       "owner-123",
				ClinicID: "clinic-123",
				Status:   true,
			},
			wantErr: true,
		},
		{
			name: "Empty ClinicID",
			clinicOwner: ClinicOwner{
				ID:     "owner-123",
				UserID: "user-123",
				Status: true,
			},
			wantErr: true,
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

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}

// Validation functions
func validateUser(user *User) error {
	if user.ID == "" {
		return errors.New("user ID is required")
	}
	if user.Email == "" {
		return errors.New("user email is required")
	}
	if user.Rol != "patient" && user.Rol != "physician" && user.Rol != "receptionist" && user.Rol != "clinic_owner" {
		return errors.New("invalid user role")
	}
	return nil
}

func validatePatient(patient *Patient) error {
	if patient.UserID == "" {
		return errors.New("patient user ID is required")
	}
	if patient.BloodType != "" && !isValidBloodType(patient.BloodType) {
		return errors.New("invalid blood type")
	}
	return nil
}

func validatePhysician(physician *Physician) error {
	if physician.UserID == "" {
		return errors.New("physician user ID is required")
	}
	if physician.LicenseNumber == "" {
		return errors.New("license number is required")
	}
	return nil
}

func validateReceptionist(receptionist *Receptionist) error {
	if receptionist.UserID == "" {
		return errors.New("receptionist user ID is required")
	}
	if receptionist.ClinicID == nil {
		return errors.New("clinic ID is required")
	}
	return nil
}

func validateClinicOwner(owner *ClinicOwner) error {
	if owner.UserID == "" {
		return errors.New("clinic owner user ID is required")
	}
	if owner.ClinicID == "" {
		return errors.New("clinic ID is required")
	}
	return nil
}

func isValidBloodType(bloodType string) bool {
	validTypes := []string{"A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-"}
	for _, t := range validTypes {
		if bloodType == t {
			return true
		}
	}
	return false
}
