package customers

import "errors"

var (
	ErrCustomerIDRequired = errors.New("customer id is required")
	ErrCustomerNotFound   = errors.New("customer not found")
	ErrCustomerInactive   = errors.New("customer is inactive")
)

type Customer struct {
	ID     string
	Active bool
}
