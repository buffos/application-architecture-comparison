package application

import (
	"testing"

	"layered-architecture/internal/domain"
	"layered-architecture/internal/infrastructure/memory"
)

func TestConvertQuoteToOrderReservesStock(t *testing.T) {
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	quoteRepo := memory.NewQuoteRepository()
	stockRepo := memory.NewStockRecordRepository()
	orderRepo := memory.NewOrderRepository()

	customerService := NewCustomerService(customerRepo)
	catalogService := NewCatalogService(productRepo)
	inventoryService := NewInventoryService(productRepo, stockRepo)
	quoteService := NewQuoteService(quoteRepo, customerRepo, productRepo, NoopPricingPluginRegistry{})
	orderService := NewOrderService(orderRepo, quoteRepo, stockRepo)

	customer, err := customerService.CreateCustomer("Acme Corp", "Preferred", "Invoice30")
	if err != nil {
		t.Fatalf("expected customer creation to succeed, got %v", err)
	}

	product, err := catalogService.CreateProduct("CHAIR-001", "Office Chair", "Standard", 10000, true)
	if err != nil {
		t.Fatalf("expected product creation to succeed, got %v", err)
	}

	if _, err := inventoryService.ReceiveStock(product.SKU, 5); err != nil {
		t.Fatalf("expected stock receive to succeed, got %v", err)
	}

	quote, err := quoteService.CreateDraftQuote(customer.ID)
	if err != nil {
		t.Fatalf("expected quote creation to succeed, got %v", err)
	}

	if _, err := quoteService.AddQuoteLine(quote.ID, product.SKU, 2); err != nil {
		t.Fatalf("expected quote line creation to succeed, got %v", err)
	}

	submittedQuote, err := quoteService.SubmitQuote(quote.ID)
	if err != nil {
		t.Fatalf("expected quote submission to succeed, got %v", err)
	}

	if submittedQuote.Status != domain.QuoteStatusApproved {
		t.Fatalf("expected quote status %s, got %s", domain.QuoteStatusApproved, submittedQuote.Status)
	}

	order, err := orderService.ConvertQuoteToOrder(quote.ID)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got %v", err)
	}

	if order.Status != domain.OrderStatusReadyForPayment {
		t.Fatalf("expected order status %s, got %s", domain.OrderStatusReadyForPayment, order.Status)
	}

	stock, err := stockRepo.FindBySKU(product.SKU)
	if err != nil {
		t.Fatalf("expected stock lookup to succeed, got %v", err)
	}

	if stock.Reserved != 2 {
		t.Fatalf("expected reserved quantity 2, got %d", stock.Reserved)
	}
}

func TestCustomBuildQuoteNeedsApprovalBeforeConversion(t *testing.T) {
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	quoteRepo := memory.NewQuoteRepository()
	stockRepo := memory.NewStockRecordRepository()
	orderRepo := memory.NewOrderRepository()

	customerService := NewCustomerService(customerRepo)
	catalogService := NewCatalogService(productRepo)
	inventoryService := NewInventoryService(productRepo, stockRepo)
	quoteService := NewQuoteService(quoteRepo, customerRepo, productRepo, NoopPricingPluginRegistry{})
	orderService := NewOrderService(orderRepo, quoteRepo, stockRepo)

	customer, err := customerService.CreateCustomer("Acme Corp", "Preferred", "Invoice30")
	if err != nil {
		t.Fatalf("expected customer creation to succeed, got %v", err)
	}

	product, err := catalogService.CreateProduct("DESK-001", "Executive Desk", "CustomBuild", 50000, true)
	if err != nil {
		t.Fatalf("expected product creation to succeed, got %v", err)
	}

	if _, err := inventoryService.ReceiveStock(product.SKU, 5); err != nil {
		t.Fatalf("expected stock receive to succeed, got %v", err)
	}

	quote, err := quoteService.CreateDraftQuote(customer.ID)
	if err != nil {
		t.Fatalf("expected quote creation to succeed, got %v", err)
	}

	if _, err := quoteService.AddQuoteLine(quote.ID, product.SKU, 1); err != nil {
		t.Fatalf("expected quote line creation to succeed, got %v", err)
	}

	submittedQuote, err := quoteService.SubmitQuote(quote.ID)
	if err != nil {
		t.Fatalf("expected quote submission to succeed, got %v", err)
	}

	if submittedQuote.Status != domain.QuoteStatusPendingApproval {
		t.Fatalf("expected quote status %s, got %s", domain.QuoteStatusPendingApproval, submittedQuote.Status)
	}

	_, err = orderService.ConvertQuoteToOrder(quote.ID)
	if err != domain.ErrQuoteNotApproved {
		t.Fatalf("expected %v, got %v", domain.ErrQuoteNotApproved, err)
	}

	approvedQuote, err := quoteService.ApproveQuote(quote.ID)
	if err != nil {
		t.Fatalf("expected quote approval to succeed, got %v", err)
	}

	if approvedQuote.Status != domain.QuoteStatusApproved {
		t.Fatalf("expected quote status %s, got %s", domain.QuoteStatusApproved, approvedQuote.Status)
	}

	order, err := orderService.ConvertQuoteToOrder(quote.ID)
	if err != nil {
		t.Fatalf("expected conversion after approval to succeed, got %v", err)
	}

	if order.SourceQuoteID != quote.ID {
		t.Fatalf("expected source quote %s, got %s", quote.ID, order.SourceQuoteID)
	}
}

func TestConvertQuoteToOrderFailsWhenStockIsInsufficient(t *testing.T) {
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	quoteRepo := memory.NewQuoteRepository()
	stockRepo := memory.NewStockRecordRepository()
	orderRepo := memory.NewOrderRepository()

	customerService := NewCustomerService(customerRepo)
	catalogService := NewCatalogService(productRepo)
	inventoryService := NewInventoryService(productRepo, stockRepo)
	quoteService := NewQuoteService(quoteRepo, customerRepo, productRepo, NoopPricingPluginRegistry{})
	orderService := NewOrderService(orderRepo, quoteRepo, stockRepo)

	customer, _ := customerService.CreateCustomer("Acme Corp", "Preferred", "Invoice30")
	product, _ := catalogService.CreateProduct("CHAIR-001", "Office Chair", "Standard", 10000, true)
	_, _ = inventoryService.ReceiveStock(product.SKU, 1)
	quote, _ := quoteService.CreateDraftQuote(customer.ID)
	_, _ = quoteService.AddQuoteLine(quote.ID, product.SKU, 2)
	_, _ = quoteService.SubmitQuote(quote.ID)

	_, err := orderService.ConvertQuoteToOrder(quote.ID)
	if err != domain.ErrInsufficientStock {
		t.Fatalf("expected %v, got %v", domain.ErrInsufficientStock, err)
	}
}
