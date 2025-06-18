package middleware

import (
	"Altheia-Backend/internal/auth"
	"Altheia-Backend/internal/db"
	"Altheia-Backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func SuperAdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}

func SuperAdminOrOwner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("access_token")

		if token == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "unauthorized: missing authentication token",
			})
		}

		userID, err := utils.ValidateJWT(token)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		database := db.GetDB()
		authRepo := auth.NewRepository(database)
		user, err := authRepo.FindByID(userID)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "user not found",
			})
		}

		if user.Rol != "super-admin" && user.Rol != "owner" {
			return c.Status(403).JSON(fiber.Map{
				"error": "access denied: super-admin or owner privileges required",
			})
		}

		if !user.Status {
			return c.Status(403).JSON(fiber.Map{
				"error": "account is deactivated",
			})
		}

		c.Locals("user_id", userID)
		c.Locals("user_role", user.Rol)

		return c.Next()
	}
}
