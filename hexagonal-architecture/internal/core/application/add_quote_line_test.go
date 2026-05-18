package application

import (
	"testing"

	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/adapters/services/pricing"
	"hexagonal-architecture/internal/core/domain"
)

func TestAddQuoteLineUsesProductLookupAndPricingPolicy(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	pricingPolicy := pricing.NewFixedPricingPolicy()

	if err := customerRepo.Save(domain.Customer{ID: "customer-001", Active: true}); err != nil {
		t.Fatalf("expected customer save to succeed, got %v", err)
	}

	if err := productRepo.Save(domain.Product{
		SKU:       "CHAIR-001",
		Name:      "Office Chair",
		Category:  "Standard",
		BasePrice: 10000,
		Available: true,
	}); err != nil {
		t.Fatalf("expected product save to succeed, got %v", err)
	}

	createQuote := NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)

	quote, err := createQuote.Execute("customer-001")
	if err != nil {
		t.Fatalf("expected quote creation to succeed, got %v", err)
	}

	updatedQuote, err := addQuoteLine.Execute(quote.ID, "CHAIR-001", 2)
	if err != nil {
		t.Fatalf("expected add quote line to succeed, got %v", err)
	}

	if len(updatedQuote.Lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(updatedQuote.Lines))
	}

	line := updatedQuote.Lines[0]
	if line.SKU != "CHAIR-001" {
		t.Fatalf("expected sku CHAIR-001, got %s", line.SKU)
	}

	if line.BaseUnitPrice != 10000 {
		t.Fatalf("expected base price 10000, got %d", line.BaseUnitPrice)
	}

	if line.AdjustedUnitPrice != 10000 {
		t.Fatalf("expected adjusted price 10000, got %d", line.AdjustedUnitPrice)
	}

	if line.LineTotal != 20000 {
		t.Fatalf("expected line total 20000, got %d", line.LineTotal)
	}
}
