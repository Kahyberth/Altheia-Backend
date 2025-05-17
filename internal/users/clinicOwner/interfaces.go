package clinicOwner

type CreateClinicOwnerDto struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Gender         string `json:"gender"`
	Phone          string `json:"phone"`
	DocumentNumber string `json:"document_number"`
	DateOfBirth    string `json:"date_of_birth"`
	Address        string `json:"address"`
}
