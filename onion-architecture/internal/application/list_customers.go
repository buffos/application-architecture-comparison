package application

type ListCustomersQuery struct {
	ActiveOnly bool
}

type ListCustomersService struct {
	customers CustomerLookup
}

func NewListCustomersService(customers CustomerLookup) ListCustomersService {
	return ListCustomersService{
		customers: customers,
	}
}

func (s ListCustomersService) Execute(query ListCustomersQuery) ([]CustomerDetails, error) {
	customers, err := s.customers.List(query.ActiveOnly)
	if err != nil {
		return nil, err
	}

	result := make([]CustomerDetails, 0, len(customers))
	for _, customer := range customers {
		result = append(result, CustomerDetails{
			CustomerID: customer.ID,
			Active:     customer.Active,
		})
	}

	return result, nil
}
