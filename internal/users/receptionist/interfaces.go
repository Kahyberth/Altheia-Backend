package receptionist

type CreateReceptionistInfo struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Gender         string `json:"gender"`
	Phone          string `json:"phone"`
	DocumentNumber string `json:"document_number"`
	ClinicID       string `json:"clinic_id"`
}
type UpdateReceptionistInfo struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}
