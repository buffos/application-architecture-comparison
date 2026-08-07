package payments

import (
	"errors"
	"strings"
)

var (
	ErrMoneyAmountNegative = errors.New("money amount cannot be negative")
	ErrCurrencyRequired    = errors.New("currency is required")
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
