package transaction

type Service interface {
	Checkout(items []CheckoutItem) (*Transaction, error)
	GetDailySalesReport() (*DailySalesReport, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Checkout(items []CheckoutItem) (*Transaction, error) {
	return s.repo.CreateTransaction(items)
}

func (s *service) GetDailySalesReport() (*DailySalesReport, error) {
	return s.repo.GetDailySalesReport()
}
