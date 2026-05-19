package application

import (
	"testing"

	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/adapters/services/approval"
	"hexagonal-architecture/internal/adapters/services/pricing"
	"hexagonal-architecture/internal/core/domain"
)

func TestQuoteConversionReport(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	orderRepo := memory.NewOrderRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	inventory := memory.NewInventoryReservationAdapter(map[string]int{
		"CHAIR-001": 5,
		"DESK-001":  5,
	})
	pricingPolicy := pricing.NewFixedPricingPolicy()
	approvalPolicy := approval.NewCategoryApprovalPolicy()

	_ = customerRepo.Save(domain.Customer{ID: "customer-001", Active: true})
	_ = productRepo.Save(domain.Product{SKU: "CHAIR-001", Name: "Office Chair", Category: "Standard", BasePrice: 10000, Available: true, ReturnWindowDays: 30})
	_ = productRepo.Save(domain.Product{SKU: "DESK-001", Name: "Executive Desk", Category: "CustomBuild", BasePrice: 50000, Available: true, ReturnWindowDays: 30})

	createQuote := NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)
	submitQuote := NewSubmitQuoteUseCase(quoteRepo, approvalPolicy)
	approveQuote := NewApproveQuoteUseCase(quoteRepo)
	convertQuote := NewConvertQuoteToOrderUseCase(quoteRepo, orderRepo, inventory)
	reportUseCase := NewGetQuoteConversionReportUseCase(quoteRepo, orderRepo)

	approvedAndConverted, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(approvedAndConverted.ID, "CHAIR-001", 1)
	_, _ = submitQuote.Execute(approvedAndConverted.ID)
	_, _ = convertQuote.Execute(approvedAndConverted.ID)

	pendingApproval, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(pendingApproval.ID, "DESK-001", 1)
	_, _ = submitQuote.Execute(pendingApproval.ID)
	_, _ = approveQuote.Execute(pendingApproval.ID)

	draftOnly, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(draftOnly.ID, "CHAIR-001", 2)

	report, err := reportUseCase.Execute()
	if err != nil {
		t.Fatalf("expected report to succeed, got %v", err)
	}

	if report.TotalQuotes != 3 {
		t.Fatalf("expected total quotes 3, got %d", report.TotalQuotes)
	}

	if report.ApprovedQuotes != 2 {
		t.Fatalf("expected approved quotes 2, got %d", report.ApprovedQuotes)
	}

	if report.ConvertedQuotes != 1 {
		t.Fatalf("expected converted quotes 1, got %d", report.ConvertedQuotes)
	}

	if report.ConversionRate != 1.0/3.0 {
		t.Fatalf("expected conversion rate %f, got %f", 1.0/3.0, report.ConversionRate)
	}
}
