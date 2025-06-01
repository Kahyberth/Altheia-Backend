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

	Physicians        []users.Physician    `gorm:"foreignKey:ClinicID" json:"physicians,omitempty"`
	Receptionists     []users.Receptionist `gorm:"foreignKey:ClinicID" json:"receptionists,omitempty"`
	ClinicInformation ClinicInformation    `gorm:"foreignKey:ClinicID;references:ID" json:"clinic_information,omitempty"`
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
