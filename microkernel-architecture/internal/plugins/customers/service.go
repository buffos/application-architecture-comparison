package customers

import "microkernel-architecture/internal/kernel"

type Service struct {
	customers Repository
}

func NewService(customers Repository) Service {
	return Service{
		customers: customers,
	}
}

func (s Service) RequireActiveCustomer(id string) error {
	customer, err := s.customers.FindByID(id)
	if err != nil {
		return err
	}

	if !customer.Active {
		return ErrCustomerInactive
	}

	return nil
}

func (s Service) GetCustomer(query kernel.GetCustomerQuery) (kernel.CustomerDetails, error) {
	customer, err := s.customers.FindByID(query.CustomerID)
	if err != nil {
		return kernel.CustomerDetails{}, err
	}

	return kernel.CustomerDetails{
		CustomerID: customer.ID,
		Active:     customer.Active,
	}, nil
}

func (s Service) ListCustomers(query kernel.ListCustomersQuery) ([]kernel.CustomerSummary, error) {
	customersList, err := s.customers.List(query.Active)
	if err != nil {
		return nil, err
	}

	results := make([]kernel.CustomerSummary, 0, len(customersList))
	for _, customer := range customersList {
		results = append(results, kernel.CustomerSummary{
			CustomerID: customer.ID,
			Active:     customer.Active,
		})
	}

	return results, nil
}
