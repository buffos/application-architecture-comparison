package usecases

import (
	"testing"

	"clean-architecture/internal/entities"
)

type stubQuoteEditor struct {
	quote entities.Quote
	err   error
	saved entities.Quote
}

func (g *stubQuoteEditor) FindByID(id string) (entities.Quote, error) {
	if g.err != nil {
		return entities.Quote{}, g.err
	}

	return g.quote, nil
}

func (g *stubQuoteEditor) Save(quote entities.Quote) error {
	g.saved = quote
	return nil
}

type stubProductGateway struct {
	product entities.Product
	err     error
}

func (g stubProductGateway) FindBySKU(sku string) (entities.Product, error) {
	if g.err != nil {
		return entities.Product{}, g.err
	}

	return g.product, nil
}

type stubAddQuoteLineOutput struct {
	output AddQuoteLineOutput
}

func (o *stubAddQuoteLineOutput) Present(output AddQuoteLineOutput) error {
	o.output = output
	return nil
}

func TestAddQuoteLineInteractorAddsLineAndSavesQuote(t *testing.T) {
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
			BasePrice: 10000,
			Available: true,
		},
	}
	output := &stubAddQuoteLineOutput{}

	interactor := NewAddQuoteLineInteractor(quotes, products, output)

	err := interactor.Execute(AddQuoteLineInput{
		QuoteID:  "quote-001",
		SKU:      "CHAIR-001",
		Quantity: 2,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(quotes.saved.Lines) != 1 {
		t.Fatalf("expected 1 saved line, got %d", len(quotes.saved.Lines))
	}

	if output.output.Lines != 1 {
		t.Fatalf("expected presenter line count 1, got %d", output.output.Lines)
	}
}

func TestAddQuoteLineInteractorRejectsUnavailableProduct(t *testing.T) {
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
			BasePrice: 10000,
			Available: false,
		},
	}
	output := &stubAddQuoteLineOutput{}

	interactor := NewAddQuoteLineInteractor(quotes, products, output)

	err := interactor.Execute(AddQuoteLineInput{
		QuoteID:  "quote-001",
		SKU:      "CHAIR-001",
		Quantity: 1,
	})
	if err != entities.ErrProductUnavailable {
		t.Fatalf("expected %v, got %v", entities.ErrProductUnavailable, err)
	}
}
