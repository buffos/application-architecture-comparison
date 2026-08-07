package catalog

import (
	"errors"
	"strings"
)

var (
	ErrPriceAmountNegative   = errors.New("price amount cannot be negative")
	ErrPriceCurrencyRequired = errors.New("price currency is required")
)

// Price is the Catalog context's money value object. It is intentionally a
// different type from quoting.Money; the composition boundary translates
// catalog facts into the currency representation used by a quote.
type Price struct {
	cents    int64
	currency string
}

func NewPrice(cents int64, currency string) (Price, error) {
	if cents < 0 {
		return Price{}, ErrPriceAmountNegative
	}

	currency = strings.ToUpper(strings.TrimSpace(currency))
	if currency == "" {
		return Price{}, ErrPriceCurrencyRequired
	}

	return Price{cents: cents, currency: currency}, nil
}

func (price Price) Cents() int64     { return price.cents }
func (price Price) Currency() string { return price.currency }
