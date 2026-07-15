package pricing

import (
	"testing"

	"microkernel-architecture/internal/kernel"
)

func TestServiceReturnsOriginalUnitPrice(t *testing.T) {
	service := NewService()

	price, err := service.UnitPriceForQuote(kernel.QuotePricingInput{
		ProductSKU:      "sku-001",
		ProductCategory: "Standard",
		Quantity:        2,
		UnitPrice:       15000,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if price != 15000 {
		t.Fatalf("expected original unit price 15000, got %d", price)
	}
}
