package usecases

import (
	"testing"

	"clean-architecture/internal/entities"
)

type stubApproveQuoteOutput struct {
	output ApproveQuoteOutput
}

func (o *stubApproveQuoteOutput) Present(output ApproveQuoteOutput) error {
	o.output = output
	return nil
}

func TestApproveQuoteInteractorApprovesPendingQuote(t *testing.T) {
	quotes := &stubQuoteEditor{
		quote: entities.Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     entities.QuoteStatusPendingApproval,
			Lines: []entities.QuoteLine{
				{SKU: "DESK-001", Quantity: 1},
			},
		},
	}
	output := &stubApproveQuoteOutput{}

	interactor := NewApproveQuoteInteractor(quotes, output)

	err := interactor.Execute(ApproveQuoteInput{QuoteID: "quote-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if quotes.saved.Status != entities.QuoteStatusApproved {
		t.Fatalf("expected saved status %s, got %s", entities.QuoteStatusApproved, quotes.saved.Status)
	}

	if output.output.Status != entities.QuoteStatusApproved {
		t.Fatalf("expected output status %s, got %s", entities.QuoteStatusApproved, output.output.Status)
	}
}

func TestApproveQuoteInteractorRejectsNonPendingQuote(t *testing.T) {
	quotes := &stubQuoteEditor{
		quote: entities.Quote{
			ID:         "quote-002",
			CustomerID: "customer-001",
			Status:     entities.QuoteStatusApproved,
			Lines: []entities.QuoteLine{
				{SKU: "CHAIR-001", Quantity: 1},
			},
		},
	}
	output := &stubApproveQuoteOutput{}

	interactor := NewApproveQuoteInteractor(quotes, output)

	err := interactor.Execute(ApproveQuoteInput{QuoteID: "quote-002"})
	if err != entities.ErrQuoteCannotTransition {
		t.Fatalf("expected %v, got %v", entities.ErrQuoteCannotTransition, err)
	}
}
