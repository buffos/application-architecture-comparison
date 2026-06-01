package usecases

import (
	"testing"

	"clean-architecture/internal/entities"
	"clean-architecture/internal/infrastructure/memory"
	pricingservice "clean-architecture/internal/infrastructure/services/pricing"
)

type stubRegisterPricingPluginOutput struct {
	output RegisterPricingPluginOutput
}

func (o *stubRegisterPricingPluginOutput) Present(output RegisterPricingPluginOutput) error {
	o.output = output
	return nil
}

type stubEnablePluginOutput struct {
	output EnablePluginOutput
}

func (o *stubEnablePluginOutput) Present(output EnablePluginOutput) error {
	o.output = output
	return nil
}

func TestEnabledPricingPluginAdjustsQuoteLinePrice(t *testing.T) {
	pluginGateway := memory.NewPluginGateway()
	registerOutput := &stubRegisterPricingPluginOutput{}
	enableOutput := &stubEnablePluginOutput{}

	registerInteractor := NewRegisterPricingPluginInteractor(pluginGateway, registerOutput)
	if err := registerInteractor.Execute(RegisterPricingPluginInput{Name: "seasonal-pricing"}); err != nil {
		t.Fatalf("expected no error registering plugin, got %v", err)
	}

	enableInteractor := NewEnablePluginInteractor(pluginGateway, enableOutput)
	if err := enableInteractor.Execute(EnablePluginInput{Name: "seasonal-pricing"}); err != nil {
		t.Fatalf("expected no error enabling plugin, got %v", err)
	}

	quotes := &stubQuoteEditor{
		quote: entities.Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     entities.QuoteStatusDraft,
		},
	}
	products := stubProductGateway{
		product: entities.Product{
			SKU:       "CHAIR-001",
			Name:      "Office Chair",
			Category:  "Standard",
			BasePrice: 10000,
			Available: true,
		},
	}
	output := &stubAddQuoteLineOutput{}
	pricingPolicy := pricingservice.NewPluginPolicy(pluginGateway)

	interactor := NewAddQuoteLineInteractor(quotes, products, pricingPolicy, output)
	err := interactor.Execute(AddQuoteLineInput{QuoteID: "quote-001", SKU: "CHAIR-001", Quantity: 2})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if quotes.saved.Lines[0].UnitPrice != 9500 {
		t.Fatalf("expected adjusted unit price 9500, got %d", quotes.saved.Lines[0].UnitPrice)
	}

	if quotes.saved.Lines[0].LineTotal != 19000 {
		t.Fatalf("expected adjusted line total 19000, got %d", quotes.saved.Lines[0].LineTotal)
	}
}
