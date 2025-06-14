package auth

import (
	"Altheia-Backend/internal/users"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	FindByEmail(email string) (*users.User, error)
	FindByID(id string) (*users.User, error)
	GetUserWithAllDetails(id string) (*users.User, error)
	ChangePassword(userID string, newHashedPassword string) error
	UpdateLastLogin(userID string) error
	CreateLoginActivity(activity *users.LoginActivity) error
	GetUserLoginActivities(userID string, limit int) ([]users.LoginActivity, error)
	MarkAllSessionsAsInactive(userID string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindByEmail(email string) (*users.User, error) {
	var user users.User
	err := r.db.Preload("Patient").Preload("Physician").Preload("Receptionist").Preload("ClinicOwner").Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *repository) FindByID(id string) (*users.User, error) {
	var user users.User
	fmt.Print("ID del usuario desde repository: ", id)
	err := r.db.Preload("Patient").Preload("Physician").Preload("Receptionist").Preload("ClinicOwner").Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *repository) GetUserWithAllDetails(id string) (*users.User, error) {
	var user users.User
	err := r.db.Preload("Patient").Preload("Physician").Preload("Receptionist").Preload("ClinicOwner").Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *repository) ChangePassword(userID string, newHashedPassword string) error {
	return r.db.Model(&users.User{}).Where("id = ?", userID).Update("password", newHashedPassword).Error
}

func (r *repository) UpdateLastLogin(userID string) error {
	return r.db.Model(&users.User{}).Where("id = ?", userID).Update("last_login", time.Now()).Error
}

func (r *repository) CreateLoginActivity(activity *users.LoginActivity) error {
	return r.db.Create(activity).Error
}

func (r *repository) GetUserLoginActivities(userID string, limit int) ([]users.LoginActivity, error) {
	var activities []users.LoginActivity
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Limit(limit).Find(&activities).Error
	return activities, err
}

func (r *repository) MarkAllSessionsAsInactive(userID string) error {
	return r.db.Model(&users.LoginActivity{}).Where("user_id = ?", userID).Update("is_current_session", false).Error
}
