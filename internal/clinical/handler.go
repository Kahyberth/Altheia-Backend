package clinical

import (
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

func (h *Handler) GetClinicsByEps(c *fiber.Ctx) error {
	epsID := c.Params("epsId")
	if epsID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "eps ID is required",
		})
	}

	page, errPage := strconv.Atoi(c.Query("page", "1"))
	size, errSize := strconv.Atoi(c.Query("size", "10"))

	if errPage != nil || errSize != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid pagination parameters",
		})
	}

	clinics, err := h.service.GetClinicsByEps(epsID, page, size)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(clinics)
}

func (h *Handler) GetClinicPersonnel(c *fiber.Ctx) error {
	clinicID := c.Params("clinicId")
	if clinicID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "clinic ID is required",
		})
	}

	personnel, err := h.service.GetClinicPersonnel(clinicID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	type BasicUser struct {
		ID             string    `json:"id"`
		Name           string    `json:"name"`
		Email          string    `json:"email"`
		Password       string    `json:"password"`
		Rol            string    `json:"rol"`
		Phone          string    `json:"phone"`
		DocumentNumber string    `json:"document_number"`
		Status         bool      `json:"status"`
		Gender         string    `json:"gender"`
		CreatedAt      time.Time `json:"createdAt"`
		UpdatedAt      time.Time `json:"updatedAt"`
		LastLogin      time.Time `json:"lastLogin"`
	}

	var response []BasicUser
	for _, u := range personnel {
		response = append(response, BasicUser{
			ID:             u.ID,
			Name:           u.Name,
			Email:          u.Email,
			Password:       u.Password,
			Rol:            u.Rol,
			Phone:          u.Phone,
			DocumentNumber: u.DocumentNumber,
			Status:         u.Status,
			Gender:         u.Gender,
			CreatedAt:      u.CreatedAt,
			UpdatedAt:      u.UpdatedAt,
			LastLogin:      u.LastLogin,
		})
	}

	return c.JSON(response)
}
