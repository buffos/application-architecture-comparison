package customer

import "errors"

var (
	ErrCustomerIDRequired      = errors.New("customer id is required")
	ErrCustomerTierInvalid     = errors.New("customer tier is invalid")
	ErrPaymentTermsInvalid     = errors.New("payment terms are invalid")
	ErrCustomerAlreadyActive   = errors.New("customer is already active")
	ErrCustomerAlreadyInactive = errors.New("customer is already inactive")
)

type CustomerID string
type CustomerTier string
type PaymentTerms string

const (
	CustomerTierStandard   CustomerTier = "Standard"
	CustomerTierPreferred  CustomerTier = "Preferred"
	CustomerTierEnterprise CustomerTier = "Enterprise"

	PaymentTermsPrepaid   PaymentTerms = "Prepaid"
	PaymentTermsInvoice30 PaymentTerms = "Invoice30"
)

// Customer is the aggregate root for the Customer bounded context.
type Customer struct {
	id           CustomerID
	tier         CustomerTier
	paymentTerms PaymentTerms
	active       bool
}

func NewCustomer(id CustomerID, tier CustomerTier, terms PaymentTerms) (Customer, error) {
	if id == "" {
		return Customer{}, ErrCustomerIDRequired
	}
	if !validTier(tier) {
		return Customer{}, ErrCustomerTierInvalid
	}
	if !validPaymentTerms(terms) {
		return Customer{}, ErrPaymentTermsInvalid
	}
	return Customer{id: id, tier: tier, paymentTerms: terms, active: true}, nil
}

func (c Customer) ID() CustomerID             { return c.id }
func (c Customer) Tier() CustomerTier         { return c.tier }
func (c Customer) PaymentTerms() PaymentTerms { return c.paymentTerms }
func (c Customer) Active() bool               { return c.active }

func (c *Customer) Deactivate() error {
	if !c.active {
		return ErrCustomerAlreadyInactive
	}
	c.active = false
	return nil
}

func (c *Customer) Activate() error {
	if c.active {
		return ErrCustomerAlreadyActive
	}
	c.active = true
	return nil
}

func (c *Customer) ChangeTier(tier CustomerTier) error {
	if !validTier(tier) {
		return ErrCustomerTierInvalid
	}
	c.tier = tier
	return nil
}

func validTier(tier CustomerTier) bool {
	switch tier {
	case CustomerTierStandard, CustomerTierPreferred, CustomerTierEnterprise:
		return true
	default:
		return false
	}
}

func validPaymentTerms(terms PaymentTerms) bool {
	switch terms {
	case PaymentTermsPrepaid, PaymentTermsInvoice30:
		return true
	default:
		return false
	}
}
