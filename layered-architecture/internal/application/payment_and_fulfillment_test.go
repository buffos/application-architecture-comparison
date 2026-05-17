package application

import (
	"testing"

	"layered-architecture/internal/domain"
	"layered-architecture/internal/infrastructure/memory"
)

func TestShipmentRequiresAcceptedPayment(t *testing.T) {
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	quoteRepo := memory.NewQuoteRepository()
	stockRepo := memory.NewStockRecordRepository()
	orderRepo := memory.NewOrderRepository()
	shipmentRepo := memory.NewShipmentRepository()

	customerService := NewCustomerService(customerRepo)
	catalogService := NewCatalogService(productRepo)
	inventoryService := NewInventoryService(productRepo, stockRepo)
	quoteService := NewQuoteService(quoteRepo, customerRepo, productRepo)
	orderService := NewOrderService(orderRepo, quoteRepo, stockRepo)
	fulfillmentService := NewFulfillmentService(orderRepo, stockRepo, shipmentRepo)

	customer, _ := customerService.CreateCustomer("Acme Corp", "Preferred", "Invoice30")
	product, _ := catalogService.CreateProduct("CHAIR-001", "Office Chair", "Standard", true)
	_, _ = inventoryService.ReceiveStock(product.SKU, 5)
	quote, _ := quoteService.CreateDraftQuote(customer.ID)
	_, _ = quoteService.AddQuoteLine(quote.ID, product.SKU, 2)
	_, _ = quoteService.SubmitQuote(quote.ID)
	order, _ := orderService.ConvertQuoteToOrder(quote.ID)

	_, err := fulfillmentService.CreateShipment(order.ID)
	if err != domain.ErrShipmentNotAllowedUntilPaymentAccepted {
		t.Fatalf("expected %v, got %v", domain.ErrShipmentNotAllowedUntilPaymentAccepted, err)
	}
}

func TestAcceptedPaymentAllowsShipmentAndConsumesReservedStock(t *testing.T) {
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	quoteRepo := memory.NewQuoteRepository()
	stockRepo := memory.NewStockRecordRepository()
	orderRepo := memory.NewOrderRepository()
	shipmentRepo := memory.NewShipmentRepository()

	customerService := NewCustomerService(customerRepo)
	catalogService := NewCatalogService(productRepo)
	inventoryService := NewInventoryService(productRepo, stockRepo)
	quoteService := NewQuoteService(quoteRepo, customerRepo, productRepo)
	orderService := NewOrderService(orderRepo, quoteRepo, stockRepo)
	paymentService := NewPaymentService(orderRepo)
	fulfillmentService := NewFulfillmentService(orderRepo, stockRepo, shipmentRepo)

	customer, _ := customerService.CreateCustomer("Acme Corp", "Preferred", "Invoice30")
	product, _ := catalogService.CreateProduct("CHAIR-001", "Office Chair", "Standard", true)
	_, _ = inventoryService.ReceiveStock(product.SKU, 5)
	quote, _ := quoteService.CreateDraftQuote(customer.ID)
	_, _ = quoteService.AddQuoteLine(quote.ID, product.SKU, 2)
	_, _ = quoteService.SubmitQuote(quote.ID)
	order, _ := orderService.ConvertQuoteToOrder(quote.ID)

	paidOrder, err := paymentService.CapturePayment(order.ID)
	if err != nil {
		t.Fatalf("expected payment capture to succeed, got %v", err)
	}

	if paidOrder.Status != domain.OrderStatusReadyForFulfillment {
		t.Fatalf("expected order status %s, got %s", domain.OrderStatusReadyForFulfillment, paidOrder.Status)
	}

	shipment, err := fulfillmentService.CreateShipment(order.ID)
	if err != nil {
		t.Fatalf("expected shipment creation to succeed, got %v", err)
	}

	if shipment.Status != domain.ShipmentStatusShipped {
		t.Fatalf("expected shipment status %s, got %s", domain.ShipmentStatusShipped, shipment.Status)
	}

	loadedOrder, err := orderRepo.FindByID(order.ID)
	if err != nil {
		t.Fatalf("expected order lookup to succeed, got %v", err)
	}

	if loadedOrder.Status != domain.OrderStatusShipped {
		t.Fatalf("expected order status %s, got %s", domain.OrderStatusShipped, loadedOrder.Status)
	}

	stock, err := stockRepo.FindBySKU(product.SKU)
	if err != nil {
		t.Fatalf("expected stock lookup to succeed, got %v", err)
	}

	if stock.OnHand != 3 || stock.Reserved != 0 {
		t.Fatalf("expected onHand=3 reserved=0, got onHand=%d reserved=%d", stock.OnHand, stock.Reserved)
	}
}
