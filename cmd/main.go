package main

import (
	"Altheia-Backend/config"
	"Altheia-Backend/internal/auth"
	"Altheia-Backend/internal/clinical"
	"Altheia-Backend/internal/clinical/appointments"
	"Altheia-Backend/internal/db"
	"Altheia-Backend/internal/middleware"
	"Altheia-Backend/internal/users"
	"Altheia-Backend/internal/users/clinicOwner"
	"Altheia-Backend/internal/users/patient"
	"Altheia-Backend/internal/users/physician"
	"Altheia-Backend/internal/users/receptionist"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database := db.GetDB()
	err := database.AutoMigrate(

		&clinical.Clinic{},
		&clinical.ClinicInformation{},
		&clinical.ClinicSchedule{},
		&clinical.ServicesOffered{},
		&clinical.Photo{},
		&clinical.EPS{},

		&users.User{},
		&users.Patient{},
		&users.Physician{},
		&users.Receptionist{},
		&users.ClinicOwner{},
		&users.LabTechnician{},
		&users.LoginActivity{},

		&clinical.MedicalHistory{},
		&clinical.MedicalConsultation{},
		&appointments.MedicalAppointment{},
		&clinical.MedicalPrescription{},
	)
	if err != nil {
		return
	}

	client := os.Getenv("CLIENT")

	authRepo := auth.NewRepository(database)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	// Patient handler
	patientRepo := patient.NewRepository(database)
	patientService := patient.NewService(patientRepo)
	patientHandler := patient.NewHandler(patientService)

	// Physician handler
	physicianRepo := physician.NewRepository(database)
	physicianService := physician.NewService(physicianRepo)
	physicianHandler := physician.NewHandler(physicianService)

	//Create Clinic handler
	clinicRepo := clinical.NewRepository(database)
	clinicService := clinical.NewService(clinicRepo)
	clinicHandler := clinical.NewHandler(clinicService)

	//Clinic Owner handler
	clinicOwnerRepo := clinicOwner.NewRepository(database)
	clinicOwnerService := clinicOwner.NewService(clinicOwnerRepo)
	clinicOwnerHandler := clinicOwner.NewHandler(clinicOwnerService)

	// Receptionist handler
	receptionistRepo := receptionist.NewRepository(database)
	receptionistService := receptionist.NewService(receptionistRepo)
	receptionistHandler := receptionist.NewHandler(receptionistService)

	// Appointment handler
	appointmentRepo := appointments.NewRepository(database)
	appointmentService := appointments.NewService(appointmentRepo)
	appointmentHandler := appointments.NewHandler(appointmentService)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     client,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Routes

	// Physician routes
	physicianGroup := app.Group("/physician")
	//physicianGroup.Use(middleware.JWTProtected())
	physicianGroup.Post("/register", physicianHandler.RegisterPhysician)
	physicianGroup.Patch("/update/:id", physicianHandler.UpdatePhysician)
	physicianGroup.Get("/getAllPaginated/", physicianHandler.GetAllPhysiciansPaginated)
	physicianGroup.Get("/getAll/", physicianHandler.GetAllPhysicians)
	physicianGroup.Get("/:id", physicianHandler.GetPhysicianById)

	//Clinical routes
	clinicGroup := app.Group("/clinic")
	clinicGroup.Post("/register", clinicHandler.CreateClinical)
	clinicGroup.Post("/create-eps", clinicHandler.CreateEps)
	clinicGroup.Get("/get-eps", clinicHandler.GetAllEps)
	clinicGroup.Post("/create-services", clinicHandler.CreateServices)
	clinicGroup.Get("/get-services", clinicHandler.GetAllServices)
	clinicGroup.Get("/by-owner/:ownerId", clinicHandler.GetClinicByOwnerID)
	clinicGroup.Get("/:clinicId", clinicHandler.GetClinicByID)
	clinicGroup.Post("/assign-services", clinicHandler.AssignServicesToClinic)
	clinicGroup.Get("/by-eps/:epsId", clinicHandler.GetClinicsByEps)
	clinicGroup.Get("/personnel/:clinicId", clinicHandler.GetClinicPersonnel)
	clinicGroup.Get("/patients/:clinicId", patientHandler.GetPatientByClinicId)

	// Medical History routes
	medicalHistoryGroup := app.Group("/medical-history")
	//medicalHistoryGroup.Use(middleware.JWTProtected())
	medicalHistoryGroup.Get("/patient/:patientId", clinicHandler.GetMedicalHistoryByPatientID)
	medicalHistoryGroup.Post("/create", clinicHandler.CreateMedicalHistory)
	medicalHistoryGroup.Post("/consultation/create", clinicHandler.CreateConsultation)
	medicalHistoryGroup.Put("/update/:historyId", clinicHandler.UpdateMedicalHistory)

	//Patient routes
	patientGroup := app.Group("/patient")
	patientGroup.Post("/register", patientHandler.RegisterPatient)
	patientGroup.Get("/getAllPaginated", patientHandler.GetAllPatientsPaginated)
	patientGroup.Get("/getAll", patientHandler.GetAllPatients)
	patientGroup.Get("/getByClinicId/:clinicId", patientHandler.GetPatientByClinicId)
	patientGroup.Patch("/update/:id", patientHandler.UpdatePatient)
	patientGroup.Post("/delete/:id", patientHandler.SoftDeletePatient)

	// Receptionist routes
	receptionistGroup := app.Group("/receptionist")
	receptionistGroup.Post("/register", receptionistHandler.RegisterReceptionist)
	receptionistGroup.Patch("/update/:id", receptionistHandler.UpdateReceptionist)
	receptionistGroup.Get("/getAll", receptionistHandler.GetAllReceptionistsPaginated)

	// Lab Technician routes

	// Clinic Owner Routes
	clinicOwnerGroup := app.Group("/clinic-owner")
	clinicOwnerGroup.Post("/register", clinicOwnerHandler.CreateClinicOwner)

	// Appointments routes
	appointmentGroup := app.Group("/appointments")
	//appointmentGroup.Use(middleware.JWTProtected())
	appointmentGroup.Post("/create", appointmentHandler.CreateAppointment)
	appointmentGroup.Get("/getAll", appointmentHandler.GetAllAppointments)
	appointmentGroup.Patch("/updateStatus/:id", appointmentHandler.UpdateAppointmentStatus)
	appointmentGroup.Get("/getAllByMedicId/:id", appointmentHandler.GetAllAppointmentsByMedicId)

	// Auth routes
	authGroup := app.Group("/auth")
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/logout", authHandler.Logout)
	authGroup.Get("/verify-token", authHandler.VerifyToken)
	authGroup.Get("/user/:id", authHandler.GetUserDetails)
	authGroup.Get("/user/:id/login-activities", authHandler.GetUserLoginActivities)
	authGroup.Get("/user/:id/exists", authHandler.CheckUserExists)
	authGroup.Delete("/user/:id", authHandler.DeleteUserCompletely)
	authGroup.Patch("/user/:id/deactivate", authHandler.DeactivateUser)
	authGroup.Patch("/user/:id/reactivate", authHandler.ReactivateUser)

	authGroup.Use(middleware.JWTProtected())
	authGroup.Post("/change-password", authHandler.ChangePassword)
	authGroup.Post("/refresh-token/:refresh_token", authHandler.RefreshTokenH)

	profile := app.Group("/profile")
	profile.Use(middleware.JWTProtected())

	_ = app.Listen(":" + config.GetEnv("PORT"))

}
