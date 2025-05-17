package clinicOwner

import "github.com/gofiber/fiber/v2"

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) CreateClinicOwner(c *fiber.Ctx) error {

	var clinicOwner CreateClinicOwnerDto
	if err := c.BodyParser(&clinicOwner); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	if err := h.service.CreateClinicOwner(clinicOwner); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "registered successfully"})

}
