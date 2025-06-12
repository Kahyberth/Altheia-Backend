package auth

import (
	"Altheia-Backend/internal/users"
	"Altheia-Backend/pkg/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{s}
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var data users.User
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	var user UserInfo

	user, accessToken, refreshToken, err := h.service.Login(data.Email, data.Password)

	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	accessTokenCookie := fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(1 * time.Hour),
		HTTPOnly: true,
		SameSite: "Lax",
		Secure:   false,
		Path:     "/",
	}

	refreshTokenCookie := fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(72 * time.Hour),
		HTTPOnly: true,
		SameSite: "Lax",
		Secure:   false,
		Path:     "/",
	}

	c.Cookie(&accessTokenCookie)
	c.Cookie(&refreshTokenCookie)

	return c.JSON(fiber.Map{"accessToken": accessToken, "refreshToken": refreshToken, "user": user})
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
		SameSite: "Lax",
		Secure:   false,
		Path:     "/",
	}

	c.Cookie(&cookie)
	return c.JSON(fiber.Map{"message": "logout successful"})
}

func (h *Handler) VerifyToken(c *fiber.Ctx) error {

	fmt.Print("Entro")

	token := c.Cookies("access_token")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token no proporcionado",
		})
	}
	data, token, err := h.service.verifyToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token inválido o expirado",
		})
	}

	fmt.Println("Token verificado", token)

	return c.JSON(fiber.Map{
		"isValid":  true,
		"userInfo": data,
		"token":    token,
	})
}
