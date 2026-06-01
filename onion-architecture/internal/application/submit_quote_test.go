package application

import (
	"testing"

	"onion-architecture/internal/domain"
)

type stubApprovalPolicy struct {
	requiresApproval bool
	err              error
}

func (p stubApprovalPolicy) RequiresApproval(quote domain.Quote) (bool, error) {
	if p.err != nil {
		return false, p.err
	}

	return p.requiresApproval, nil
}

func TestSubmitQuoteServiceApprovesQuoteWhenPolicyDoesNotRequireReview(t *testing.T) {
	quotes := &stubQuoteStore{
		quote: domain.Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     domain.QuoteStatusDraft,
			Lines: []domain.QuoteLine{
				{
					ProductSKU:      "sku-001",
					ProductName:     "Desk",
					ProductCategory: "Standard",
					Quantity:        2,
					UnitPrice:       15000,
				},
			},
		},
	}

	service := NewSubmitQuoteService(quotes, stubApprovalPolicy{})

	result, err := service.Execute(SubmitQuoteCommand{QuoteID: "quote-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != domain.QuoteStatusApproved {
		t.Fatalf("expected status %s, got %s", domain.QuoteStatusApproved, result.Status)
	}

	if quotes.saved.Status != domain.QuoteStatusApproved {
		t.Fatalf("expected saved status %s, got %s", domain.QuoteStatusApproved, quotes.saved.Status)
	}
}

func TestSubmitQuoteServiceRejectsEmptyQuote(t *testing.T) {
	quotes := &stubQuoteStore{
		quote: domain.Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     domain.QuoteStatusDraft,
		},
	}

	service := NewSubmitQuoteService(quotes, stubApprovalPolicy{})

	_, err := service.Execute(SubmitQuoteCommand{QuoteID: "quote-001"})
	if err != domain.ErrQuoteCannotBeSubmittedWithoutLines {
		t.Fatalf("expected %v, got %v", domain.ErrQuoteCannotBeSubmittedWithoutLines, err)
	}
}

func TestSubmitQuoteServiceMarksQuotePendingApprovalWhenPolicyRequiresReview(t *testing.T) {
	quotes := &stubQuoteStore{
		quote: domain.Quote{
			ID:         "quote-002",
			CustomerID: "customer-001",
			Status:     domain.QuoteStatusDraft,
			Lines: []domain.QuoteLine{
				{
					ProductSKU:      "sku-002",
					ProductName:     "Custom Desk",
					ProductCategory: "CustomBuild",
					Quantity:        1,
					UnitPrice:       45000,
				},
			},
		},
	}

	service := NewSubmitQuoteService(quotes, stubApprovalPolicy{requiresApproval: true})

	result, err := service.Execute(SubmitQuoteCommand{QuoteID: "quote-002"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != domain.QuoteStatusPendingApproval {
		t.Fatalf("expected status %s, got %s", domain.QuoteStatusPendingApproval, result.Status)
	}
}
