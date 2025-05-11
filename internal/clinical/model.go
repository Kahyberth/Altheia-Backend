package clinical

import (
	"Altheia-Backend/internal/users/patient"
	"Altheia-Backend/internal/users/physician"
	"gorm.io/gorm"
	"time"
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

	Patient       patient.Patient       `gorm:"foreignKey:PatientId"`
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
	Physician      physician.Physician   `gorm:"foreignKey:PhysicianId"`
	Prescriptions  []MedicalPrescription `gorm:"foreignKey:ConsultationId"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type MedicalAppointment struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	PatientId   string    `json:"patient_id"`
	PhysicianId string    `json:"physician_id"`
	DateTime    time.Time `json:"date_time"`
	Status      string    `json:"status"`
	Reason      string    `json:"reason"`
	Remarks     string    `json:"remarks"`

	Patient   patient.Patient     `gorm:"foreignKey:PatientId"`
	Physician physician.Physician `gorm:"foreignKey:PhysicianId"`

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
