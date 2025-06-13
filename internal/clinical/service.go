package clinical

type Service interface {
	CreateClinical(createClinicDto CreateClinicDTO) error
	CreateEps(epsDto CreateEpsDto) error
	GetAllEps(page int, pagSize int) ([]EPS, error)
	CreateServicesOffered(servicesOffered CreateServicesDto) error
	GetAllServicesOffered(page int, pagSize int) ([]ServicesOffered, error)
	GetClinicByOwnerID(ownerID string) (*ClinicCompleteInfoResponse, error)
	AssignServicesToClinic(dto AssignServicesClinicDTO) error
	GetClinicsByEps(epsID string, page int, pageSize int) ([]Clinic, error)
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

func (s *service) CreateServicesOffered(servicesOffered CreateServicesDto) error {
	err := s.repo.CreateServices(servicesOffered)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetAllServicesOffered(page int, pagSize int) ([]ServicesOffered, error) {
	var servicesOffered []ServicesOffered
	servicesOffered, serviceError := s.repo.GetAllServices(page, pagSize)
	if serviceError != nil {
		return nil, serviceError
	}
	return servicesOffered, nil
}

func (s *service) GetAllEps(page int, pagSize int) ([]EPS, error) {
	var eps []EPS
	eps, epsError := s.repo.GetAllEps(page, pagSize)
	if epsError != nil {
		return eps, epsError
	}
	return eps, nil
}

func (s *service) GetClinicByOwnerID(ownerID string) (*ClinicCompleteInfoResponse, error) {
	return s.repo.GetClinicByOwnerID(ownerID)
}

func (s *service) AssignServicesToClinic(dto AssignServicesClinicDTO) error {
	err := s.repo.AssignServicesToClinic(dto)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetClinicsByEps(epsID string, page int, pageSize int) ([]Clinic, error) {
	return s.repo.GetClinicsByEps(epsID, page, pageSize)
}
