package application

import (
	"testing"

	"layered-architecture/internal/infrastructure/memory"
)

type testPricingPlugin struct{}

func (testPricingPlugin) Key() string {
	return "seasonal-discount"
}

func (testPricingPlugin) Adjust(input PricingPluginInput) (PricingAdjustment, bool) {
	return PricingAdjustment{
		Label:         "seasonal-discount",
		AdjustedPrice: input.BasePrice - ((input.BasePrice * 10) / 100),
	}, true
}

type testPricingPluginRegistry struct{}

func (testPricingPluginRegistry) EnabledPricingPlugins() []PricingPlugin {
	return []PricingPlugin{testPricingPlugin{}}
}

func TestPricingPluginAdjustsQuoteLinePrice(t *testing.T) {
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	quoteRepo := memory.NewQuoteRepository()

	customerService := NewCustomerService(customerRepo)
	catalogService := NewCatalogService(productRepo)
	quoteService := NewQuoteService(quoteRepo, customerRepo, productRepo, testPricingPluginRegistry{})

	customer, _ := customerService.CreateCustomer("Acme Corp", "Preferred", "Invoice30")
	product, _ := catalogService.CreateProduct("CHAIR-001", "Office Chair", "Standard", 10000, true)
	quote, _ := quoteService.CreateDraftQuote(customer.ID)

	updatedQuote, err := quoteService.AddQuoteLine(quote.ID, product.SKU, 2)
	if err != nil {
		t.Fatalf("expected add line to succeed, got %v", err)
	}

	line := updatedQuote.Lines[0]
	if line.BaseUnitPrice != 10000 {
		t.Fatalf("expected base price 10000, got %d", line.BaseUnitPrice)
	}

	if line.AdjustedUnitPrice != 9000 {
		t.Fatalf("expected adjusted price 9000, got %d", line.AdjustedUnitPrice)
	}

	if line.LineTotal != 18000 {
		t.Fatalf("expected line total 18000, got %d", line.LineTotal)
	}

	if len(line.PricingAdjustments) != 1 || line.PricingAdjustments[0] != "seasonal-discount" {
		t.Fatalf("expected seasonal adjustment, got %+v", line.PricingAdjustments)
	}
}
