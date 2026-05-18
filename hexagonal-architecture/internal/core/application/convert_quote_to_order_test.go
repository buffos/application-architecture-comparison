package application

import (
	"testing"

	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/adapters/services/approval"
	"hexagonal-architecture/internal/adapters/services/pricing"
	"hexagonal-architecture/internal/core/domain"
)

func TestConvertQuoteToOrderReservesInventoryAndCreatesOrder(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	orderRepo := memory.NewOrderRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	inventory := memory.NewInventoryReservationAdapter(map[string]int{
		"CHAIR-001": 5,
	})
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

	createQuote := NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)
	submitQuote := NewSubmitQuoteUseCase(quoteRepo, approvalPolicy)
	convertQuote := NewConvertQuoteToOrderUseCase(quoteRepo, orderRepo, inventory)

	quote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(quote.ID, "CHAIR-001", 2)
	_, _ = submitQuote.Execute(quote.ID)

	order, err := convertQuote.Execute(quote.ID)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got %v", err)
	}

	if order.Status != domain.OrderStatusReadyForPayment {
		t.Fatalf("expected order status %s, got %s", domain.OrderStatusReadyForPayment, order.Status)
	}

	if inventory.Available("CHAIR-001") != 3 {
		t.Fatalf("expected remaining stock 3, got %d", inventory.Available("CHAIR-001"))
	}
}

func TestConvertQuoteToOrderRequiresApprovedQuote(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	orderRepo := memory.NewOrderRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	inventory := memory.NewInventoryReservationAdapter(map[string]int{
		"DESK-001": 5,
	})
	pricingPolicy := pricing.NewFixedPricingPolicy()

	_ = customerRepo.Save(domain.Customer{ID: "customer-001", Active: true})
	_ = productRepo.Save(domain.Product{
		SKU:       "DESK-001",
		Name:      "Executive Desk",
		Category:  "CustomBuild",
		BasePrice: 50000,
		Available: true,
	})

	createQuote := NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)
	convertQuote := NewConvertQuoteToOrderUseCase(quoteRepo, orderRepo, inventory)

	quote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(quote.ID, "DESK-001", 1)

	_, err := convertQuote.Execute(quote.ID)
	if err != domain.ErrQuoteNotApproved {
		t.Fatalf("expected %v, got %v", domain.ErrQuoteNotApproved, err)
	}
}
