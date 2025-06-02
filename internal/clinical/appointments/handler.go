package appointments

import (
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{s}
}

func (h *Handler) CreateAppointment(c *fiber.Ctx) error {
	var createAppointmentDTO CreateAppointmentDTO
	if err := c.BodyParser(&createAppointmentDTO); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	createAppointmentError := h.service.CreateAppointment(createAppointmentDTO)
	if createAppointmentError != nil {
		return createAppointmentError
	}
	return nil
}

func (h *Handler) GetAllAppointments(c *fiber.Ctx) error {
	var appointment []MedicalAppointment

	appointment, appointmentError := h.service.GetAllAppointments()

	if appointmentError != nil {
		return appointmentError
	}

	return c.JSON(appointment)
}

func (h *Handler) UpdateAppointmentStatus(c *fiber.Ctx) error {
	appointmentId := c.Params("id")
	status := c.Params("status")

	err := h.service.UpdateAppointmentStatus(appointmentId, AppointmentStatus(status))
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) GetAllAppointmentsByMedicId(c *fiber.Ctx) error {
	physicianId := c.Params("id")
	var appointment []AppointmentWithNamesDTO

	appointment, appointmentError := h.service.GetAllAppointmentsByMedicId(physicianId)

	if appointmentError != nil {
		return appointmentError
	}

	return c.JSON(appointment)
}
