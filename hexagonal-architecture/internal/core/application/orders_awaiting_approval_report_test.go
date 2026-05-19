package application

import (
	"testing"

	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/adapters/services/approval"
	"hexagonal-architecture/internal/adapters/services/pricing"
	"hexagonal-architecture/internal/core/domain"
)

func TestOrdersAwaitingApprovalReport(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	pricingPolicy := pricing.NewFixedPricingPolicy()
	approvalPolicy := approval.NewCategoryApprovalPolicy()

	_ = customerRepo.Save(domain.Customer{ID: "customer-001", Active: true})
	_ = customerRepo.Save(domain.Customer{ID: "customer-002", Active: true})
	_ = productRepo.Save(domain.Product{SKU: "DESK-001", Name: "Executive Desk", Category: "CustomBuild", BasePrice: 50000, Available: true, ReturnWindowDays: 30})
	_ = productRepo.Save(domain.Product{SKU: "CHAIR-001", Name: "Office Chair", Category: "Standard", BasePrice: 10000, Available: true, ReturnWindowDays: 30})

	createQuote := NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)
	submitQuote := NewSubmitQuoteUseCase(quoteRepo, approvalPolicy)
	reportUseCase := NewGetOrdersAwaitingApprovalReportUseCase(quoteRepo)

	pendingA, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(pendingA.ID, "DESK-001", 1)
	_, _ = submitQuote.Execute(pendingA.ID)

	pendingB, _ := createQuote.Execute("customer-002")
	_, _ = addQuoteLine.Execute(pendingB.ID, "DESK-001", 2)
	_, _ = submitQuote.Execute(pendingB.ID)

	approved, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(approved.ID, "CHAIR-001", 1)
	_, _ = submitQuote.Execute(approved.ID)

	report, err := reportUseCase.Execute()
	if err != nil {
		t.Fatalf("expected report to succeed, got %v", err)
	}

	if len(report) != 2 {
		t.Fatalf("expected 2 approval queue rows, got %d", len(report))
	}

	if report[0].QuoteID != pendingA.ID || report[0].CustomerID != "customer-001" || report[0].LineCount != 1 || report[0].TotalAmount != 50000 {
		t.Fatalf("unexpected first row: %+v", report[0])
	}

	if report[1].QuoteID != pendingB.ID || report[1].CustomerID != "customer-002" || report[1].LineCount != 1 || report[1].TotalAmount != 100000 {
		t.Fatalf("unexpected second row: %+v", report[1])
	}
}
