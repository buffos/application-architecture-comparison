package application

import "layered-architecture/internal/domain"

type CustomerRepository interface {
	Save(customer domain.Customer) error
	FindByID(id string) (domain.Customer, error)
}

type CustomerService struct {
	repo CustomerRepository
}

func NewCustomerService(repo CustomerRepository) CustomerService {
	return CustomerService{repo: repo}
}

func (s CustomerService) CreateCustomer(name string, tier string, paymentTerms string) (domain.Customer, error) {
	customer, err := domain.NewCustomer(name, tier, paymentTerms)
	if err != nil {
		return domain.Customer{}, err
	}

	if err := s.repo.Save(customer); err != nil {
		return domain.Customer{}, err
	}

	return customer, nil
}
