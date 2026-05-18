package application

import (
	"testing"

	"layered-architecture/internal/domain"
	"layered-architecture/internal/infrastructure/memory"
)

func TestAcceptedReturnRestocksInventory(t *testing.T) {
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	quoteRepo := memory.NewQuoteRepository()
	stockRepo := memory.NewStockRecordRepository()
	orderRepo := memory.NewOrderRepository()
	shipmentRepo := memory.NewShipmentRepository()
	returnRepo := memory.NewReturnRequestRepository()

	customerService := NewCustomerService(customerRepo)
	catalogService := NewCatalogService(productRepo)
	inventoryService := NewInventoryService(productRepo, stockRepo)
	quoteService := NewQuoteService(quoteRepo, customerRepo, productRepo, NoopPricingPluginRegistry{})
	orderService := NewOrderService(orderRepo, quoteRepo, stockRepo)
	paymentService := NewPaymentService(orderRepo)
	fulfillmentService := NewFulfillmentService(orderRepo, stockRepo, shipmentRepo)
	returnService := NewReturnService(orderRepo, stockRepo, returnRepo)

	customer, _ := customerService.CreateCustomer("Acme Corp", "Preferred", "Invoice30")
	product, _ := catalogService.CreateProduct("CHAIR-001", "Office Chair", "Standard", 10000, true)
	_, _ = inventoryService.ReceiveStock(product.SKU, 5)
	quote, _ := quoteService.CreateDraftQuote(customer.ID)
	_, _ = quoteService.AddQuoteLine(quote.ID, product.SKU, 2)
	_, _ = quoteService.SubmitQuote(quote.ID)
	order, _ := orderService.ConvertQuoteToOrder(quote.ID)
	_, _ = paymentService.CapturePayment(order.ID)
	_, _ = fulfillmentService.CreateShipment(order.ID)

	request, err := returnService.RequestReturn(order.ID, "Damaged")
	if err != nil {
		t.Fatalf("expected return request to succeed, got %v", err)
	}

	accepted, err := returnService.AcceptReturn(request.ID)
	if err != nil {
		t.Fatalf("expected return acceptance to succeed, got %v", err)
	}

	if accepted.Status != domain.ReturnStatusAccepted {
		t.Fatalf("expected return status %s, got %s", domain.ReturnStatusAccepted, accepted.Status)
	}

	stock, err := stockRepo.FindBySKU(product.SKU)
	if err != nil {
		t.Fatalf("expected stock lookup to succeed, got %v", err)
	}

	if stock.OnHand != 5 || stock.Reserved != 0 {
		t.Fatalf("expected onHand=5 reserved=0, got onHand=%d reserved=%d", stock.OnHand, stock.Reserved)
	}
}

func TestClearanceItemCannotBeReturned(t *testing.T) {
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	quoteRepo := memory.NewQuoteRepository()
	stockRepo := memory.NewStockRecordRepository()
	orderRepo := memory.NewOrderRepository()
	shipmentRepo := memory.NewShipmentRepository()
	returnRepo := memory.NewReturnRequestRepository()

	customerService := NewCustomerService(customerRepo)
	catalogService := NewCatalogService(productRepo)
	inventoryService := NewInventoryService(productRepo, stockRepo)
	quoteService := NewQuoteService(quoteRepo, customerRepo, productRepo, NoopPricingPluginRegistry{})
	orderService := NewOrderService(orderRepo, quoteRepo, stockRepo)
	paymentService := NewPaymentService(orderRepo)
	fulfillmentService := NewFulfillmentService(orderRepo, stockRepo, shipmentRepo)
	returnService := NewReturnService(orderRepo, stockRepo, returnRepo)

	customer, _ := customerService.CreateCustomer("Acme Corp", "Preferred", "Invoice30")
	product, _ := catalogService.CreateProduct("LAMP-001", "Clearance Lamp", "Clearance", 4000, true)
	_, _ = inventoryService.ReceiveStock(product.SKU, 5)
	quote, _ := quoteService.CreateDraftQuote(customer.ID)
	_, _ = quoteService.AddQuoteLine(quote.ID, product.SKU, 1)
	_, _ = quoteService.SubmitQuote(quote.ID)
	order, _ := orderService.ConvertQuoteToOrder(quote.ID)
	_, _ = paymentService.CapturePayment(order.ID)
	_, _ = fulfillmentService.CreateShipment(order.ID)

	_, err := returnService.RequestReturn(order.ID, "No longer needed")
	if err != domain.ErrReturnNotEligible {
		t.Fatalf("expected %v, got %v", domain.ErrReturnNotEligible, err)
	}
}
