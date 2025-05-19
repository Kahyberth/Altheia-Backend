package clinical

type CreateClinicDTO struct {
	OwnerName           string `json:"owner_name"`
	OwnerEmail          string `json:"owner_email"`
	OwnerPhone          string `json:"owner_phone"`
	OwnerPosition       string `json:"owner_position"`
	OwnerDocumentNumber string `json:"owner_document_number"`
	OwnerGender         string `json:"owner_gender"`

	Name        string `json:"name"`
	Email       string `json:"email"`
	Description string `json:"description"`
	Phone       string `json:"phone"`
	Website     string `json:"website"`

	Address    string `json:"address"`
	Country    string `json:"country"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`

	MemberCount     int      `json:"member_count"`
	ServicesOffered []string `json:"services_offered"`

	AcceptedEPS []string `json:"accepted_eps"`
}

type CreateEpsDto struct {
	Eps []string `json:"eps"`
}

type CreateServicesDto struct {
	ServicesOffered []string `json:"services_offered"`
}
