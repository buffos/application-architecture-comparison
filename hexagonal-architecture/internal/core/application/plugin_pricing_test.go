package application

import (
	"testing"

	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/adapters/services/pricing"
	"hexagonal-architecture/internal/core/domain"
)

func TestEnabledPricingPluginChangesQuotePricing(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	pluginRepo := memory.NewPluginRepository()

	_ = customerRepo.Save(domain.Customer{ID: "customer-001", Active: true})
	_ = productRepo.Save(domain.Product{SKU: "CHAIR-001", Name: "Office Chair", Category: "Standard", BasePrice: 10000, Available: true, ReturnWindowDays: 30})

	registerPlugin := NewRegisterPricingPluginUseCase(pluginRepo)
	enablePlugin := NewEnablePluginUseCase(pluginRepo)
	listPlugins := NewListPluginsUseCase(pluginRepo)
	pricingPolicy := pricing.NewPluginAwarePolicy(pricing.NewFixedPricingPolicy(), pluginRepo)
	createQuote := NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)

	_, err := registerPlugin.Execute("seasonal-pricing", 5)
	if err != nil {
		t.Fatalf("expected plugin registration to succeed, got %v", err)
	}

	plugin, err := enablePlugin.Execute("seasonal-pricing")
	if err != nil {
		t.Fatalf("expected plugin enable to succeed, got %v", err)
	}

	if plugin.Status != domain.PluginStatusEnabled {
		t.Fatalf("expected plugin status %s, got %s", domain.PluginStatusEnabled, plugin.Status)
	}

	quote, _ := createQuote.Execute("customer-001")
	quote, err = addQuoteLine.Execute(quote.ID, "CHAIR-001", 2)
	if err != nil {
		t.Fatalf("expected add quote line to succeed, got %v", err)
	}

	if quote.Lines[0].AdjustedUnitPrice != 9500 {
		t.Fatalf("expected adjusted unit price 9500, got %d", quote.Lines[0].AdjustedUnitPrice)
	}

	plugins, err := listPlugins.Execute()
	if err != nil {
		t.Fatalf("expected list plugins to succeed, got %v", err)
	}

	if len(plugins) != 1 || plugins[0].Key != "seasonal-pricing" {
		t.Fatalf("expected one plugin seasonal-pricing, got %+v", plugins)
	}
}
