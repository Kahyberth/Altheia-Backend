package clinicOwner

type Service interface {
	CreateClinicOwner(creaClinicOwnerDto CreateClinicOwnerDto) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) CreateClinicOwner(creaClinicOwnerDto CreateClinicOwnerDto) error {
	s.CreateClinicOwner(creaClinicOwnerDto)
	return nil
}
