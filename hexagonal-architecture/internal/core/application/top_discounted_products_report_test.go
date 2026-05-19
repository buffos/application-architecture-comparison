package application

import (
	"testing"

	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/core/domain"
)

func TestTopDiscountedProductsReport(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	pricingPolicy := discountedPricingPolicy{
		adjustedBySKU: map[string]int{
			"CHAIR-001": 9000,
			"DESK-001":  45000,
			"LAMP-001":  4000,
		},
	}

	_ = customerRepo.Save(domain.Customer{ID: "customer-001", Active: true})
	_ = productRepo.Save(domain.Product{SKU: "CHAIR-001", Name: "Office Chair", Category: "Standard", BasePrice: 10000, Available: true, ReturnWindowDays: 30})
	_ = productRepo.Save(domain.Product{SKU: "DESK-001", Name: "Executive Desk", Category: "CustomBuild", BasePrice: 50000, Available: true, ReturnWindowDays: 30})
	_ = productRepo.Save(domain.Product{SKU: "LAMP-001", Name: "Desk Lamp", Category: "Standard", BasePrice: 4000, Available: true, ReturnWindowDays: 30})

	createQuote := NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)
	reportUseCase := NewGetTopDiscountedProductsReportUseCase(quoteRepo)

	quoteA, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(quoteA.ID, "CHAIR-001", 2)
	_, _ = addQuoteLine.Execute(quoteA.ID, "DESK-001", 1)

	quoteB, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(quoteB.ID, "CHAIR-001", 1)
	_, _ = addQuoteLine.Execute(quoteB.ID, "LAMP-001", 3)

	report, err := reportUseCase.Execute()
	if err != nil {
		t.Fatalf("expected report to succeed, got %v", err)
	}

	if len(report) != 3 {
		t.Fatalf("expected 3 report rows, got %d", len(report))
	}

	if report[0].SKU != "DESK-001" || report[0].TotalDiscountAmount != 5000 || report[0].QuotedQuantity != 1 {
		t.Fatalf("unexpected first row: %+v", report[0])
	}

	if report[1].SKU != "CHAIR-001" || report[1].TotalDiscountAmount != 3000 || report[1].QuotedQuantity != 3 {
		t.Fatalf("unexpected second row: %+v", report[1])
	}

	if report[2].SKU != "LAMP-001" || report[2].TotalDiscountAmount != 0 || report[2].AverageDiscountRate != 0 {
		t.Fatalf("unexpected third row: %+v", report[2])
	}

	if report[0].AverageDiscountRate != 0.1 {
		t.Fatalf("expected desk average discount rate 0.1, got %f", report[0].AverageDiscountRate)
	}
}

type discountedPricingPolicy struct {
	adjustedBySKU map[string]int
}

func (p discountedPricingPolicy) Price(product domain.Product, quantity int) (int, error) {
	if adjusted, ok := p.adjustedBySKU[product.SKU]; ok {
		return adjusted, nil
	}

	return product.BasePrice, nil
}
