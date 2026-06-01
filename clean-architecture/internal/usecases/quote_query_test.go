package usecases

import (
	"testing"

	"clean-architecture/internal/entities"
)

type stubQuoteLister struct {
	quotes []entities.Quote
	err    error
	status string
}

func (g *stubQuoteLister) ListByStatus(status string) ([]entities.Quote, error) {
	g.status = status
	if g.err != nil {
		return nil, g.err
	}

	return g.quotes, nil
}

type stubListQuotesOutput struct {
	output ListQuotesOutput
}

func (o *stubListQuotesOutput) Present(output ListQuotesOutput) error {
	o.output = output
	return nil
}

func TestListQuotesInteractorFiltersByStatus(t *testing.T) {
	quotes := &stubQuoteLister{
		quotes: []entities.Quote{
			{
				ID:         "quote-001",
				CustomerID: "customer-001",
				Status:     entities.QuoteStatusApproved,
				Lines:      []entities.QuoteLine{{SKU: "CHAIR-001"}},
			},
			{
				ID:         "quote-002",
				CustomerID: "customer-002",
				Status:     entities.QuoteStatusApproved,
				Lines:      []entities.QuoteLine{{SKU: "DESK-001"}},
			},
		},
	}
	output := &stubListQuotesOutput{}
	interactor := NewListQuotesInteractor(quotes, output)

	err := interactor.Execute(ListQuotesInput{Status: entities.QuoteStatusApproved})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if quotes.status != entities.QuoteStatusApproved {
		t.Fatalf("expected status filter %s, got %s", entities.QuoteStatusApproved, quotes.status)
	}

	if output.output.Count != 2 {
		t.Fatalf("expected 2 quotes, got %d", output.output.Count)
	}

	if output.output.Quotes[0].QuoteID != "quote-001" {
		t.Fatalf("expected first quote id quote-001, got %s", output.output.Quotes[0].QuoteID)
	}
}
