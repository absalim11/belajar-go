package product

type Service interface {
	GetAll(nameFilter string) ([]ProductDetail, error)
	GetByID(id int) (*ProductDetail, error)
	Create(req CreateProductRequest) (*Product, error)
	Update(id int, req UpdateProductRequest) (*Product, error)
	Delete(id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAll(nameFilter string) ([]ProductDetail, error) {
	return s.repo.GetAll(nameFilter)
}

func (s *service) GetByID(id int) (*ProductDetail, error) {
	return s.repo.GetByID(id)
}

func (s *service) Create(req CreateProductRequest) (*Product, error) {
	return s.repo.Create(req)
}

func (s *service) Update(id int, req UpdateProductRequest) (*Product, error) {
	return s.repo.Update(id, req)
}

func (s *service) Delete(id int) error {
	return s.repo.Delete(id)
}
