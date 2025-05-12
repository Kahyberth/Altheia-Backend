package clinical

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{s}
}

func (h *Handler) CreateClinical(c *fiber.Ctx) error {

	var createClinicDto CreateClinicDTO

	if err := c.BodyParser(&createClinicDto); err != nil {
		return err
	}

	createClinicError := h.service.CreateClinical(createClinicDto)

	if createClinicError != nil {
		return createClinicError
	}
	return nil
}

func (h *Handler) CreateEps(c *fiber.Ctx) error {
	var createEpsDto CreateEpsDto
	if err := c.BodyParser(&createEpsDto); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	createEpsError := h.service.CreateEps(createEpsDto)
	if createEpsError != nil {
		return createEpsError
	}
	return nil
}

func (h *Handler) GetAllEps(c *fiber.Ctx) error {
	var eps []EPS
	page, errPage := strconv.ParseInt(c.Query("page"), 10, 16)
	pageSize, errSize := strconv.ParseInt(c.Query("size"), 10, 16)

	if errPage != nil {
		return errPage
	}

	if errSize != nil {
		return errSize
	}

	eps, epsError := h.service.GetAllEps(int(page), int(pageSize))
	if epsError != nil {
		return epsError
	}

	return c.JSON(eps)
}
