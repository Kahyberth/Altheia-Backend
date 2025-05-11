package main

import (
	"Altheia-Backend/config"
	"Altheia-Backend/internal/auth"
	"Altheia-Backend/internal/db"
	"Altheia-Backend/internal/middleware"
	"Altheia-Backend/internal/users"
	"Altheia-Backend/internal/users/physician"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"os"
)

func main() {
	database := db.GetDB()
	err := database.AutoMigrate(
		&users.User{},
		&users.Patient{},
		&users.Physician{},
	)
	if err != nil {
		return
	}

	client := os.Getenv("CLIENT")

	authRepo := auth.NewRepository(database)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	// Patient handler
	patientRepo := physician.NewRepository(database)
	patientService := physician.NewService(patientRepo)
	patientHandler := physician.NewHandler(patientService)

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
	physicianGroup.Post("/register", patientHandler.RegisterPhysician)
	//physicianGroup.Put("/update/:id", patientHandler.UpdatePhysician)
	//physicianGroup.Get("/get/:id", patientHandler.GetPhysicianByID)

	// Auth routes
	authGroup := app.Group("/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/logout", authHandler.Logout)
	authGroup.Get("/verify-token", authHandler.VerifyToken)
	authGroup.Use(middleware.JWTProtected())
	authGroup.Post("/refresh-token/:refresh_token", authHandler.RefreshTokenH)

	profile := app.Group("/profile")
	profile.Use(middleware.JWTProtected())
	profile.Get("/", authHandler.Profile)

	app.Listen(":" + config.GetEnv("PORT"))

}
