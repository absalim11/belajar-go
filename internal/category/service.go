package category

type Service interface {
	GetAll() ([]Category, error)
	GetByID(id int) (*Category, error)
	Create(req CreateCategoryRequest) (*Category, error)
	Update(id int, req UpdateCategoryRequest) (*Category, error)
	Delete(id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAll() ([]Category, error) {
	return s.repo.GetAll()
}

func (s *service) GetByID(id int) (*Category, error) {
	return s.repo.GetByID(id)
}

func (s *service) Create(req CreateCategoryRequest) (*Category, error) {
	return s.repo.Create(req)
}

func (s *service) Update(id int, req UpdateCategoryRequest) (*Category, error) {
	return s.repo.Update(id, req)
}

func (s *service) Delete(id int) error {
	return s.repo.Delete(id)
}
