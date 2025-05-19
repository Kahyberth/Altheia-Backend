package main

import (
	"Altheia-Backend/config"
	"Altheia-Backend/internal/auth"
	"Altheia-Backend/internal/clinical"
	"Altheia-Backend/internal/db"
	"Altheia-Backend/internal/middleware"
	"Altheia-Backend/internal/users"
	"Altheia-Backend/internal/users/clinicOwner"
	"Altheia-Backend/internal/users/patient"
	"Altheia-Backend/internal/users/physician"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"os"
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

		&clinical.MedicalHistory{},
		&clinical.MedicalConsultation{},
		&clinical.MedicalAppointment{},
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

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     client,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// Routes

	// Physician routes
	physicianGroup := app.Group("/physician")
	//physicianGroup.Use(middleware.JWTProtected())
	physicianGroup.Post("/register", physicianHandler.RegisterPhysician)
	physicianGroup.Patch("/update/:id", physicianHandler.UpdatePhysician)
	physicianGroup.Get("/getAll/", physicianHandler.GetAllPhysiciansPaginated)
	physicianGroup.Get("/:id", physicianHandler.GetPhysicianById)

	//Clinical routes
	clinicGroup := app.Group("/clinic")
	clinicGroup.Post("/register", clinicHandler.CreateClinical)
	clinicGroup.Post("/create-eps", clinicHandler.CreateEps)
	clinicGroup.Get("/get-eps", clinicHandler.GetAllEps)
	clinicGroup.Post("/create-services", clinicHandler.CreateServices)
	clinicGroup.Get("/get-services", clinicHandler.GetAllServices)

	//Patient routes
	patientGroup := app.Group("/patient")
	patientGroup.Post("/register", patientHandler.RegisterPatient)
	patientGroup.Get("/getAll", patientHandler.GetAllPatientsPaginated)
	patientGroup.Patch("/update/:id", patientHandler.UpdatePatient)
	patientGroup.Post("/delete/:id", patientHandler.SoftDeletePatient)

	// Clinic Owner Routes
	clinicOwnerGroup := app.Group("/clinic-owner")
	clinicOwnerGroup.Post("/register", clinicOwnerHandler.CreateClinicOwner)

	// Auth routes
	authGroup := app.Group("/auth")
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/logout", authHandler.Logout)
	authGroup.Get("/verify-token", authHandler.VerifyToken)
	authGroup.Use(middleware.JWTProtected())
	authGroup.Post("/refresh-token/:refresh_token", authHandler.RefreshTokenH)

	profile := app.Group("/profile")
	profile.Use(middleware.JWTProtected())

	_ = app.Listen(":" + config.GetEnv("PORT"))

}
