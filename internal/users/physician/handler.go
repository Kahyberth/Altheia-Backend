package physician

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{s}
}

func (h *Handler) RegisterPhysician(c *fiber.Ctx) error {
	var physician CreatePhysicianInfo
	if err := c.BodyParser(&physician); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	if err := h.service.RegisterPhysician(physician); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "registered successfully"})
}

func (h *Handler) UpdatePhysician(c *fiber.Ctx) error {
	var physician UpdatePhysicianInfo
	if err := c.BodyParser(&physician); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	id := c.Params("id")
	fmt.Println(id)
	if err := h.service.UpdatePhysician(id, physician); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "updated successfully"})
}

func (h *Handler) SoftDeletePhysician(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.SoftDeletePhysician(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "deleted successfully"})
}

func (h *Handler) GetAllPhysiciansPaginated(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	result, err := h.service.GetAllPhysiciansPaginated(page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func (h *Handler) GetPhysicianById(c *fiber.Ctx) error {
	id := c.Params("id")
	user, _ := h.service.GetPhysicianByID(id)
	return c.JSON(user)
}
