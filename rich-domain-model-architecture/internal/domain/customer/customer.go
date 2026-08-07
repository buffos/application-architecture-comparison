package customer

import "errors"

var (
	ErrCustomerIDRequired      = errors.New("customer id is required")
	ErrCustomerTierInvalid     = errors.New("customer tier is invalid")
	ErrPaymentTermsInvalid     = errors.New("payment terms are invalid")
	ErrCustomerInactive        = errors.New("inactive customer cannot create a quote")
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

// Customer is the aggregate root for the Customer domain. Its state is
// private so customer rules cannot be bypassed by public field mutation.
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

func (customer Customer) ID() CustomerID             { return customer.id }
func (customer Customer) Tier() CustomerTier         { return customer.tier }
func (customer Customer) PaymentTerms() PaymentTerms { return customer.paymentTerms }
func (customer Customer) Active() bool               { return customer.active }

func (customer Customer) EnsureCanCreateQuote() error {
	if !customer.active {
		return ErrCustomerInactive
	}
	return nil
}

func (customer *Customer) Deactivate() error {
	if !customer.active {
		return ErrCustomerAlreadyInactive
	}

	customer.active = false
	return nil
}

func (customer *Customer) Activate() error {
	if customer.active {
		return ErrCustomerAlreadyActive
	}

	customer.active = true
	return nil
}

func (customer *Customer) ChangeTier(tier CustomerTier) error {
	if !validTier(tier) {
		return ErrCustomerTierInvalid
	}

	customer.tier = tier
	return nil
}

func (customer *Customer) ChangePaymentTerms(terms PaymentTerms) error {
	if !validPaymentTerms(terms) {
		return ErrPaymentTermsInvalid
	}

	customer.paymentTerms = terms
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
