package auth

import (
	"Altheia-Backend/internal/users"
	"fmt"
	"gorm.io/gorm"
)

type Repository interface {
	FindByEmail(email string) (*users.User, error)
	FindByID(id string) (*users.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindByEmail(email string) (*users.User, error) {
	var user users.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *repository) FindByID(id string) (*users.User, error) {
	var user users.User
	fmt.Print("ID del usuario desde repository: ", id)
	err := r.db.Where("id = ?", id).First(&user).Error
	return &user, err
}
