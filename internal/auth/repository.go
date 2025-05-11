package auth

import (
	"Altheia-Backend/internal/users"
	"fmt"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *users.User) error
	FindByEmail(email string) (*users.User, error)
	FindByID(id string) (*users.User, error)
}

type UpdateInfo struct {
	Name               string
	Email              string
	Password           string
	Gender             string
	PhysicianSpecialty string
	Phone              string
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(user *users.User) error {
	return r.db.Create(user).Error
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

func (r *repository) Update(user *users.User) error {
	return r.db.Save(user).Error
}
