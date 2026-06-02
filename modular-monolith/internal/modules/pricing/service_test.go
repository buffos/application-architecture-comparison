package pricing

import (
	"testing"

	"modular-monolith/internal/modules/products"
)

type stubPluginReader struct {
	enabled bool
	err     error
}

func (r stubPluginReader) IsEnabled(pluginID string) (bool, error) {
	return r.enabled, r.err
}

func TestUnitPriceUsesBasePriceWhenPluginDisabled(t *testing.T) {
	service := NewService(stubPluginReader{})

	price, err := service.UnitPrice(products.ProductForQuote{SKU: "sku-001", UnitPrice: 10000})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if price != 10000 {
		t.Fatalf("expected 10000, got %d", price)
	}
}

func TestUnitPriceAppliesSeasonalPricingWhenPluginEnabled(t *testing.T) {
	service := NewService(stubPluginReader{enabled: true})

	price, err := service.UnitPrice(products.ProductForQuote{SKU: "sku-001", UnitPrice: 10000})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if price != 9500 {
		t.Fatalf("expected 9500, got %d", price)
	}
}
