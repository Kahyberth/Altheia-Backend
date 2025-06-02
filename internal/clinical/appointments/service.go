package appointments

type Service interface {
	CreateAppointment(createAppointmentDTO CreateAppointmentDTO) error
	GetAllAppointments() ([]MedicalAppointment, error)
	UpdateAppointmentStatus(appointmentId string, status AppointmentStatus) error
	GetAllAppointmentsByMedicId(medicId string) ([]AppointmentWithNamesDTO, error)
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
