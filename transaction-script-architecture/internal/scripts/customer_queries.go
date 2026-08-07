package scripts

import (
	"sort"

	"transaction-script-architecture/internal/data"
)

func GetCustomer(store *data.Store, customerID string) (data.Customer, error) {
	if store == nil {
		return data.Customer{}, ErrStoreRequired
	}
	if customerID == "" {
		return data.Customer{}, ErrCustomerIDNeeded
	}

	customer, ok := store.Customers[customerID]
	if !ok {
		return data.Customer{}, ErrCustomerNotFound
	}

	return customer, nil
}

func ListCustomers(store *data.Store, activeOnly bool) ([]data.Customer, error) {
	if store == nil {
		return nil, ErrStoreRequired
	}

	customers := make([]data.Customer, 0, len(store.Customers))
	for _, customer := range store.Customers {
		if activeOnly && !customer.Active {
			continue
		}
		customers = append(customers, customer)
	}

	sort.Slice(customers, func(i, j int) bool {
		return customers[i].ID < customers[j].ID
	})

	return customers, nil
}
