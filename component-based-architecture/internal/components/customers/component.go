package customers

// Component owns customer behavior and its in-memory state for this lesson.
type Component struct {
	customers map[string]Customer
}

func NewComponent() *Component {
	return &Component{
		customers: make(map[string]Customer),
	}
}

func (c *Component) Register(customer Customer) error {
	if customer.ID == "" {
		return ErrCustomerIDRequired
	}

	c.customers[customer.ID] = customer
	return nil
}

func (c *Component) RequireActiveCustomer(id string) error {
	customer, ok := c.customers[id]
	if !ok {
		return ErrCustomerNotFound
	}

	if !customer.Active {
		return ErrCustomerInactive
	}

	return nil
}

var _ CustomerDirectory = (*Component)(nil)
