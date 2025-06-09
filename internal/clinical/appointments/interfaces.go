package appointments

type AppointmentStatus string

const (
	AppointmentStatusPending   AppointmentStatus = "pending"
	AppointmentStatusConfirmed AppointmentStatus = "confirmed"
	AppointmentStatusCancelled AppointmentStatus = "cancelled"
	AppointmentStatusCompleted AppointmentStatus = "completed"
	AppointmentStatusNoShow    AppointmentStatus = "no_show"
)

type CreateAppointmentDTO struct {
	PatientId   string `json:"patient_id"`
	PhysicianId string `json:"physician_id"`
	ClinicId    string `json:"clinic_id"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	Status      string `json:"status"`
	Reason      string `json:"reason"`
}

type AppointmentWithNamesDTO struct {
	MedicalAppointment
	PatientName     string `json:"patient_name"`
	PatientGender   string `json:"patient_gender"`
	PatientEmail    string `json:"patient_email"`
	PatientPhone    string `json:"patient_phone"`
	PhysicianName   string `json:"physician_name"`
	PhysicianGender string `json:"physician_gender"`
	PhysicianEmail  string `json:"physician_email"`
	PhysicianPhone  string `json:"physician_phone"`
}
