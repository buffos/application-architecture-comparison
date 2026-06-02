package memory

import (
	"sync"

	"microkernel-architecture/internal/plugins/customers"
)

type CustomerRepository struct {
	mu        sync.RWMutex
	customers map[string]customers.Customer
}

func NewCustomerRepository() *CustomerRepository {
	return &CustomerRepository{
		customers: make(map[string]customers.Customer),
	}
}

func (r *CustomerRepository) Save(customer customers.Customer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.customers[customer.ID] = customer
	return nil
}

func (r *CustomerRepository) FindByID(id string) (customers.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	customer, ok := r.customers[id]
	if !ok {
		return customers.Customer{}, customers.ErrCustomerNotFound
	}

	return customer, nil
}
