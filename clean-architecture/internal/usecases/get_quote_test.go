package usecases

import (
	"testing"

	"clean-architecture/internal/entities"
)

type stubQuoteReader struct {
	quote entities.Quote
	err   error
}

func (g stubQuoteReader) FindByID(id string) (entities.Quote, error) {
	if g.err != nil {
		return entities.Quote{}, g.err
	}

	return g.quote, nil
}

type stubGetQuoteOutput struct {
	output GetQuoteOutput
}

func (o *stubGetQuoteOutput) Present(output GetQuoteOutput) error {
	o.output = output
	return nil
}

func TestGetQuoteInteractorLoadsQuoteAndPresentsIt(t *testing.T) {
	quotes := stubQuoteReader{
		quote: entities.Quote{
			ID:         "quote-123",
			CustomerID: "customer-001",
			Status:     entities.QuoteStatusDraft,
			Lines: []entities.QuoteLine{
				{SKU: "CHAIR-001", Quantity: 2},
			},
		},
	}
	output := &stubGetQuoteOutput{}

	interactor := NewGetQuoteInteractor(quotes, output)

	err := interactor.Execute(GetQuoteInput{QuoteID: "quote-123"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.output.QuoteID != "quote-123" {
		t.Fatalf("expected quote id quote-123, got %s", output.output.QuoteID)
	}

	if output.output.CustomerID != "customer-001" {
		t.Fatalf("expected customer id customer-001, got %s", output.output.CustomerID)
	}

	if output.output.Lines != 1 {
		t.Fatalf("expected 1 line, got %d", output.output.Lines)
	}
}

func TestGetQuoteInteractorReturnsNotFound(t *testing.T) {
	quotes := stubQuoteReader{
		err: entities.ErrQuoteNotFound,
	}
	output := &stubGetQuoteOutput{}

	interactor := NewGetQuoteInteractor(quotes, output)

	err := interactor.Execute(GetQuoteInput{QuoteID: "quote-missing"})
	if err != entities.ErrQuoteNotFound {
		t.Fatalf("expected %v, got %v", entities.ErrQuoteNotFound, err)
	}
}
