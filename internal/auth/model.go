package auth

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	Name      string         `json:"name"`
	Email     string         `gorm:"unique" json:"email"`
	Password  string         `json:"password"`
	Rol       string         `json:"rol"`
	Status    bool           `json:"status"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	LastLogin time.Time `json:"lastLogin"`

	Patient   *Patient   `gorm:"foreignKey:UserId" json:"patient,omitempty"`
	Physician *Physician `gorm:"foreignKey:UserId" json:"physician,omitempty"`
}

type Patient struct {
	ID             string `gorm:"primaryKey" json:"id"`
	UserId         string `json:"user_id"`
	DocumentNumber string `json:"document_number"`
	DateOfBirth    string `json:"date_of_birth"`
	Gender         string `json:"gender"`
	Address        string `json:"address"`
	Phone          string `json:"phone"`
	Eps            string `json:"eps"`
	BloodType      string `json:"blood_type"`

	User           *User                `gorm:"foreignKey:UserId" json:"user,omitempty"`
	MedicalHistory []MedicalHistory     `gorm:"foreignKey:PatientId"`
	Appointments   []MedicalAppointment `gorm:"foreignKey:PatientId"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Physician struct {
	ID                 string `gorm:"primaryKey" json:"id"`
	UserId             string `json:"user_id"`
	Gender             string `json:"gender"`
	PhysicianSpecialty string `json:"physician_specialty"`
	LicenseNumber      string `json:"license_number"`
	Phone              string `json:"phone"`
	Status             bool   `json:"status"`

	User          *User                 `gorm:"foreignKey:UserId" json:"user,omitempty"`
	Consultations []MedicalConsultation `gorm:"foreignKey:PhysicianId"`
	Appointments  []MedicalAppointment  `gorm:"foreignKey:PhysicianId"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type MedicalHistory struct {
	ID            string    `gorm:"primaryKey" json:"id"`
	PatientId     string    `json:"patient_id"`
	ConsultReason string    `json:"consult_reason"`
	PersonalInfo  string    `json:"personal_info"`
	FamilyInfo    string    `json:"family_info"`
	Allergies     string    `json:"allergies"`
	Observations  string    `json:"observations"`
	LastUpdate    time.Time `json:"last_update"`

	Patient       Patient               `gorm:"foreignKey:PatientId"`
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
	Physician      Physician             `gorm:"foreignKey:PhysicianId"`
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

	Patient   Patient   `gorm:"foreignKey:PatientId"`
	Physician Physician `gorm:"foreignKey:PhysicianId"`

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
