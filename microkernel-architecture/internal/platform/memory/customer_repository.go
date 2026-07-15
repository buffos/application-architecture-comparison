package memory

import (
	"slices"
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

func (r *CustomerRepository) List(active *bool) ([]customers.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	results := make([]customers.Customer, 0)
	for _, customer := range r.customers {
		if active != nil && customer.Active != *active {
			continue
		}
		results = append(results, customer)
	}

	slices.SortFunc(results, func(a customers.Customer, b customers.Customer) int {
		if a.ID < b.ID {
			return -1
		}
		if a.ID > b.ID {
			return 1
		}
		return 0
	})

	return results, nil
}
