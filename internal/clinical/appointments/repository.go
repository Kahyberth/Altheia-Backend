package appointments

import (
	"Altheia-Backend/internal/clinical"
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
	CancelAppointment(appointmentId string) error
	GetAllAppointmentsByUserId(userId string) ([]AppointmentWithNamesDTO, error)
	ReScheduleAppointment(appointmentId string, newDateTime time.Time) error
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
		var clinic clinical.ClinicInformation
		// Buscar nombre del paciente
		if err := r.db.Model(&users.User{}).Where("id = ?", appt.Patient.UserID).First(&patient).Error; err != nil {
			patient.Name = ""
		}
		// Buscar nombre del médico
		if err := r.db.Model(&users.User{}).Where("id = ?", appt.Physician.UserID).First(&physician).Error; err != nil {
			physician.Name = ""
		}
		// Buscar nombre de la clínica del médico
		if err := r.db.Model(&clinical.ClinicInformation{}).Where("clinic_id = ?", appt.Physician.ClinicID).First(&clinic).Error; err != nil {
			clinic.ClinicName = ""
			clinic.City = ""
			clinic.Address = ""
			clinic.ClinicID = ""
		}

		result = append(result, AppointmentWithNamesDTO{
			MedicalAppointment: appt,
			PatientName:        patient.Name,
			PatientGender:      patient.Gender,
			PatientEmail:       patient.Email,
			PatientPhone:       patient.Phone,
			PhysicianName:      physician.Name,
			PhysicianGender:    physician.Gender,
			PhysicianEmail:     physician.Email,
			PhysicianPhone:     physician.Phone,
			ClinicName:         clinic.ClinicName,
			ClinicCity:         clinic.City,
			ClinicAddress:      clinic.Address,
			ClinicId:           clinic.ClinicID,
		})
	}
	return result, nil
}

func (r *repository) CreateAppointment(appointment CreateAppointmentDTO) error {
	nanoId, _ := gonanoid.Nanoid()

	// Configurar la zona horaria de Bogotá
	loc, err := time.LoadLocation("America/Bogota")
	if err != nil {
		return fmt.Errorf("error al cargar la zona horaria: %w", err)
	}

	dateTimeStr := fmt.Sprintf("%s %s", appointment.Date, appointment.Time)

	// Parsear la fecha especificando la zona horaria de Bogotá
	dateTime, err := time.ParseInLocation("2006-01-02 15:04", dateTimeStr, loc)
	if err != nil {
		return fmt.Errorf("formato de fecha u hora inválido: %w", err)
	}

	// Primero verificamos que existan el paciente y el médico
	var patient users.Patient
	var physician users.Physician

	if err := r.db.Where("id = ?", appointment.PatientId).First(&patient).Error; err != nil {
		return fmt.Errorf("paciente no encontrado: %w", err)
	}

	if err := r.db.Where("id = ?", appointment.PhysicianId).First(&physician).Error; err != nil {
		return fmt.Errorf("médico no encontrado: %w", err)
	}

	newAppointment := MedicalAppointment{
		ID:          nanoId,
		PatientId:   appointment.PatientId,
		PhysicianId: appointment.PhysicianId,
		DateTime:    dateTime,
		Status:      string(AppointmentStatusPending),
		Reason:      appointment.Reason,
		Patient:     patient,
		Physician:   physician,
	}

	if err := r.db.Create(&newAppointment).Error; err != nil {
		return fmt.Errorf("error al crear la cita médica: %w", err)
	}

	return nil
}

func (r *repository) ReScheduleAppointment(appointmentId string, newDateTime time.Time) error {
	// Configurar la zona horaria de Bogotá
	loc, err := time.LoadLocation("America/Bogota")
	if err != nil {
		return fmt.Errorf("error al cargar la zona horaria: %w", err)
	}

	// Convertir la nueva fecha y hora a la zona horaria de Bogotá
	newDateTime = newDateTime.In(loc)

	err = r.db.Model(&MedicalAppointment{}).
		Where("id = ?", appointmentId).
		Update("date_time", newDateTime).
		Update("status", string(AppointmentStatusRescheduled)).Error

	if err != nil {
		return fmt.Errorf("error al reprogramar la cita médica: %w", err)
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

func (r *repository) CancelAppointment(appointmentId string) error {
	err := r.db.Model(&MedicalAppointment{}).
		Where("id = ?", appointmentId).
		Update("status", AppointmentStatusCancelled).Error

	if err != nil {
		return fmt.Errorf("error al cancelar la cita médica: %w", err)
	}

	return nil
}

func (r *repository) GetAllAppointmentsByUserId(userId string) ([]AppointmentWithNamesDTO, error) {
	var appointments []MedicalAppointment

	// Buscar todas las citas donde el usuario es paciente o médico usando joins
	err := r.db.Joins("JOIN patients ON medical_appointments.patient_id = patients.id").
		Joins("JOIN physicians ON medical_appointments.physician_id = physicians.id").
		Where("patients.user_id = ? OR physicians.user_id = ?", userId, userId).
		Preload("Patient").
		Preload("Physician").
		Find(&appointments).Error

	if err != nil {
		return nil, fmt.Errorf("error al obtener las citas médicas: %w", err)
	}

	var result []AppointmentWithNamesDTO
	for _, appt := range appointments {
		var patientUser users.User
		var physicianUser users.User
		var clinic clinical.ClinicInformation

		// Obtener datos del paciente
		if err := r.db.Model(&users.User{}).Where("id = ?", appt.Patient.UserID).First(&patientUser).Error; err != nil {
			continue
		}

		// Obtener datos del médico
		if err := r.db.Model(&users.User{}).Where("id = ?", appt.Physician.UserID).First(&physicianUser).Error; err != nil {
			continue
		}

		// Obtener datos de la clínica
		_ = r.db.Model(&clinical.ClinicInformation{}).Where("clinic_id = ?", appt.Physician.ClinicID).First(&clinic)

		result = append(result, AppointmentWithNamesDTO{
			MedicalAppointment: appt,
			PatientName:        patientUser.Name,
			PatientGender:      patientUser.Gender,
			PatientEmail:       patientUser.Email,
			PatientPhone:       patientUser.Phone,
			PhysicianName:      physicianUser.Name,
			PhysicianGender:    physicianUser.Gender,
			PhysicianEmail:     physicianUser.Email,
			PhysicianPhone:     physicianUser.Phone,
			ClinicName:         clinic.ClinicName,
			ClinicCity:         clinic.City,
			ClinicAddress:      clinic.Address,
			ClinicId:           clinic.ClinicID,
		})
	}

	return result, nil
}
