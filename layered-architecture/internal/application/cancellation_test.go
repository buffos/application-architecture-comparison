package application

import (
	"testing"

	"layered-architecture/internal/domain"
	"layered-architecture/internal/infrastructure/memory"
)

func TestCancelOrderReleasesReservedStock(t *testing.T) {
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
	_, _ = inventoryService.ReceiveStock(product.SKU, 5)
	quote, _ := quoteService.CreateDraftQuote(customer.ID)
	_, _ = quoteService.AddQuoteLine(quote.ID, product.SKU, 2)
	_, _ = quoteService.SubmitQuote(quote.ID)
	order, _ := orderService.ConvertQuoteToOrder(quote.ID)

	cancelledOrder, err := orderService.CancelOrder(order.ID)
	if err != nil {
		t.Fatalf("expected cancel to succeed, got %v", err)
	}

	if cancelledOrder.Status != domain.OrderStatusCancelled {
		t.Fatalf("expected order status %s, got %s", domain.OrderStatusCancelled, cancelledOrder.Status)
	}

	stock, err := stockRepo.FindBySKU(product.SKU)
	if err != nil {
		t.Fatalf("expected stock lookup to succeed, got %v", err)
	}

	if stock.Reserved != 0 || stock.OnHand != 5 {
		t.Fatalf("expected onHand=5 reserved=0, got onHand=%d reserved=%d", stock.OnHand, stock.Reserved)
	}
}

func TestCancelShippedOrderFails(t *testing.T) {
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	quoteRepo := memory.NewQuoteRepository()
	stockRepo := memory.NewStockRecordRepository()
	orderRepo := memory.NewOrderRepository()
	shipmentRepo := memory.NewShipmentRepository()

	customerService := NewCustomerService(customerRepo)
	catalogService := NewCatalogService(productRepo)
	inventoryService := NewInventoryService(productRepo, stockRepo)
	quoteService := NewQuoteService(quoteRepo, customerRepo, productRepo, NoopPricingPluginRegistry{})
	orderService := NewOrderService(orderRepo, quoteRepo, stockRepo)
	paymentService := NewPaymentService(orderRepo)
	fulfillmentService := NewFulfillmentService(orderRepo, stockRepo, shipmentRepo)

	customer, _ := customerService.CreateCustomer("Acme Corp", "Preferred", "Invoice30")
	product, _ := catalogService.CreateProduct("CHAIR-001", "Office Chair", "Standard", 10000, true)
	_, _ = inventoryService.ReceiveStock(product.SKU, 5)
	quote, _ := quoteService.CreateDraftQuote(customer.ID)
	_, _ = quoteService.AddQuoteLine(quote.ID, product.SKU, 2)
	_, _ = quoteService.SubmitQuote(quote.ID)
	order, _ := orderService.ConvertQuoteToOrder(quote.ID)
	_, _ = paymentService.CapturePayment(order.ID)
	_, _ = fulfillmentService.CreateShipment(order.ID)

	_, err := orderService.CancelOrder(order.ID)
	if err != domain.ErrOrderAlreadyShipped {
		t.Fatalf("expected %v, got %v", domain.ErrOrderAlreadyShipped, err)
	}
}
