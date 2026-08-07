package records

import "sort"

// GetOrder is the named detail-query form of FindOrder.
func GetOrder(db *Database, id string) (*Order, error) {
	return FindOrder(db, id)
}

// ListOrders returns reconstructed Active Record snapshots ordered by order
// ID. An empty status lists every order.
func ListOrders(db *Database, status string) ([]*Order, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}

	ids := make([]string, 0, len(db.orders))
	for id, row := range db.orders {
		if status != "" && row.Status != status {
			continue
		}
		ids = append(ids, id)
	}
	sort.Strings(ids)

	orders := make([]*Order, 0, len(ids))
	for _, id := range ids {
		order, err := FindOrder(db, id)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
