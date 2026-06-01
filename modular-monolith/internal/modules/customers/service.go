package customers

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
