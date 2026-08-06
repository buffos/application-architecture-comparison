package pricing

import (
	"component-based-architecture/internal/components/plugins"
	"component-based-architecture/internal/components/products"
)

type QuotePricer interface {
	UnitPrice(product products.ProductForQuote) (int, error)
}

type Component struct{ plugins plugins.Reader }

func NewComponent(plugins plugins.Reader) *Component { return &Component{plugins: plugins} }

func (c *Component) UnitPrice(product products.ProductForQuote) (int, error) {
	unitPrice := product.UnitPrice
	enabled, err := c.plugins.IsEnabled("seasonal-pricing")
	if err == nil && enabled {
		unitPrice -= unitPrice * 5 / 100
	}
	return unitPrice, nil
}

var _ QuotePricer = (*Component)(nil)
