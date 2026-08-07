package records

import "sort"

// GetCustomer is the named detail-query form of FindCustomer.
func GetCustomer(db *Database, id string) (*Customer, error) {
	return FindCustomer(db, id)
}

// ListCustomers returns reconstructed customer snapshots ordered by customer
// ID. When activeOnly is true, inactive customers are omitted.
func ListCustomers(db *Database, activeOnly bool) ([]*Customer, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}

	ids := make([]string, 0, len(db.customers))
	for id, row := range db.customers {
		if activeOnly && !row.Active {
			continue
		}
		ids = append(ids, id)
	}
	sort.Strings(ids)

	customers := make([]*Customer, 0, len(ids))
	for _, id := range ids {
		customer, err := FindCustomer(db, id)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	return customers, nil
}
