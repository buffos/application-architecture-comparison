package application

import (
	"testing"

	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/core/domain"
)

func TestCreateDraftQuoteRequiresActiveCustomer(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	customerRepo := memory.NewCustomerRepository()
	useCase := NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)

	if _, err := useCase.Execute("customer-001"); err != domain.ErrCustomerNotFound {
		t.Fatalf("expected %v, got %v", domain.ErrCustomerNotFound, err)
	}

	if err := customerRepo.Save(domain.Customer{ID: "customer-001", Active: false}); err != nil {
		t.Fatalf("expected customer save to succeed, got %v", err)
	}

	if _, err := useCase.Execute("customer-001"); err != domain.ErrCustomerInactive {
		t.Fatalf("expected %v, got %v", domain.ErrCustomerInactive, err)
	}

	if err := customerRepo.Save(domain.Customer{ID: "customer-001", Active: true}); err != nil {
		t.Fatalf("expected customer save to succeed, got %v", err)
	}

	quote, err := useCase.Execute("customer-001")
	if err != nil {
		t.Fatalf("expected quote creation to succeed, got %v", err)
	}

	if quote.CustomerID != "customer-001" {
		t.Fatalf("expected customer id customer-001, got %s", quote.CustomerID)
	}
}
