package quoting

import (
	"errors"
	"testing"
)

func TestQuotePricingServiceAppliesCustomerTierDiscounts(t *testing.T) {
	basePrice, err := NewMoney(15000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	service := NewQuotePricingService()
	tests := []struct {
		name      string
		tier      CustomerPricingTier
		wantCents int64
	}{
		{name: "standard", tier: CustomerPricingTierStandard, wantCents: 15000},
		{name: "preferred", tier: CustomerPricingTierPreferred, wantCents: 14250},
		{name: "enterprise", tier: CustomerPricingTierEnterprise, wantCents: 13500},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			line, err := service.PriceLine(ProductPricingInput{SKU: "sku-001", BasePrice: basePrice}, test.tier, 2)
			if err != nil {
				t.Fatal(err)
			}
			if line.UnitPrice().Cents() != test.wantCents || line.Quantity() != 2 {
				t.Fatalf("line price=%d quantity=%d, want price=%d quantity=2", line.UnitPrice().Cents(), line.Quantity(), test.wantCents)
			}
		})
	}
}

func TestQuotePricingServiceRejectsInvalidTier(t *testing.T) {
	basePrice, err := NewMoney(15000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	_, err = NewQuotePricingService().PriceLine(ProductPricingInput{SKU: "sku-001", BasePrice: basePrice}, CustomerPricingTier("Unknown"), 1)
	if !errors.Is(err, ErrPricingTierInvalid) {
		t.Fatalf("invalid tier returned %v", err)
	}
}
