package application

import "onion-architecture/internal/domain"

type CustomerLookup interface {
	FindByID(id string) (domain.Customer, error)
	List(activeOnly bool) ([]domain.Customer, error)
}

type GetCustomerQuery struct {
	CustomerID string
}

type CustomerDetails struct {
	CustomerID string
	Active     bool
}

type GetCustomerService struct {
	customers CustomerLookup
}

func NewGetCustomerService(customers CustomerLookup) GetCustomerService {
	return GetCustomerService{
		customers: customers,
	}
}

func (s GetCustomerService) Execute(query GetCustomerQuery) (CustomerDetails, error) {
	customer, err := s.customers.FindByID(query.CustomerID)
	if err != nil {
		return CustomerDetails{}, err
	}

	return CustomerDetails{
		CustomerID: customer.ID,
		Active:     customer.Active,
	}, nil
}
