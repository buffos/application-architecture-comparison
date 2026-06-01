package domain

import "errors"

var ErrProductNotFound = errors.New("product not found")
var ErrProductInactive = errors.New("product is inactive")

type Product struct {
	SKU       string
	Name      string
	Category  string
	Active    bool
	UnitPrice int
	ReturnWindowDays int
}

func (p Product) EnsureActive() error {
	if !p.Active {
		return ErrProductInactive
	}

	return nil
}
