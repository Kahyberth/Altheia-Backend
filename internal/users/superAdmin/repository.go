package superAdmin

import (
	"Altheia-Backend/internal/users"

	"gorm.io/gorm"
)

type Repository interface {
	Create(user *users.User) error
	Update(userID string, updateData UpdateSuperAdminInfo) error
	GetByID(userID string) (*users.User, error)
	GetAll(page, limit int) ([]users.User, int64, error)
	SoftDelete(userID string) error
	ValidateUserExists(userID string) error
	GetDeactivatedUsers(page, limit int) ([]users.User, int64, error)
	GetClinicOwners(page, limit int) ([]users.User, int64, error)
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

func (r *repository) Update(userID string, updateData UpdateSuperAdminInfo) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	userUpdate := map[string]interface{}{}
	if updateData.Name != "" {
		userUpdate["name"] = updateData.Name
	}
	if updateData.Email != "" {
		userUpdate["email"] = updateData.Email
	}
	if updateData.Phone != "" {
		userUpdate["phone"] = updateData.Phone
	}
	if updateData.DocumentNumber != "" {
		userUpdate["document_number"] = updateData.DocumentNumber
	}
	if updateData.Gender != "" {
		userUpdate["gender"] = updateData.Gender
	}

	if len(userUpdate) > 0 {
		if err := tx.Model(&users.User{}).Where("id = ?", userID).Updates(userUpdate).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	superAdminUpdate := map[string]interface{}{}
	if updateData.Permissions != "" {
		superAdminUpdate["permissions"] = updateData.Permissions
	}

	if len(superAdminUpdate) > 0 {
		if err := tx.Model(&users.SuperAdmin{}).Where("user_id = ?", userID).Updates(superAdminUpdate).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *repository) GetByID(userID string) (*users.User, error) {
	var user users.User
	err := r.db.Preload("SuperAdmin").Where("id = ?", userID).First(&user).Error
	return &user, err
}

func (r *repository) GetAll(page, limit int) ([]users.User, int64, error) {
	var superAdmins []users.User
	var total int64

	offset := (page - 1) * limit

	r.db.Model(&users.User{}).Where("rol = ?", "super-admin").Count(&total)

	err := r.db.Preload("SuperAdmin").
		Where("rol = ?", "super-admin").
		Offset(offset).
		Limit(limit).
		Find(&superAdmins).Error

	return superAdmins, total, err
}

func (r *repository) SoftDelete(userID string) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Delete(&users.SuperAdmin{}, "user_id = ?", userID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&users.User{}, "id = ?", userID).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *repository) ValidateUserExists(userID string) error {
	var user users.User
	return r.db.Where("id = ? AND rol = ?", userID, "super-admin").First(&user).Error
}

func (r *repository) GetDeactivatedUsers(page, limit int) ([]users.User, int64, error) {
	var deactivatedUsers []users.User
	var total int64

	offset := (page - 1) * limit

	r.db.Model(&users.User{}).Where("status = ?", false).Count(&total)

	err := r.db.Preload("Patient").
		Preload("Physician").
		Preload("Receptionist").
		Preload("ClinicOwner").
		Preload("LabTechnician").
		Preload("SuperAdmin").
		Where("status = ?", false).
		Offset(offset).
		Limit(limit).
		Order("updated_at DESC").
		Find(&deactivatedUsers).Error

	return deactivatedUsers, total, err
}

func (r *repository) GetClinicOwners(page, limit int) ([]users.User, int64, error) {
	var clinicOwners []users.User
	var total int64

	offset := (page - 1) * limit

	r.db.Model(&users.User{}).Where("rol = ?", "owner").Count(&total)

	err := r.db.Preload("ClinicOwner").
		Where("rol = ?", "owner").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&clinicOwners).Error

	return clinicOwners, total, err
}
