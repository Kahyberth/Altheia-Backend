package clinical

type Service interface {
	CreateClinical(createClinicDto CreateClinicDTO) error
	CreateEps(epsDto CreateEpsDto) error
	GetAllEps(page int, pagSize int) ([]EPS, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) CreateClinical(createClinicDto CreateClinicDTO) error {
	err := s.repo.CreateClinic(createClinicDto)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) CreateEps(epsDto CreateEpsDto) error {
	err := s.repo.CreateEps(epsDto)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetAllEps(page int, pagSize int) ([]EPS, error) {
	var eps []EPS
	eps, epsError := s.repo.GetAllEps(page, pagSize)

	if epsError != nil {
		return eps, epsError
	}

	return eps, nil
}
