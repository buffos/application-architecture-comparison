package pricing

import (
	"modular-monolith/internal/modules/plugins"
	"modular-monolith/internal/modules/products"
)

type QuotePricer interface {
	UnitPrice(product products.ProductForQuote) (int, error)
}

type Service struct {
	plugins plugins.Reader
}

func NewService(plugins plugins.Reader) Service {
	return Service{plugins: plugins}
}

func (s Service) UnitPrice(product products.ProductForQuote) (int, error) {
	unitPrice := product.UnitPrice

	enabled, err := s.plugins.IsEnabled("seasonal-pricing")
	if err == nil && enabled {
		unitPrice = unitPrice - (unitPrice * 5 / 100)
	}

	return unitPrice, nil
}
