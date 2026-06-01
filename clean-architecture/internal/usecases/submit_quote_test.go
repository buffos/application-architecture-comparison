package usecases

import (
	"testing"

	"clean-architecture/internal/entities"
)

type stubSubmitQuoteOutput struct {
	output SubmitQuoteOutput
}

func (o *stubSubmitQuoteOutput) Present(output SubmitQuoteOutput) error {
	o.output = output
	return nil
}

func TestSubmitQuoteInteractorSubmitsDraftWithLines(t *testing.T) {
	quotes := &stubQuoteEditor{
		quote: entities.Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     entities.QuoteStatusDraft,
			Lines: []entities.QuoteLine{
				{SKU: "CHAIR-001", Quantity: 2},
			},
		},
	}
	output := &stubSubmitQuoteOutput{}

	interactor := NewSubmitQuoteInteractor(quotes, output)

	err := interactor.Execute(SubmitQuoteInput{QuoteID: "quote-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if quotes.saved.Status != entities.QuoteStatusSubmitted {
		t.Fatalf("expected saved status %s, got %s", entities.QuoteStatusSubmitted, quotes.saved.Status)
	}

	if output.output.Status != entities.QuoteStatusSubmitted {
		t.Fatalf("expected presenter status %s, got %s", entities.QuoteStatusSubmitted, output.output.Status)
	}
}

func TestSubmitQuoteInteractorRejectsQuoteWithoutLines(t *testing.T) {
	quotes := &stubQuoteEditor{
		quote: entities.Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     entities.QuoteStatusDraft,
		},
	}
	output := &stubSubmitQuoteOutput{}

	interactor := NewSubmitQuoteInteractor(quotes, output)

	err := interactor.Execute(SubmitQuoteInput{QuoteID: "quote-001"})
	if err != entities.ErrQuoteCannotSubmitWithoutLines {
		t.Fatalf("expected %v, got %v", entities.ErrQuoteCannotSubmitWithoutLines, err)
	}
}
