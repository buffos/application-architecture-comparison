package records

import "errors"

var (
	ErrDatabaseRequired   = errors.New("database is required")
	ErrCustomerIDRequired = errors.New("customer id is required")
	ErrCustomerNotFound   = errors.New("customer not found")
	ErrCustomerInactive   = errors.New("customer is inactive")
)

// Customer is an Active Record. It contains the customer fields and retains
// the database needed to load or save those fields.
type Customer struct {
	db *Database

	ID     string
	Active bool
}

// NewCustomer creates a new, unsaved Customer Active Record.
func NewCustomer(db *Database, id string, active bool) *Customer {
	return &Customer{
		db:     db,
		ID:     id,
		Active: active,
	}
}

// FindCustomer loads a Customer Active Record from the customer table.
func FindCustomer(db *Database, id string) (*Customer, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}
	if id == "" {
		return nil, ErrCustomerIDRequired
	}

	row, ok := db.customers[id]
	if !ok {
		return nil, ErrCustomerNotFound
	}

	return &Customer{
		db:     db,
		ID:     row.ID,
		Active: row.Active,
	}, nil
}

// Save writes the current Customer Active Record to its table.
func (customer *Customer) Save() error {
	if customer == nil || customer.db == nil {
		return ErrDatabaseRequired
	}
	if customer.ID == "" {
		return ErrCustomerIDRequired
	}

	customer.db.customers[customer.ID] = customerRow{
		ID:     customer.ID,
		Active: customer.Active,
	}
	return nil
}
