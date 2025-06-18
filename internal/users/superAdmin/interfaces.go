package superAdmin

import "time"

type CreateSuperAdminInfo struct {
	Name           string `json:"name" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=8"`
	Phone          string `json:"phone" validate:"required"`
	DocumentNumber string `json:"document_number" validate:"required"`
	Gender         string `json:"gender" validate:"required,oneof=male female other"`
	Permissions    string `json:"permissions,omitempty"`
}

type UpdateSuperAdminInfo struct {
	Name           string `json:"name,omitempty"`
	Email          string `json:"email,omitempty"`
	Phone          string `json:"phone,omitempty"`
	DocumentNumber string `json:"document_number,omitempty"`
	Gender         string `json:"gender,omitempty"`
	Permissions    string `json:"permissions,omitempty"`
}

type SuperAdminResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	DocumentNumber string    `json:"document_number"`
	Gender         string    `json:"gender"`
	Status         bool      `json:"status"`
	Permissions    string    `json:"permissions"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	LastLogin      time.Time `json:"last_login"`
}

type DeactivatedUserResponse struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Email          string                 `json:"email"`
	Role           string                 `json:"role"`
	Phone          string                 `json:"phone"`
	DocumentNumber string                 `json:"document_number"`
	Gender         string                 `json:"gender"`
	Status         bool                   `json:"status"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	LastLogin      time.Time              `json:"last_login"`
	RoleDetails    map[string]interface{} `json:"role_details"`
}

type ClinicOwnerResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	DocumentNumber string    `json:"document_number"`
	Gender         string    `json:"gender"`
	Status         bool      `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	LastLogin      time.Time `json:"last_login"`
	ClinicOwnerID  string    `json:"clinic_owner_id"`
	ClinicID       string    `json:"clinic_id"`
	OwnerStatus    bool      `json:"owner_status"`
	OwnerCreatedAt time.Time `json:"owner_created_at"`
}
