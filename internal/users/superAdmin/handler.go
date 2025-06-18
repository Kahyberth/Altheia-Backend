package superAdmin

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{s}
}

func (h *Handler) RegisterSuperAdmin(c *fiber.Ctx) error {
	var data CreateSuperAdminInfo
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if err := h.service.RegisterSuperAdmin(data); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Super Admin registered successfully",
	})
}

func (h *Handler) UpdateSuperAdmin(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	var data UpdateSuperAdminInfo
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if err := h.service.UpdateSuperAdmin(userID, data); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Super Admin updated successfully",
	})
}

func (h *Handler) GetSuperAdminByID(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	superAdmin, err := h.service.GetSuperAdminByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(superAdmin)
}

func (h *Handler) GetAllSuperAdminsPaginated(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	result, err := h.service.GetAllSuperAdminsPaginated(page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}

func (h *Handler) SoftDeleteSuperAdmin(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	if err := h.service.SoftDelete(userID); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Super Admin deleted successfully",
	})
}

func (h *Handler) GetAllSystemData(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "This endpoint provides access to all system data",
		"note":    "Only super-admin users can access this endpoint",
		"data": fiber.Map{
			"users":           "All users in the system",
			"clinics":         "All clinics in the system",
			"appointments":    "All appointments in the system",
			"medical_records": "All medical records in the system",
			"system_logs":     "All system logs",
			"analytics":       "System analytics and reports",
		},
	})
}

func (h *Handler) GetDeactivatedUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	result, err := h.service.GetDeactivatedUsersPaginated(page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to retrieve deactivated users",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Deactivated users retrieved successfully",
		"data":    result,
	})
}

func (h *Handler) GetClinicOwners(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	result, err := h.service.GetClinicOwnersPaginated(page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to retrieve clinic owners",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Clinic owners retrieved successfully",
		"data":    result,
	})
}
