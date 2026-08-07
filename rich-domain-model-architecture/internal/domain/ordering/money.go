package ordering

import (
	"errors"
	"strings"
)

var (
	ErrMoneyAmountNegative = errors.New("money amount cannot be negative")
	ErrCurrencyRequired    = errors.New("currency is required")
	ErrCurrencyMismatch    = errors.New("money currencies must match")
)

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

func (money Money) Cents() int64     { return money.cents }
func (money Money) Currency() string { return money.currency }

func (money Money) Add(other Money) (Money, error) {
	if money.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}
	return Money{cents: money.cents + other.cents, currency: money.currency}, nil
}

func (money Money) Multiply(quantity int) (Money, error) {
	if quantity < 0 {
		return Money{}, ErrMoneyAmountNegative
	}
	return Money{cents: money.cents * int64(quantity), currency: money.currency}, nil
}
