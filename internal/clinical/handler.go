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

func (h *Handler) GetClinicByID(c *fiber.Ctx) error {
	clinicID := c.Params("clinicId")

	if clinicID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "clinic ID is required",
		})
	}

	clinicInfo, err := h.service.GetClinicByID(clinicID)
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

	personnelResponse, err := h.service.GetClinicPersonnel(clinicID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(personnelResponse)
}

func (h *Handler) GetMedicalHistoryByPatientID(c *fiber.Ctx) error {
	patientID := c.Params("patientId")

	if patientID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Patient ID is required",
		})
	}

	comprehensiveResponse, err := h.service.GetMedicalHistoryComprehensive(patientID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(comprehensiveResponse)
}

func (h *Handler) CreateMedicalHistory(c *fiber.Ctx) error {
	var dto CreateMedicalHistoryDTO

	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if dto.PatientId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Patient ID is required",
		})
	}

	if len(dto.Prescriptions) > 0 && dto.PhysicianId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Physician ID is required when prescriptions are provided",
		})
	}

	response, err := h.service.CreateMedicalHistoryComprehensive(dto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}

func (h *Handler) CreateConsultation(c *fiber.Ctx) error {
	var dto CreateConsultationDTO

	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if dto.PatientId == "" || dto.PhysicianId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Patient ID and Physician ID are required",
		})
	}

	err := h.service.CreateConsultation(dto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Consultation created successfully",
	})
}

func (h *Handler) UpdateMedicalHistory(c *fiber.Ctx) error {
	historyID := c.Params("historyId")

	if historyID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Medical history ID is required",
		})
	}

	var dto UpdateMedicalHistoryDTO

	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err := h.service.UpdateMedicalHistory(historyID, dto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Medical history updated successfully",
	})
}

func (h *Handler) GetClinicMedicalHistoriesPaginated(c *fiber.Ctx) error {
	clinicID := c.Params("clinicId")
	if clinicID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Clinic ID is required",
		})
	}

	page := c.QueryInt("page", 1)
	size := c.QueryInt("size", 10)

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	response, err := h.service.GetClinicMedicalHistoriesPaginated(clinicID, page, size)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}

func (h *Handler) AddDocumentsToMedicalHistory(c *fiber.Ctx) error {
	var dto AddDocumentsToMedicalHistoryDTO

	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if dto.MedicalHistoryId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Medical history ID is required",
		})
	}

	if len(dto.Documents) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "At least one document is required",
		})
	}

	response, err := h.service.AddDocumentsToMedicalHistory(dto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}

func (h *Handler) AddDocumentsToConsultation(c *fiber.Ctx) error {
	var dto AddDocumentsToConsultationDTO

	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if dto.ConsultationId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Consultation ID is required",
		})
	}

	if len(dto.Documents) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "At least one document is required",
		})
	}

	response, err := h.service.AddDocumentsToConsultation(dto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}

func (h *Handler) GetDocumentsByMedicalHistory(c *fiber.Ctx) error {
	medicalHistoryId := c.Params("medicalHistoryId")
	if medicalHistoryId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Medical history ID is required",
		})
	}

	documents, err := h.service.GetDocumentsByMedicalHistory(medicalHistoryId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"documents": documents,
		"count":     len(documents),
	})
}

func (h *Handler) GetDocumentsByConsultation(c *fiber.Ctx) error {
	consultationId := c.Params("consultationId")
	if consultationId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Consultation ID is required",
		})
	}

	documents, err := h.service.GetDocumentsByConsultation(consultationId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"documents": documents,
		"count":     len(documents),
	})
}
