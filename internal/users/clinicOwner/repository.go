package clinicOwner

import (
	"Altheia-Backend/internal/users"
	"gorm.io/gorm"
	"time"
)

type Repository interface {
	Create(createClinicOwnerDto CreateClinicOwnerDto) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(createClinicOwnerDto CreateClinicOwnerDto) error {

	newUser := users.User{
		ID:             "",
		Name:           createClinicOwnerDto.Name,
		Email:          createClinicOwnerDto.Email,
		Password:       createClinicOwnerDto.Password,
		Rol:            "staff",
		Phone:          createClinicOwnerDto.Phone,
		DocumentNumber: createClinicOwnerDto.DocumentNumber,
		Status:         true,
		Gender:         createClinicOwnerDto.Gender,
		CreatedAt:      time.Time{},
		UpdatedAt:      time.Time{},
		DeletedAt:      gorm.DeletedAt{},
		LastLogin:      time.Time{},
		Patient:        users.Patient{},
		Physician:      users.Physician{},
		Receptionist:   users.Receptionist{},
		ClinicOwner:    users.ClinicOwner{},
	}

	r.db.Create(&newUser)
	return nil
}
