package payments

import (
	"errors"
	"strings"
)

var (
	ErrAmountNegative   = errors.New("payment amount cannot be negative")
	ErrCurrencyRequired = errors.New("currency is required")
)

type Money struct {
	cents    int64
	currency string
}

func NewMoney(cents int64, currency string) (Money, error) {
	if cents < 0 {
		return Money{}, ErrAmountNegative
	}
	currency = strings.ToUpper(strings.TrimSpace(currency))
	if currency == "" {
		return Money{}, ErrCurrencyRequired
	}
	return Money{cents: cents, currency: currency}, nil
}

func (m Money) Cents() int64     { return m.cents }
func (m Money) Currency() string { return m.currency }
