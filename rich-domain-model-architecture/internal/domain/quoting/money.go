package quoting

import (
	"errors"
	"strings"
)

var (
	ErrMoneyAmountNegative = errors.New("money amount cannot be negative")
	ErrCurrencyRequired    = errors.New("currency is required")
	ErrCurrencyMismatch    = errors.New("money currencies must match")
)

// Money is a value object. Its state is private so callers can only create
// valid amounts and combine compatible currencies through its methods.
type Money struct {
	cents    int64
	currency string
}

func NewMoney(cents int64, currency string) (Money, error) {
	if cents < 0 {
		return Money{}, ErrMoneyAmountNegative
	}

	currency = strings.ToUpper(strings.TrimSpace(currency))
	if currency == "" {
		return Money{}, ErrCurrencyRequired
	}

	return Money{cents: cents, currency: currency}, nil
}

func (m Money) Cents() int64     { return m.cents }
func (m Money) Currency() string { return m.currency }

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}

	return Money{cents: m.cents + other.cents, currency: m.currency}, nil
}

func (m Money) Multiply(quantity int) (Money, error) {
	if quantity < 0 {
		return Money{}, ErrMoneyAmountNegative
	}

	return Money{cents: m.cents * int64(quantity), currency: m.currency}, nil
}
