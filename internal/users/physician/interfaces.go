package physician

import "time"

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

type ResultPhysicians struct {
	UserId             string    `json:"user_id"`
	PhysicianID        string    `json:"physician_id"`
	Name               string    `json:"name"`
	Email              string    `json:"email"`
	Rol                string    `json:"rol"`
	UserStatus         bool      `json:"user_status"`
	Gender             string    `json:"gender"`
	LastLogin          time.Time `json:"last_login"`
	PhysicianSpecialty string    `json:"physician_specialty"`
	PhysicianStatus    bool      `json:"physician_status"`
}
