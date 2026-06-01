package entities

import "errors"

var ErrProductNotFound = errors.New("product not found")
var ErrProductUnavailable = errors.New("product is unavailable")

type Product struct {
	SKU        string
	Name       string
	BasePrice  int
	Available  bool
}

func (p Product) EnsureAvailable() error {
	if !p.Available {
		return ErrProductUnavailable
	}

	return nil
}
