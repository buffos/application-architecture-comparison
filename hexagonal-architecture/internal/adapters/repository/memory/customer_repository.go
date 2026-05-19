package memory

import (
	"sync"

	"hexagonal-architecture/internal/core/domain"
)

type CustomerRepository struct {
	mu        sync.RWMutex
	customers map[string]domain.Customer
}

func NewCustomerRepository() *CustomerRepository {
	return &CustomerRepository{
		customers: make(map[string]domain.Customer),
	}
}

func (r *CustomerRepository) Save(customer domain.Customer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.customers[customer.ID] = customer
	return nil
}

func (r *CustomerRepository) FindByID(id string) (domain.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	customer, ok := r.customers[id]
	if !ok {
		return domain.Customer{}, domain.ErrCustomerNotFound
	}

	return customer, nil
}

func (r *CustomerRepository) List(activeOnly bool) ([]domain.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	customers := make([]domain.Customer, 0, len(r.customers))
	for _, customer := range r.customers {
		if activeOnly && !customer.Active {
			continue
		}
		customers = append(customers, customer)
	}

	return customers, nil
}
