package customers

type GetCustomerQuery struct {
	CustomerID string
}

type CustomerDetails struct {
	CustomerID string
	Active     bool
}

type ListCustomersQuery struct {
	ActiveOnly bool
}

func (s Service) GetCustomer(query GetCustomerQuery) (CustomerDetails, error) {
	customer, err := s.customers.FindByID(query.CustomerID)
	if err != nil {
		return CustomerDetails{}, err
	}

	return CustomerDetails{
		CustomerID: customer.ID,
		Active:     customer.Active,
	}, nil
}

func (s Service) ListCustomers(query ListCustomersQuery) ([]CustomerDetails, error) {
	customers, err := s.customers.List(query.ActiveOnly)
	if err != nil {
		return nil, err
	}

	list := make([]CustomerDetails, 0, len(customers))
	for _, customer := range customers {
		list = append(list, CustomerDetails{
			CustomerID: customer.ID,
			Active:     customer.Active,
		})
	}

	return list, nil
}
