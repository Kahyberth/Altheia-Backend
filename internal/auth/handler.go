package auth

import (
	"Altheia-Backend/pkg/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

type RegisterPatient struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Role        string `json:"role"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
	Phone       string `json:"phone"`
}

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{s}
}

func (h *Handler) RegisterPatient(c *fiber.Ctx) error {
	var user RegisterPatient
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	newUser := User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Rol:      user.Role,
		Gender:   user.Gender,
		Phone:    user.Phone,
		Status:   true,
	}

	if err := h.service.RegisterPatient(&newUser); err != nil {
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
		SameSite: "Lax",
		Secure:   false,
		Path:     "/",
	}

	c.Cookie(&cookie)
	return c.JSON(fiber.Map{"message": "logout successful"})
}

func (h *Handler) VerifyToken(c *fiber.Ctx) error {
	fmt.Print("Se esta verificando...")
	token := c.Cookies("access_token")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token no proporcionado",
		})
	}
	data, err := h.service.verifyToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token inválido o expirado",
		})
	}
	fmt.Print("Se verifico el token!")
	return c.JSON(fiber.Map{
		"isValid":      true,
		"access_token": data,
	})
}
