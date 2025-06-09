package patient

type CreatePatientInfo struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Gender         string `json:"gender"`
	Phone          string `json:"phone"`
	DocumentNumber string `json:"document_number"`
	DateOfBirth    string `json:"date_of_birth"`
	Address        string `json:"address"`
	Eps            string `json:"eps"`
	BloodType      string `json:"blood_type"`
}

type UpdatePatientInfo struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Eps      string `json:"eps"`
}

type GetPatientWithName struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	DateOfBirth string `json:"date_of_birth"`
	Address     string `json:"address"`
	Eps         string `json:"eps"`
	BloodType   string `json:"blood_type"`
	Status      bool   `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
