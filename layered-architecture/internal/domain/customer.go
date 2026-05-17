package domain

import (
	"errors"
	"fmt"
	"sync/atomic"
)

var customerSequence uint64

var ErrCustomerNameRequired = errors.New("customer name is required")
var ErrCustomerTierRequired = errors.New("customer tier is required")
var ErrCustomerPaymentTermsRequired = errors.New("customer payment terms are required")
var ErrCustomerNotFound = errors.New("customer not found")
var ErrCustomerInactive = errors.New("customer is inactive")

type Customer struct {
	ID           string
	Name         string
	Tier         string
	PaymentTerms string
	Active       bool
}

func NewCustomer(name string, tier string, paymentTerms string) (Customer, error) {
	if name == "" {
		return Customer{}, ErrCustomerNameRequired
	}

	if tier == "" {
		return Customer{}, ErrCustomerTierRequired
	}

	if paymentTerms == "" {
		return Customer{}, ErrCustomerPaymentTermsRequired
	}

	id := atomic.AddUint64(&customerSequence, 1)

	return Customer{
		ID:           fmt.Sprintf("cust-%03d", id),
		Name:         name,
		Tier:         tier,
		PaymentTerms: paymentTerms,
		Active:       true,
	}, nil
}
