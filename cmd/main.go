package main

import (
	"Altheia-Backend/config"
	"Altheia-Backend/internal/auth"
	"Altheia-Backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadEnv()
	db := config.ConnectDB()
	err := db.AutoMigrate(&auth.User{})
	if err != nil {
		return
	}

	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	app := fiber.New()

	authGroup := app.Group("/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)

	profile := app.Group("/profile")
	profile.Use(middleware.JWTProtected())
	profile.Get("/", authHandler.Profile)

	app.Listen(":" + config.GetEnv("PORT"))

}
