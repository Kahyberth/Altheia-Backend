package appointments

import (
	"Altheia-Backend/internal/users"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid"
	"gorm.io/gorm"
	"time"
)

type Repository interface {
	CreateAppointment(appointment CreateAppointmentDTO) error
	GetAllAppointments() ([]MedicalAppointment, error)
	UpdateAppointmentStatus(appointmentId string, status AppointmentStatus) error
	GetAllAppointmentsByMedicId(physicianId string) ([]AppointmentWithNamesDTO, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository { return &repository{db} }

func (r *repository) GetAllAppointmentsByMedicId(physicianId string) ([]AppointmentWithNamesDTO, error) {
	var appointments []MedicalAppointment
	err := r.db.Model(&MedicalAppointment{}).
		Where("physician_id = ?", physicianId).
		Preload("Patient").
		Preload("Physician").
		Find(&appointments).Error
	if err != nil {
		return nil, fmt.Errorf("error al obtener las citas médicas por ID de médico: %w", err)
	}

	var result []AppointmentWithNamesDTO
	for _, appt := range appointments {
		var patient users.User
		var physician users.User
		// Buscar nombre del paciente
		if err := r.db.Model(&users.User{}).Where("id = ?", appt.Patient.UserID).First(&patient).Error; err != nil {
			patient.Name = ""
		}
		// Buscar nombre del médico
		if err := r.db.Model(&users.User{}).Where("id = ?", appt.Physician.UserID).First(&physician).Error; err != nil {
			physician.Name = ""
		}
		result = append(result, AppointmentWithNamesDTO{
			MedicalAppointment: appt,
			PatientName:        patient.Name,
			PhysicianName:      physician.Name,
		})
	}
	return result, nil
}

func (r *repository) CreateAppointment(appointment CreateAppointmentDTO) error {
	nanoId, _ := gonanoid.Nanoid()

	dateTimeStr := fmt.Sprintf("%s %s", appointment.Date, appointment.Time)

	dateTime, err := time.Parse("2006-01-02 15:04", dateTimeStr)
	if err != nil {
		return fmt.Errorf("formato de fecha u hora inválido: %w", err)
	}

	newAppointment := MedicalAppointment{
		ID:          nanoId,
		PatientId:   appointment.PatientId,
		PhysicianId: appointment.PhysicianId,
		DateTime:    dateTime,
		Status:      string(AppointmentStatusPending),
		Reason:      appointment.Reason,
		Patient:     users.Patient{},
		Physician:   users.Physician{},
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		DeletedAt:   gorm.DeletedAt{},
	}

	err = r.db.Create(&newAppointment).Error

	if err != nil {
		return fmt.Errorf("error al crear la cita médica: %w", err)
	}

	return nil
}

func (r *repository) GetAllAppointments() ([]MedicalAppointment, error) {
	var appointments []MedicalAppointment

	err := r.db.Model(&MedicalAppointment{}).
		Preload("Patient").
		Preload("Physician").
		Find(&appointments).Error

	if err != nil {
		return nil, fmt.Errorf("error al obtener las citas médicas: %w", err)
	}

	return appointments, nil
}

func (r *repository) UpdateAppointmentStatus(appointmentId string, status AppointmentStatus) error {
	err := r.db.Model(&MedicalAppointment{}).
		Where("id = ?", appointmentId).
		Update("status", status).Error

	if err != nil {
		return fmt.Errorf("error al actualizar el estado de la cita médica: %w", err)
	}

	return nil
}
