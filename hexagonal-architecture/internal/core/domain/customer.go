package domain

import "errors"

var ErrCustomerNotFound = errors.New("customer not found")
var ErrCustomerInactive = errors.New("customer is inactive")

type Customer struct {
	ID     string
	Active bool
}
