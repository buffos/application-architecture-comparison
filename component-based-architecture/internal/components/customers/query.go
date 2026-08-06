package customers

// Reader is the public read contract provided by Customers.
type Reader interface {
	GetCustomer(query GetCustomerQuery) (CustomerDetails, error)
	ListCustomers(query ListCustomersQuery) []CustomerSummary
}

type GetCustomerQuery struct{ ID string }
type ListCustomersQuery struct{ ActiveOnly bool }

type CustomerDetails struct {
	ID     string
	Active bool
}

type CustomerSummary struct {
	ID     string
	Active bool
}

func (c *Component) GetCustomer(query GetCustomerQuery) (CustomerDetails, error) {
	customer, ok := c.customers[query.ID]
	if !ok {
		return CustomerDetails{}, ErrCustomerNotFound
	}
	return CustomerDetails{ID: customer.ID, Active: customer.Active}, nil
}

func (c *Component) ListCustomers(query ListCustomersQuery) []CustomerSummary {
	customers := make([]CustomerSummary, 0, len(c.customers))
	for _, customer := range c.customers {
		if query.ActiveOnly && !customer.Active {
			continue
		}
		customers = append(customers, CustomerSummary{ID: customer.ID, Active: customer.Active})
	}
	return customers
}
