package patient

import (
	"Altheia-Backend/internal/users"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *users.Patient) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(user *users.Patient) error {
	r.db.Create(user)
	return nil
}
