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

func TestPaymentAndShipmentWorkflowUsesPorts(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	orderRepo := memory.NewOrderRepository()
	shipmentRepo := memory.NewShipmentRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	inventory := memory.NewInventoryReservationAdapter(map[string]int{
		"CHAIR-001": 5,
	})
	pricingPolicy := pricing.NewFixedPricingPolicy()
	approvalPolicy := approval.NewCategoryApprovalPolicy()
	paymentGateway := payment.NewAcceptAllGateway()
	shipmentClock := timeadapter.NewFixedClock(time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC))

	_ = customerRepo.Save(domain.Customer{ID: "customer-001", Active: true})
	_ = productRepo.Save(domain.Product{
		SKU:              "CHAIR-001",
		Name:             "Office Chair",
		Category:         "Standard",
		BasePrice:        10000,
		Available:        true,
		ReturnWindowDays: 30,
	})

	createQuote := NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)
	submitQuote := NewSubmitQuoteUseCase(quoteRepo, approvalPolicy)
	convertQuote := NewConvertQuoteToOrderUseCase(quoteRepo, orderRepo, inventory)
	capturePayment := NewCapturePaymentUseCase(orderRepo, paymentGateway)
	createShipment := NewCreateShipmentUseCase(orderRepo, shipmentRepo, inventory, shipmentClock)

	quote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(quote.ID, "CHAIR-001", 2)
	_, _ = submitQuote.Execute(quote.ID)
	order, _ := convertQuote.Execute(quote.ID)

	paidOrder, err := capturePayment.Execute(order.ID)
	if err != nil {
		t.Fatalf("expected payment capture to succeed, got %v", err)
	}

	if paidOrder.Status != domain.OrderStatusReadyForFulfillment {
		t.Fatalf("expected order status %s, got %s", domain.OrderStatusReadyForFulfillment, paidOrder.Status)
	}

	shipment, err := createShipment.Execute(order.ID)
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

	if !loadedOrder.ShippedAt.Equal(shipmentClock.Now()) {
		t.Fatalf("expected shippedAt %v, got %v", shipmentClock.Now(), loadedOrder.ShippedAt)
	}

	if inventory.Available("CHAIR-001") != 1 {
		t.Fatalf("expected remaining stock 1, got %d", inventory.Available("CHAIR-001"))
	}
}
