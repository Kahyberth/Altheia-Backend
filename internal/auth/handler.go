package auth

import (
	"Altheia-Backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"time"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{s}
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "cannot parse JSON"})
	}
	if err := h.service.Register(&user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "registered successfully"})
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var data User
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	accessToken, refreshToken, err := h.service.Login(data.Email, data.Password)

	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	accessTokenCookie := fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(1 * time.Hour),
		HTTPOnly: true,
		SameSite: "Lax",
		Path:     "/",
	}

	refreshTokenCookie := fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(72 * time.Hour),
		HTTPOnly: true,
		SameSite: "Lax",
		Path:     "/",
	}

	c.Cookie(&accessTokenCookie)
	c.Cookie(&refreshTokenCookie)

	return c.JSON(fiber.Map{"accessToken": accessToken, "refreshToken": refreshToken})
}

func (h *Handler) Profile(c *fiber.Ctx) error {
	id := c.Locals("user_id").(string)
	user, err := h.service.GetProfile(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}
	return c.JSON(user)
}

func (h *Handler) RefreshTokenH(c *fiber.Ctx) error {
	c.AllParams()
	refreshToken := c.Params("refresh_token")
	accessToken, refreshToken, err := utils.RefreshToken(refreshToken)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	return c.JSON(fiber.Map{"refresh_token": refreshToken, "access_token": accessToken, "message": "refresh token successfully"})
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		Path:     "/",
	}

	c.Cookie(&cookie)
	return c.JSON(fiber.Map{"message": "logout successful"})
}
