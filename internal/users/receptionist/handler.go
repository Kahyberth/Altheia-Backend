package receptionist

import "github.com/gofiber/fiber/v2"

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{s}
}

func (h *Handler) RegisterReceptionist(c *fiber.Ctx) error {
	var receptionist CreateReceptionistInfo
	if err := c.BodyParser(&receptionist); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	if err := h.service.RegisterReceptionist(receptionist); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "registered successfully"})
}

func (h *Handler) UpdateReceptionist(c *fiber.Ctx) error {
	var receptionist UpdateReceptionistInfo
	if err := c.BodyParser(&receptionist); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	id := c.Params("id")
	if err := h.service.UpdateReceptionist(id, receptionist); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "updated successfully"})
}

func (h *Handler) SoftDeleteReceptionist(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.SoftDelete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "deleted successfully"})
}

func (h *Handler) GetAllReceptionistsPaginated(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	result, err := h.service.GetAllReceptionistPaginated(page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}
