package quoting

import (
	"errors"
	"testing"
)

func TestQuotePricingServiceAppliesCustomerTierDiscounts(t *testing.T) {
	price, err := NewMoney(10000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	input := ProductPricingInput{SKU: "sku-001", ProductName: "Desk", Category: ProductCategoryStandard, BasePrice: price}
	service := NewQuotePricingService()

	tests := []struct {
		name     string
		tier     CustomerPricingTier
		expected int64
	}{
		{name: "standard", tier: CustomerPricingTierStandard, expected: 10000},
		{name: "preferred", tier: CustomerPricingTierPreferred, expected: 9500},
		{name: "enterprise", tier: CustomerPricingTierEnterprise, expected: 9000},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			line, err := service.PriceLine(input, test.tier, 2)
			if err != nil {
				t.Fatal(err)
			}
			if line.UnitPrice().Cents() != test.expected {
				t.Fatalf("unit price = %d, want %d", line.UnitPrice().Cents(), test.expected)
			}
			if line.ProductCategory() != ProductCategoryStandard || line.ProductName() != "Desk" {
				t.Fatalf("line snapshot lost product facts: category=%s name=%s", line.ProductCategory(), line.ProductName())
			}
		})
	}
}

func TestQuotePricingServiceRejectsUnknownTiers(t *testing.T) {
	price, err := NewMoney(10000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	_, err = NewQuotePricingService().PriceLine(ProductPricingInput{SKU: "sku-001", Category: ProductCategoryStandard, BasePrice: price}, CustomerPricingTier("Unknown"), 1)
	if !errors.Is(err, ErrPricingTierInvalid) {
		t.Fatalf("unknown tier returned %v", err)
	}
}
