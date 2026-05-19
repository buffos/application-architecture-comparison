package application

import (
	"testing"

	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/adapters/services/approval"
	"hexagonal-architecture/internal/adapters/services/pricing"
	"hexagonal-architecture/internal/core/domain"
)

func TestListQuotesByStatus(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	pricingPolicy := pricing.NewFixedPricingPolicy()
	approvalPolicy := approval.NewCategoryApprovalPolicy()

	_ = customerRepo.Save(domain.Customer{ID: "customer-001", Active: true})
	_ = productRepo.Save(domain.Product{SKU: "CHAIR-001", Name: "Office Chair", Category: "Standard", BasePrice: 10000, Available: true, ReturnWindowDays: 30})
	_ = productRepo.Save(domain.Product{SKU: "DESK-001", Name: "Executive Desk", Category: "CustomBuild", BasePrice: 50000, Available: true, ReturnWindowDays: 30})

	createQuote := NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)
	submitQuote := NewSubmitQuoteUseCase(quoteRepo, approvalPolicy)
	listQuotes := NewListQuotesUseCase(quoteRepo)

	standardQuote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(standardQuote.ID, "CHAIR-001", 1)
	_, _ = submitQuote.Execute(standardQuote.ID)

	customQuote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(customQuote.ID, "DESK-001", 1)
	_, _ = submitQuote.Execute(customQuote.ID)

	pendingApprovalQuotes, err := listQuotes.Execute(domain.QuoteStatusPendingApproval)
	if err != nil {
		t.Fatalf("expected list pending approval quotes to succeed, got %v", err)
	}

	if len(pendingApprovalQuotes) != 1 || pendingApprovalQuotes[0].ID != customQuote.ID {
		t.Fatalf("expected one pending approval quote %s, got %+v", customQuote.ID, pendingApprovalQuotes)
	}

	approvedQuotes, err := listQuotes.Execute(domain.QuoteStatusApproved)
	if err != nil {
		t.Fatalf("expected list approved quotes to succeed, got %v", err)
	}

	if len(approvedQuotes) != 1 || approvedQuotes[0].ID != standardQuote.ID {
		t.Fatalf("expected one approved quote %s, got %+v", standardQuote.ID, approvedQuotes)
	}
}
