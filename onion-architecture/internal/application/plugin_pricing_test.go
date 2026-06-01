package application

import (
	"testing"

	"onion-architecture/internal/domain"
	"onion-architecture/internal/infrastructure/memory"
	"onion-architecture/internal/infrastructure/services/pricing"
)

func TestRegisterEnableAndListPlugins(t *testing.T) {
	plugins := memory.NewPluginRepository()

	register := NewRegisterPricingPluginService(plugins)
	enable := NewEnablePluginService(plugins)
	list := NewListPluginsService(plugins)

	_, err := register.Execute(RegisterPricingPluginCommand{Name: "seasonal-pricing"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	enabled, err := enable.Execute(EnablePluginCommand{Name: "seasonal-pricing"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !enabled.Enabled {
		t.Fatalf("expected plugin to be enabled")
	}

	registered, err := list.Execute()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(registered) != 1 {
		t.Fatalf("expected 1 plugin, got %d", len(registered))
	}
}

func TestAddQuoteLineServiceUsesEnabledPricingPlugin(t *testing.T) {
	quotes := &stubQuoteStore{
		quote: domain.Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     domain.QuoteStatusDraft,
		},
	}

	products := stubProductLookup{
		product: domain.Product{
			SKU:       "sku-001",
			Name:      "Desk",
			Category:  "Standard",
			Active:    true,
			UnitPrice: 10000,
		},
	}

	plugins := memory.NewPluginRepository()
	register := NewRegisterPricingPluginService(plugins)
	enable := NewEnablePluginService(plugins)

	if _, err := register.Execute(RegisterPricingPluginCommand{Name: "seasonal-pricing"}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, err := enable.Execute(EnablePluginCommand{Name: "seasonal-pricing"}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	policy := pricing.NewPluginPolicy(pricing.NewFixedPolicy(), plugins)
	service := NewAddQuoteLineService(quotes, products, policy)

	_, err := service.Execute(AddQuoteLineCommand{
		QuoteID:    "quote-001",
		ProductSKU: "sku-001",
		Quantity:   1,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(quotes.saved.Lines) != 1 {
		t.Fatalf("expected one saved line, got %d", len(quotes.saved.Lines))
	}

	if quotes.saved.Lines[0].UnitPrice != 9500 {
		t.Fatalf("expected adjusted unit price 9500, got %d", quotes.saved.Lines[0].UnitPrice)
	}
}
