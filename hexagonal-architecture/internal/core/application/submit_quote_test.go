package application

import (
	"testing"

	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/adapters/services/approval"
	"hexagonal-architecture/internal/adapters/services/pricing"
	"hexagonal-architecture/internal/core/domain"
)

func TestSubmitQuoteUsesApprovalPolicy(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	pricingPolicy := pricing.NewFixedPricingPolicy()
	approvalPolicy := approval.NewCategoryApprovalPolicy()

	_ = customerRepo.Save(domain.Customer{ID: "customer-001", Active: true})
	_ = productRepo.Save(domain.Product{
		SKU:       "CHAIR-001",
		Name:      "Office Chair",
		Category:  "Standard",
		BasePrice: 10000,
		Available: true,
	})
	_ = productRepo.Save(domain.Product{
		SKU:       "DESK-001",
		Name:      "Executive Desk",
		Category:  "CustomBuild",
		BasePrice: 50000,
		Available: true,
	})

	createQuote := NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)
	submitQuote := NewSubmitQuoteUseCase(quoteRepo, approvalPolicy)
	approveQuote := NewApproveQuoteUseCase(quoteRepo)

	standardQuote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(standardQuote.ID, "CHAIR-001", 1)
	submittedStandard, err := submitQuote.Execute(standardQuote.ID)
	if err != nil {
		t.Fatalf("expected standard quote submission to succeed, got %v", err)
	}

	if submittedStandard.Status != domain.QuoteStatusApproved {
		t.Fatalf("expected standard quote status %s, got %s", domain.QuoteStatusApproved, submittedStandard.Status)
	}

	customQuote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(customQuote.ID, "DESK-001", 1)
	submittedCustom, err := submitQuote.Execute(customQuote.ID)
	if err != nil {
		t.Fatalf("expected custom quote submission to succeed, got %v", err)
	}

	if submittedCustom.Status != domain.QuoteStatusPendingApproval {
		t.Fatalf("expected custom quote status %s, got %s", domain.QuoteStatusPendingApproval, submittedCustom.Status)
	}

	approvedCustom, err := approveQuote.Execute(customQuote.ID)
	if err != nil {
		t.Fatalf("expected quote approval to succeed, got %v", err)
	}

	if approvedCustom.Status != domain.QuoteStatusApproved {
		t.Fatalf("expected approved quote status %s, got %s", domain.QuoteStatusApproved, approvedCustom.Status)
	}
}
