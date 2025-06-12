package clinical

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

func (h *Handler) CreateServices(c *fiber.Ctx) error {
	var createService CreateServicesDto
	if err := c.BodyParser(&createService); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	createEpsError := h.service.CreateServicesOffered(createService)
	if createEpsError != nil {
		return createEpsError
	}
	return nil
}

func (h *Handler) GetAllServices(c *fiber.Ctx) error {
	var servicesOffered []ServicesOffered
	page, errPage := strconv.ParseInt(c.Query("page"), 10, 16)
	pageSize, errSize := strconv.ParseInt(c.Query("size"), 10, 16)

	if errPage != nil {
		return errPage
	}

	if errSize != nil {
		return errSize
	}

	servicesOffered, epsError := h.service.GetAllServicesOffered(int(page), int(pageSize))
	if epsError != nil {
		return epsError
	}

	return c.JSON(servicesOffered)
}

func (h *Handler) GetClinicByOwnerID(c *fiber.Ctx) error {
	ownerID := c.Params("ownerId")

	if ownerID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "owner ID is required",
		})
	}

	clinicInfo, err := h.service.GetClinicByOwnerID(ownerID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(clinicInfo)
}

func (h *Handler) AssignServicesToClinic(c *fiber.Ctx) error {
	var dto AssignServicesClinicDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if dto.ClinicID == "" || len(dto.Services) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "clinic_id and services are required",
		})
	}

	if err := h.service.AssignServicesToClinic(dto); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "services assigned to clinic successfully",
	})
}
