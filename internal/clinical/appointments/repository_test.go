package appointments

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	return gormDB, mock
}

func TestRepository_CreateAppointment(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewRepository(db)

	tests := []struct {
		name        string
		appointment CreateAppointmentDTO
		wantErr     bool
	}{
		{
			name: "Valid Appointment",
			appointment: CreateAppointmentDTO{
				PatientId:   "patient-123",
				PhysicianId: "physician-123",
				Date:        "2024-03-20",
				Time:        "14:30",
				Reason:      "Regular checkup",
			},
			wantErr: false,
		},
		{
			name: "Invalid Date Format",
			appointment: CreateAppointmentDTO{
				PatientId:   "patient-123",
				PhysicianId: "physician-123",
				Date:        "invalid-date",
				Time:        "14:30",
				Reason:      "Regular checkup",
			},
			wantErr: true,
		},
		{
			name: "Invalid Time Format",
			appointment: CreateAppointmentDTO{
				PatientId:   "patient-123",
				PhysicianId: "physician-123",
				Date:        "2024-03-20",
				Time:        "invalid-time",
				Reason:      "Regular checkup",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "medical_appointments" ("id","patient_id","physician_id","date_time","status","reason") VALUES ($1,$2,$3,$4,$5,$6)`).
					WithArgs(
						sqlmock.AnyArg(), // ID
						tt.appointment.PatientId,
						tt.appointment.PhysicianId,
						sqlmock.AnyArg(), // DateTime
						string(AppointmentStatusPending),
						tt.appointment.Reason,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			}

			err := repo.CreateAppointment(tt.appointment)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAppointment() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil && !tt.wantErr {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestRepository_GetAllAppointments(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewRepository(db)

	expectedAppointments := []MedicalAppointment{
		{
			ID:          "appt-1",
			PatientId:   "patient-1",
			PhysicianId: "physician-1",
			DateTime:    time.Now().Add(24 * time.Hour),
			Status:      string(AppointmentStatusPending),
			Reason:      "Checkup 1",
		},
		{
			ID:          "appt-2",
			PatientId:   "patient-2",
			PhysicianId: "physician-2",
			DateTime:    time.Now().Add(48 * time.Hour),
			Status:      string(AppointmentStatusConfirmed),
			Reason:      "Checkup 2",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "patient_id", "physician_id", "date_time", "status", "reason", "created_at", "updated_at", "deleted_at"}).
		AddRow(expectedAppointments[0].ID, expectedAppointments[0].PatientId, expectedAppointments[0].PhysicianId,
			expectedAppointments[0].DateTime, expectedAppointments[0].Status, expectedAppointments[0].Reason,
			time.Now(), time.Now(), nil).
		AddRow(expectedAppointments[1].ID, expectedAppointments[1].PatientId, expectedAppointments[1].PhysicianId,
			expectedAppointments[1].DateTime, expectedAppointments[1].Status, expectedAppointments[1].Reason,
			time.Now(), time.Now(), nil)

	mock.ExpectQuery("SELECT \\* FROM \"medical_appointments\" WHERE \"medical_appointments\".\"deleted_at\" IS NULL").
		WillReturnRows(rows)

	appointments, err := repo.GetAllAppointments()
	if err != nil {
		t.Errorf("GetAllAppointments() error = %v", err)
	}

	if len(appointments) != len(expectedAppointments) {
		t.Errorf("GetAllAppointments() got %d appointments, want %d", len(appointments), len(expectedAppointments))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestRepository_UpdateAppointmentStatus(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewRepository(db)

	tests := []struct {
		name          string
		appointmentId string
		status        AppointmentStatus
		wantErr       bool
	}{
		{
			name:          "Valid Status Update",
			appointmentId: "appt-123",
			status:        AppointmentStatusConfirmed,
			wantErr:       false,
		},
		{
			name:          "Non-existent Appointment",
			appointmentId: "non-existent",
			status:        AppointmentStatusCancelled,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE medical_appointments").
					WithArgs(string(tt.status), sqlmock.AnyArg(), tt.appointmentId).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			} else {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE medical_appointments").
					WithArgs(string(tt.status), sqlmock.AnyArg(), tt.appointmentId).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			}

			err := repo.UpdateAppointmentStatus(tt.appointmentId, tt.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateAppointmentStatus() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil && !tt.wantErr {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestRepository_GetAllAppointmentsByMedicId(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewRepository(db)

	physicianId := "physician-123"
	expectedAppointments := []AppointmentWithNamesDTO{
		{
			MedicalAppointment: MedicalAppointment{
				ID:          "appt-1",
				PatientId:   "patient-1",
				PhysicianId: physicianId,
				DateTime:    time.Now().Add(24 * time.Hour),
				Status:      string(AppointmentStatusPending),
				Reason:      "Checkup 1",
			},
			PatientName:   "John Doe",
			PhysicianName: "Dr. Smith",
		},
	}

	// Mock the main appointments query
	rows := sqlmock.NewRows([]string{"id", "patient_id", "physician_id", "date_time", "status", "reason", "created_at", "updated_at", "deleted_at"}).
		AddRow(expectedAppointments[0].ID, expectedAppointments[0].PatientId, expectedAppointments[0].PhysicianId,
			expectedAppointments[0].DateTime, expectedAppointments[0].Status, expectedAppointments[0].Reason,
			time.Now(), time.Now(), nil)

	mock.ExpectQuery("SELECT \\* FROM \"medical_appointments\" WHERE physician_id = \\$1 AND \"medical_appointments\".\"deleted_at\" IS NULL").
		WithArgs(physicianId).
		WillReturnRows(rows)

	// Mock the patient user query
	patientRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow("user-1", expectedAppointments[0].PatientName)
	mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE \"users\".\"deleted_at\" IS NULL AND \"users\".\"id\" = \\$1").
		WithArgs(expectedAppointments[0].PatientId).
		WillReturnRows(patientRows)

	// Mock the physician user query
	physicianRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow("user-2", expectedAppointments[0].PhysicianName)
	mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE \"users\".\"deleted_at\" IS NULL AND \"users\".\"id\" = \\$1").
		WithArgs(expectedAppointments[0].PhysicianId).
		WillReturnRows(physicianRows)

	appointments, err := repo.GetAllAppointmentsByMedicId(physicianId)
	if err != nil {
		t.Errorf("GetAllAppointmentsByMedicId() error = %v", err)
	}

	if len(appointments) != len(expectedAppointments) {
		t.Errorf("GetAllAppointmentsByMedicId() got %d appointments, want %d. Appointments: %+v", len(appointments), len(expectedAppointments), appointments)
	}

	if len(appointments) > 0 && appointments[0].PatientName != expectedAppointments[0].PatientName {
		t.Errorf("GetAllAppointmentsByMedicId() got patient name %s, want %s",
			appointments[0].PatientName, expectedAppointments[0].PatientName)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
