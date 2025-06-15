package auth

import (
	"Altheia-Backend/internal/clinical"
	"Altheia-Backend/internal/clinical/appointments"
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
	DeleteUserCompletely(userID string) error
	DeactivateUser(userID string) error
	ReactivateUser(userID string) error
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

func (r *repository) DeleteUserCompletely(userID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		var user users.User
		if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
			return fmt.Errorf("user not found: %v", err)
		}

		fmt.Printf("Deleting user: %s with role: %s\n", userID, user.Rol)

		switch user.Rol {
		case "patient":
			var patient users.Patient
			if err := tx.Where("user_id = ?", userID).First(&patient).Error; err == nil {
				fmt.Printf("Found patient with ID: %s\n", patient.ID)

				result := tx.Exec(`
					DELETE FROM medical_prescriptions 
					WHERE consultation_id IN (
						SELECT id FROM medical_consultations 
						WHERE medical_history_id IN (
							SELECT id FROM medical_histories WHERE patient_id = ?
						)
					)
				`, patient.ID)
				fmt.Printf("Deleted %d prescriptions\n", result.RowsAffected)

				result = tx.Exec(`
					DELETE FROM medical_consultations 
					WHERE medical_history_id IN (
						SELECT id FROM medical_histories WHERE patient_id = ?
					)
				`, patient.ID)
				fmt.Printf("Deleted %d consultations\n", result.RowsAffected)

				result = tx.Where("patient_id = ?", patient.ID).Delete(&clinical.MedicalHistory{})
				if result.Error != nil {
					return fmt.Errorf("failed to delete medical histories: %v", result.Error)
				}
				fmt.Printf("Deleted %d medical histories\n", result.RowsAffected)

				result = tx.Where("patient_id = ?", patient.ID).Delete(&appointments.MedicalAppointment{})
				if result.Error != nil {
					return fmt.Errorf("failed to delete patient appointments: %v", result.Error)
				}
				fmt.Printf("Deleted %d patient appointments\n", result.RowsAffected)
			}

			result := tx.Where("user_id = ?", userID).Delete(&users.Patient{})
			if result.Error != nil {
				return fmt.Errorf("failed to delete patient data: %v", result.Error)
			}
			fmt.Printf("Deleted patient record, rows affected: %d\n", result.RowsAffected)

		case "physician":
			var physician users.Physician
			if err := tx.Where("user_id = ?", userID).First(&physician).Error; err == nil {
				fmt.Printf("Found physician with ID: %s\n", physician.ID)

				result := tx.Exec(`
					DELETE FROM medical_prescriptions 
					WHERE consultation_id IN (
						SELECT id FROM medical_consultations WHERE physician_id = ?
					)
				`, physician.ID)
				fmt.Printf("Deleted %d physician prescriptions\n", result.RowsAffected)

				result = tx.Where("physician_id = ?", physician.ID).Delete(&clinical.MedicalConsultation{})
				if result.Error != nil {
					return fmt.Errorf("failed to delete physician consultations: %v", result.Error)
				}
				fmt.Printf("Deleted %d physician consultations\n", result.RowsAffected)

				result = tx.Where("physician_id = ?", physician.ID).Delete(&appointments.MedicalAppointment{})
				if result.Error != nil {
					return fmt.Errorf("failed to delete physician appointments: %v", result.Error)
				}
				fmt.Printf("Deleted %d physician appointments\n", result.RowsAffected)
			}

			result := tx.Where("user_id = ?", userID).Delete(&users.Physician{})
			if result.Error != nil {
				return fmt.Errorf("failed to delete physician data: %v", result.Error)
			}
			fmt.Printf("Deleted physician record, rows affected: %d\n", result.RowsAffected)

		case "receptionist":

			result := tx.Where("user_id = ?", userID).Delete(&users.Receptionist{})
			if result.Error != nil {
				return fmt.Errorf("failed to delete receptionist data: %v", result.Error)
			}
			fmt.Printf("Deleted receptionist record, rows affected: %d\n", result.RowsAffected)

		case "owner":
			var clinicOwner users.ClinicOwner
			if err := tx.Where("user_id = ?", userID).First(&clinicOwner).Error; err == nil {

				var clinicExists bool
				tx.Model(&clinical.Clinic{}).Select("count(*) > 0").Where("id = ?", clinicOwner.ClinicID).Find(&clinicExists)
				if clinicExists {
					return fmt.Errorf("cannot delete clinic owner while clinic still exists. Please delete or transfer clinic ownership first")
				}
			}

			result := tx.Where("user_id = ?", userID).Delete(&users.ClinicOwner{})
			if result.Error != nil {
				return fmt.Errorf("failed to delete clinic owner data: %v", result.Error)
			}
			fmt.Printf("Deleted clinic owner record, rows affected: %d\n", result.RowsAffected)

		case "lab_technician":

			result := tx.Where("user_id = ?", userID).Delete(&users.LabTechnician{})
			if result.Error != nil {
				return fmt.Errorf("failed to delete lab technician data: %v", result.Error)
			}
			fmt.Printf("Deleted lab technician record, rows affected: %d\n", result.RowsAffected)
		}

		result := tx.Where("user_id = ?", userID).Delete(&users.LoginActivity{})
		if result.Error != nil {
			return fmt.Errorf("failed to delete login activities: %v", result.Error)
		}
		fmt.Printf("Deleted %d login activities\n", result.RowsAffected)

		result = tx.Where("id = ?", userID).Delete(&users.User{})
		if result.Error != nil {
			return fmt.Errorf("failed to delete user: %v", result.Error)
		}
		fmt.Printf("Deleted user record, rows affected: %d\n", result.RowsAffected)

		fmt.Printf("Successfully deleted user %s and all related data\n", userID)
		return nil
	})
}

func (r *repository) DeactivateUser(userID string) error {
	return r.db.Model(&users.User{}).Where("id = ?", userID).Update("status", false).Error
}

func (r *repository) ReactivateUser(userID string) error {
	return r.db.Model(&users.User{}).Where("id = ?", userID).Update("status", true).Error
}
