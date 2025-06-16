package appointments

import "time"

type Service interface {
	CreateAppointment(createAppointmentDTO CreateAppointmentDTO) error
	GetAllAppointments() ([]MedicalAppointment, error)
	UpdateAppointmentStatus(appointmentId string, status AppointmentStatus) error
	GetAllAppointmentsByMedicId(medicId string) ([]AppointmentWithNamesDTO, error)
	GetAllAppointmentsByUserId(userId string) ([]AppointmentWithNamesDTO, error)
	CancelAppointment(appointmentId string) error
	RescheduleAppointment(appointmentId string, newDateTime time.Time) error
}
type service struct {
	repo Repository
}

func (s *service) UpdateAppointmentStatus(appointmentId string, status AppointmentStatus) error {
	err := s.repo.UpdateAppointmentStatus(appointmentId, status)
	if err != nil {
		return err
	}
	return nil
}

func NewService(r Repository) Service { return &service{r} }

func (s *service) CreateAppointment(createAppointmentDTO CreateAppointmentDTO) error {
	err := s.repo.CreateAppointment(createAppointmentDTO)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) RescheduleAppointment(appointmentId string, newDateTime time.Time) error {
	err := s.repo.ReScheduleAppointment(appointmentId, newDateTime)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) CancelAppointment(appointmentId string) error {
	err := s.repo.CancelAppointment(appointmentId)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetAllAppointments() ([]MedicalAppointment, error) {
	var appointment []MedicalAppointment

	appointment, err := s.repo.GetAllAppointments()

	if err != nil {
		return appointment, err
	}

	return appointment, nil
}

func (s *service) GetAllAppointmentsByMedicId(medicId string) ([]AppointmentWithNamesDTO, error) {
	var appointment []AppointmentWithNamesDTO

	appointment, err := s.repo.GetAllAppointmentsByMedicId(medicId)

	if err != nil {
		return appointment, err
	}

	return appointment, nil
}

func (s *service) GetAllAppointmentsByUserId(userId string) ([]AppointmentWithNamesDTO, error) {
	var appointment []AppointmentWithNamesDTO

	appointment, err := s.repo.GetAllAppointmentsByUserId(userId)

	if err != nil {
		return appointment, err
	}

	return appointment, nil
}
