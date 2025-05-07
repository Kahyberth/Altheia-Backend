package main

import (
	"Altheia-Backend/config"
	"Altheia-Backend/internal/auth"
	"Altheia-Backend/internal/db"
	"Altheia-Backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func main() {
	database := db.GetDB()
	err := database.AutoMigrate(&auth.User{})
	if err != nil {
		return
	}

	authRepo := auth.NewRepository(database)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	app := fiber.New()

	authGroup := app.Group("/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/refresh-token/:refresh_token", authHandler.RefreshTokenH)

	profile := app.Group("/profile")
	profile.Use(middleware.JWTProtected())
	profile.Get("/", authHandler.Profile)

	app.Listen(":" + config.GetEnv("PORT"))

}
