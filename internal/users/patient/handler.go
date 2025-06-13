package patient

import (
	"Altheia-Backend/internal/users"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) RegisterPatient(c *fiber.Ctx) error {
	var patient CreatePatientInfo
	if err := c.BodyParser(&patient); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	if err := h.service.RegisterPatient(patient); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "registered successfully"})
}

func (h *Handler) UpdatePatient(c *fiber.Ctx) error {
	var patient UpdatePatientInfo
	if err := c.BodyParser(&patient); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	id := c.Params("id")
	if err := h.service.UpdatePatient(id, patient); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "updated successfully"})
}

func (h *Handler) SoftDeletePatient(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.SoftDeletePatient(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "deleted successfully"})
}

func (h *Handler) GetAllPatientsPaginated(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	result, err := h.service.GetAllPatientsPaginated(page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func (h *Handler) GetAllPatients(c *fiber.Ctx) error {
	patients, err := h.service.GetAllPatients()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(patients)
}

func (h *Handler) GetPatientByClinicId(c *fiber.Ctx) error {
	clinicId := c.Params("clinicId")

	page := c.QueryInt("page", 0)
	limit := c.QueryInt("limit", 0)

	if page > 0 && limit > 0 {
		result, err := h.service.GetPatientByClinicIdPaginated(clinicId, page, limit)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		if patients, ok := result.Result.([]users.Patient); ok {
			result.Result = buildResponse(patients)
		}

		return c.JSON(result)
	}

	patients, err := h.service.GetPatientByClinicId(clinicId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(buildResponse(patients))
}

type BasicUser struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Rol            string `json:"rol"`
	Phone          string `json:"phone"`
	DocumentNumber string `json:"document_number"`
	Status         bool   `json:"status"`
	Gender         string `json:"gender"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
	LastLogin      string `json:"lastLogin"`
}

type PatientWithUser struct {
	users.Patient
	User BasicUser `json:"user"`
}

func buildResponse(patients []users.Patient) []PatientWithUser {
	out := make([]PatientWithUser, 0, len(patients))
	for _, p := range patients {
		if p.User == nil {
			continue
		}
		u := p.User
		out = append(out, PatientWithUser{
			Patient: p,
			User: BasicUser{
				ID:             u.ID,
				Name:           u.Name,
				Email:          u.Email,
				Password:       u.Password,
				Rol:            u.Rol,
				Phone:          u.Phone,
				DocumentNumber: u.DocumentNumber,
				Status:         u.Status,
				Gender:         u.Gender,
				CreatedAt:      u.CreatedAt.Format(time.RFC3339),
				UpdatedAt:      u.UpdatedAt.Format(time.RFC3339),
				LastLogin:      u.LastLogin.Format(time.RFC3339),
			},
		})
	}
	return out
}
