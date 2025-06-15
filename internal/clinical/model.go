package clinical

import (
	"Altheia-Backend/internal/users"
	"time"

	"gorm.io/gorm"
)

type MedicalHistory struct {
	ID            string    `gorm:"primaryKey" json:"id"`
	PatientId     string    `json:"patient_id"`
	ConsultReason string    `json:"consult_reason"`
	PersonalInfo  string    `json:"personal_info"`
	FamilyInfo    string    `json:"family_info"`
	Allergies     string    `json:"allergies"`
	Observations  string    `json:"observations"`
	LastUpdate    time.Time `json:"last_update"`

	Patient       users.Patient         `gorm:"foreignKey:PatientId"`
	Consultations []MedicalConsultation `gorm:"foreignKey:MedicalHistoryId"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type MedicalConsultation struct {
	ID               string    `gorm:"primaryKey" json:"id"`
	MedicalHistoryId string    `json:"medical_history_id"`
	PhysicianId      string    `json:"physician_id"`
	ConsultDate      time.Time `json:"consult_date"`
	Symptoms         string    `json:"symptoms"`
	Diagnosis        string    `json:"diagnosis"`
	Treatment        string    `json:"treatment"`
	Notes            string    `json:"notes"`

	MedicalHistory MedicalHistory        `gorm:"foreignKey:MedicalHistoryId"`
	Physician      users.Physician       `gorm:"foreignKey:PhysicianId"`
	Prescriptions  []MedicalPrescription `gorm:"foreignKey:ConsultationId"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type MedicalPrescription struct {
	ID             string    `gorm:"primaryKey" json:"id"`
	ConsultationId string    `json:"consultation_id"`
	Medicine       string    `json:"medicine"`
	Dosage         string    `json:"dosage"`
	Frequency      string    `json:"frequency"`
	Duration       string    `json:"duration"`
	Instructions   string    `json:"instructions"`
	IssuedAt       time.Time `json:"issued_at"`

	Consultation MedicalConsultation `gorm:"foreignKey:ConsultationId"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Clinic struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	Status    bool           `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UserID    string         `gorm:"not null;index" json:"user_id"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Physicians        []users.Physician     `gorm:"foreignKey:ClinicID" json:"physicians,omitempty"`
	Receptionists     []users.Receptionist  `gorm:"foreignKey:ClinicID" json:"receptionists,omitempty"`
	LabTechnicians    []users.LabTechnician `gorm:"foreignKey:ClinicID" json:"lab_technicians,omitempty"`
	Patients          []users.Patient       `gorm:"foreignKey:ClinicID" json:"patients,omitempty"`
	ClinicInformation ClinicInformation     `gorm:"foreignKey:ClinicID;references:ID" json:"clinic_information,omitempty"`
}

type ClinicSchedule struct {
	ID       string `gorm:"primaryKey" json:"id"`
	ClinicID string `gorm:"not null" json:"clinic_id"`
	Day      string `gorm:"not null" json:"day"`
	Open     bool   `json:"open"`
	From     string `json:"from"`
	To       string `json:"to"`
}

type ClinicInformation struct {
	ClinicID          string            `gorm:"primaryKey" json:"clinic_id"`
	ClinicEmail       string            `json:"clinic_email"`
	ClinicName        string            `json:"clinic_name"`
	ClinicPhone       string            `json:"clinic_phone"`
	ClinicDescription string            `json:"clinic_description"`
	ClinicWebsite     string            `json:"clinic_website"`
	EmployeeCount     int               `json:"employee_count"`
	ServicesOffered   []ServicesOffered `gorm:"many2many:clinic_services;" json:"services offered,omitempty"`
	EpsOffered        []EPS             `gorm:"many2many:clinic_eps;" json:"eps offered,omitempty"`
	Photos            []Photo           `gorm:"foreignKey:ClinicID" json:"photos,omitempty"`

	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code,omitempty"`
	Country    string `json:"country"`
}

type EPS struct {
	ID   string `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique" json:"name"`
}

type ServicesOffered struct {
	ID   string `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique" json:"name"`
}

type Photo struct {
	ID       string `gorm:"primaryKey" json:"id"`
	URL      string `json:"url"`
	ClinicID string `json:"clinic_id"`
}

type PersonnelResponse struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Email          string                 `json:"email"`
	Role           string                 `json:"role"`
	Phone          string                 `json:"phone"`
	DocumentNumber string                 `json:"document_number"`
	Status         bool                   `json:"status"`
	Gender         string                 `json:"gender"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	LastLogin      time.Time              `json:"last_login"`
	RoleDetails    map[string]interface{} `json:"role_details"`
}

type ClinicPersonnelResponse struct {
	Personnel []PersonnelResponse `json:"personnel"`
	Count     int                 `json:"count"`
}

type CreateMedicalHistoryDTO struct {
	PatientId     string `json:"patient_id" validate:"required"`
	ConsultReason string `json:"consult_reason"`
	PersonalInfo  string `json:"personal_info"`
	FamilyInfo    string `json:"family_info"`
	Allergies     string `json:"allergies"`
	Observations  string `json:"observations"`
}

type UpdateMedicalHistoryDTO struct {
	ConsultReason string `json:"consult_reason"`
	PersonalInfo  string `json:"personal_info"`
	FamilyInfo    string `json:"family_info"`
	Allergies     string `json:"allergies"`
	Observations  string `json:"observations"`
}

type CreateConsultationDTO struct {
	MedicalHistoryId *string                 `json:"medical_history_id"`
	PatientId        string                  `json:"patient_id" validate:"required"`
	PhysicianId      string                  `json:"physician_id" validate:"required"`
	Symptoms         string                  `json:"symptoms"`
	Diagnosis        string                  `json:"diagnosis"`
	Treatment        string                  `json:"treatment"`
	Notes            string                  `json:"notes"`
	Prescriptions    []CreatePrescriptionDTO `json:"prescriptions,omitempty"`

	// Optional fields to update medical history when creating consultation
	UpdateMedicalHistory bool   `json:"update_medical_history,omitempty"`
	ConsultReason        string `json:"consult_reason,omitempty"`
	PersonalInfo         string `json:"personal_info,omitempty"`
	FamilyInfo           string `json:"family_info,omitempty"`
	Allergies            string `json:"allergies,omitempty"`
	Observations         string `json:"observations,omitempty"`
}

type CreatePrescriptionDTO struct {
	Medicine     string `json:"medicine" validate:"required"`
	Dosage       string `json:"dosage" validate:"required"`
	Frequency    string `json:"frequency" validate:"required"`
	Duration     string `json:"duration" validate:"required"`
	Instructions string `json:"instructions"`
}

type MedicalHistoryResponseDTO struct {
	ID                 string                            `json:"id"`
	PatientId          string                            `json:"patient_id"`
	PatientName        string                            `json:"patient_name"`
	ConsultReason      string                            `json:"consult_reason"`
	PersonalInfo       string                            `json:"personal_info"`
	FamilyInfo         string                            `json:"family_info"`
	Allergies          string                            `json:"allergies"`
	Observations       string                            `json:"observations"`
	LastUpdate         time.Time                         `json:"last_update"`
	CreatedAt          time.Time                         `json:"created_at"`
	UpdatedAt          time.Time                         `json:"updated_at"`
	Consultations      []EnhancedConsultationResponseDTO `json:"consultations"`
	TotalConsultations int                               `json:"total_consultations"`
}

type ConsultationResponseDTO struct {
	MedicalConsultation
	PhysicianName string                    `json:"physician_name"`
	PhysicianInfo PhysicianInfoDTO          `json:"physician_info"`
	Prescriptions []PrescriptionResponseDTO `json:"prescriptions"`
}

type PhysicianInfoDTO struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Email              string `json:"email"`
	Phone              string `json:"phone"`
	DocumentNumber     string `json:"document_number"`
	PhysicianSpecialty string `json:"physician_specialty"`
	LicenseNumber      string `json:"license_number"`
	Gender             string `json:"gender"`
}

type PrescriptionResponseDTO struct {
	ID             string    `json:"id"`
	ConsultationId string    `json:"consultation_id"`
	Medicine       string    `json:"medicine"`
	Dosage         string    `json:"dosage"`
	Frequency      string    `json:"frequency"`
	Duration       string    `json:"duration"`
	Instructions   string    `json:"instructions"`
	IssuedAt       time.Time `json:"issued_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ConsultationMetadata struct {
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	ConsultDate      time.Time `json:"consult_date"`
	DurationMinutes  int       `json:"duration_minutes,omitempty"`
	ConsultationType string    `json:"consultation_type,omitempty"`
}

type EnhancedConsultationResponseDTO struct {
	ID               string                    `json:"id"`
	MedicalHistoryId string                    `json:"medical_history_id"`
	Symptoms         string                    `json:"symptoms"`
	Diagnosis        string                    `json:"diagnosis"`
	Treatment        string                    `json:"treatment"`
	Notes            string                    `json:"notes"`
	PhysicianInfo    PhysicianInfoDTO          `json:"physician_info"`
	Metadata         ConsultationMetadata      `json:"metadata"`
	Prescriptions    []PrescriptionResponseDTO `json:"prescriptions"`
}
