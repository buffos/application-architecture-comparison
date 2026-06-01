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

type stubApprovalPolicy struct {
	requiresApproval bool
	err              error
}

func (p stubApprovalPolicy) RequiresApproval(quote entities.Quote) (bool, error) {
	if p.err != nil {
		return false, p.err
	}

	return p.requiresApproval, nil
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

	interactor := NewSubmitQuoteInteractor(quotes, stubApprovalPolicy{}, output)

	err := interactor.Execute(SubmitQuoteInput{QuoteID: "quote-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if quotes.saved.Status != entities.QuoteStatusApproved {
		t.Fatalf("expected saved status %s, got %s", entities.QuoteStatusApproved, quotes.saved.Status)
	}

	if output.output.Status != entities.QuoteStatusApproved {
		t.Fatalf("expected presenter status %s, got %s", entities.QuoteStatusApproved, output.output.Status)
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

	interactor := NewSubmitQuoteInteractor(quotes, stubApprovalPolicy{}, output)

	err := interactor.Execute(SubmitQuoteInput{QuoteID: "quote-001"})
	if err != entities.ErrQuoteCannotSubmitWithoutLines {
		t.Fatalf("expected %v, got %v", entities.ErrQuoteCannotSubmitWithoutLines, err)
	}
}

func TestSubmitQuoteInteractorMarksCustomQuotePendingApproval(t *testing.T) {
	quotes := &stubQuoteEditor{
		quote: entities.Quote{
			ID:         "quote-002",
			CustomerID: "customer-001",
			Status:     entities.QuoteStatusDraft,
			Lines: []entities.QuoteLine{
				{SKU: "DESK-001", ProductCategory: "CustomBuild", Quantity: 1},
			},
		},
	}
	output := &stubSubmitQuoteOutput{}

	interactor := NewSubmitQuoteInteractor(quotes, stubApprovalPolicy{requiresApproval: true}, output)

	err := interactor.Execute(SubmitQuoteInput{QuoteID: "quote-002"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if quotes.saved.Status != entities.QuoteStatusPendingApproval {
		t.Fatalf("expected saved status %s, got %s", entities.QuoteStatusPendingApproval, quotes.saved.Status)
	}

	if output.output.Status != entities.QuoteStatusPendingApproval {
		t.Fatalf("expected presenter status %s, got %s", entities.QuoteStatusPendingApproval, output.output.Status)
	}
}
