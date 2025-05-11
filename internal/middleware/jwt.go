package middleware

import (
	"Altheia-Backend/pkg/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {

		token := c.Cookies("access_token")

		fmt.Print("Cookies: ", token)

		if token == "" {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized: missing cookie"})
		}

		userID, err := utils.ValidateJWT(token)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
		}

		c.Locals("user_id", userID)
		return c.Next()
	}
}
