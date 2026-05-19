package application

import (
	"testing"
	"time"

	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/adapters/services/approval"
	"hexagonal-architecture/internal/adapters/services/payment"
	"hexagonal-architecture/internal/adapters/services/pricing"
	timeadapter "hexagonal-architecture/internal/adapters/services/time"
	"hexagonal-architecture/internal/core/domain"
)

func TestPartialShipmentUpdatesOrderStateAndRemainingQuantities(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	orderRepo := memory.NewOrderRepository()
	shipmentRepo := memory.NewShipmentRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	inventory := memory.NewInventoryReservationAdapter(map[string]int{"CHAIR-001": 5})
	pricingPolicy := pricing.NewFixedPricingPolicy()
	approvalPolicy := approval.NewCategoryApprovalPolicy()
	paymentGateway := payment.NewAcceptAllGateway()
	shipmentClock := timeadapter.NewFixedClock(time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC))

	_ = customerRepo.Save(domain.Customer{ID: "customer-001", Active: true})
	_ = productRepo.Save(domain.Product{SKU: "CHAIR-001", Name: "Office Chair", Category: "Standard", BasePrice: 10000, Available: true, ReturnWindowDays: 30})

	createQuote := NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)
	submitQuote := NewSubmitQuoteUseCase(quoteRepo, approvalPolicy)
	convertQuote := NewConvertQuoteToOrderUseCase(quoteRepo, orderRepo, inventory)
	capturePayment := NewCapturePaymentUseCase(orderRepo, paymentGateway)
	createShipment := NewCreateShipmentUseCase(orderRepo, shipmentRepo, inventory, shipmentClock)

	quote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(quote.ID, "CHAIR-001", 3)
	_, _ = submitQuote.Execute(quote.ID)
	order, _ := convertQuote.Execute(quote.ID)
	_, _ = capturePayment.Execute(order.ID)

	firstShipment, err := createShipment.Execute(order.ID, domain.ShipmentLine{SKU: "CHAIR-001", Quantity: 1})
	if err != nil {
		t.Fatalf("expected partial shipment to succeed, got %v", err)
	}

	if firstShipment.Status != domain.ShipmentStatusPartiallyShipped {
		t.Fatalf("expected shipment status %s, got %s", domain.ShipmentStatusPartiallyShipped, firstShipment.Status)
	}

	loadedOrder, err := orderRepo.FindByID(order.ID)
	if err != nil {
		t.Fatalf("expected order lookup to succeed, got %v", err)
	}

	if loadedOrder.Status != domain.OrderStatusPartiallyShipped {
		t.Fatalf("expected order status %s, got %s", domain.OrderStatusPartiallyShipped, loadedOrder.Status)
	}

	if loadedOrder.Lines[0].ShippedQuantity != 1 {
		t.Fatalf("expected shipped quantity 1, got %d", loadedOrder.Lines[0].ShippedQuantity)
	}

	secondShipment, err := createShipment.Execute(order.ID)
	if err != nil {
		t.Fatalf("expected final shipment to succeed, got %v", err)
	}

	if secondShipment.Status != domain.ShipmentStatusShipped {
		t.Fatalf("expected shipment status %s, got %s", domain.ShipmentStatusShipped, secondShipment.Status)
	}

	loadedOrder, err = orderRepo.FindByID(order.ID)
	if err != nil {
		t.Fatalf("expected order lookup to succeed, got %v", err)
	}

	if loadedOrder.Status != domain.OrderStatusShipped || loadedOrder.Lines[0].ShippedQuantity != 3 {
		t.Fatalf("unexpected final order state: %+v", loadedOrder)
	}
}
