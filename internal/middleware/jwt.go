package middleware

import (
	"Altheia-Backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		userID, err := utils.ValidateJWT(token)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}
		c.Locals("user_id", userID)
		return c.Next()
	}
}
