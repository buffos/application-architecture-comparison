package seasonalpricing

import (
	"testing"

	"microkernel-architecture/internal/kernel"
	"microkernel-architecture/internal/plugins/pricing"
)

func TestServiceDiscountsCustomBuildPrice(t *testing.T) {
	service := NewService(pricing.NewService(), 10)

	price, err := service.UnitPriceForQuote(kernel.QuotePricingInput{
		ProductSKU:      "sku-002",
		ProductCategory: "CustomBuild",
		Quantity:        2,
		UnitPrice:       45000,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if price != 40500 {
		t.Fatalf("expected discounted price 40500, got %d", price)
	}
}
