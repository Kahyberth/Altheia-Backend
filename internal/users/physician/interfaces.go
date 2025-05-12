package physician

type CreatePhysicianInfo struct {
	Name                string `json:"name"`
	Email               string `json:"email"`
	Password            string `json:"password"`
	Gender              string `json:"gender"`
	Phone               string `json:"phone"`
	DocumentNumber      string `json:"document_number"`
	PhysicianSpeciality string `json:"physician_specialty"`
	LicenseNumber       string `json:"license_number"`
}

type UpdatePhysicianInfo struct {
	Name               string `json:"name"`
	Password           string `json:"password"`
	Phone              string `json:"phone"`
	PhysicianSpecialty string `json:"physician_specialty"`
}
