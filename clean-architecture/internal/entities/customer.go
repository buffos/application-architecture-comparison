package entities

import "errors"

var ErrCustomerInactive = errors.New("customer is inactive")
var ErrCustomerNotFound = errors.New("customer not found")

type Customer struct {
	ID     string
	Active bool
}

func (c Customer) EnsureActive() error {
	if !c.Active {
		return ErrCustomerInactive
	}

	return nil
}
