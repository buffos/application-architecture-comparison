package pricing

import (
	"testing"

	"domain-driven-design-architecture/internal/domain/quoting"
)

func TestSeasonalPolicyComposesTierPricing(t *testing.T) {
	base, err := quoting.NewMoney(10000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	service := quoting.NewQuotePricingServiceWithPolicy(NewSeasonalPolicy(quoting.TierPricingPolicy{}, 10))
	line, err := service.PriceLine(quoting.ProductPricingInput{SKU: "sku-001", BasePrice: base}, quoting.CustomerPricingTierPreferred, 1)
	if err != nil {
		t.Fatal(err)
	}
	if line.UnitPrice().Cents() != 8550 {
		t.Fatalf("price = %d, want 8550", line.UnitPrice().Cents())
	}
}
