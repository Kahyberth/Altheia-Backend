package auth

import (
	"Altheia-Backend/internal/users"
	"Altheia-Backend/pkg/utils"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
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

	userAgent := c.Get("User-Agent")
	ipAddress := utils.GetClientIP(
		c.Get("X-Forwarded-For"),
		c.Get("X-Real-IP"),
		c.IP(),
	)

	var user UserInfo

	user, accessToken, refreshToken, err := h.service.LoginWithActivity(data.Email, data.Password, userAgent, ipAddress)

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

func (h *Handler) GetUserDetails(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "User ID is required"})
	}

	userDetails, err := h.service.GetUserDetails(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(userDetails)
}

func (h *Handler) ChangePassword(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var request ChangePasswordRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.ChangePassword(userIDStr, request); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Password changed successfully"})
}

func (h *Handler) GetUserLoginActivities(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "User ID is required"})
	}

	limit := 10
	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil {
			if parsedLimit > 0 && parsedLimit <= 50 {
				limit = parsedLimit
			}
		}
	}

	activities, err := h.service.GetUserLoginActivities(userID, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve login activities"})
	}

	return c.JSON(fiber.Map{
		"activities": activities,
		"count":      len(activities),
	})
}
