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

func (h *Handler) RescheduleAppointment(c *fiber.Ctx) error {
	appointmentId := c.Params("id")
	var newDateTimeDTO NewDateTimeDTO

	if err := c.BodyParser(&newDateTimeDTO); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	rescheduleAppointmentError := h.service.RescheduleAppointment(appointmentId, newDateTimeDTO.NewDateTime)
	if rescheduleAppointmentError != nil {
		return rescheduleAppointmentError
	}
	return nil
}

func (h *Handler) CancelAppointment(c *fiber.Ctx) error {
	appointmentId := c.Params("id")

	cancelAppointmentError := h.service.CancelAppointment(appointmentId)
	if cancelAppointmentError != nil {
		return cancelAppointmentError
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

func (h *Handler) GetAllAppointmentsByUserId(c *fiber.Ctx) error {
	userId := c.Params("id")
	var appointment []AppointmentWithNamesDTO

	appointment, appointmentError := h.service.GetAllAppointmentsByUserId(userId)

	if appointmentError != nil {
		return appointmentError
	}

	return c.JSON(appointment)
}
