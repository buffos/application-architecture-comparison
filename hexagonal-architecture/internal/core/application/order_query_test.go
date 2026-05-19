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

func TestGetAndListOrders(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	orderRepo := memory.NewOrderRepository()
	shipmentRepo := memory.NewShipmentRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	inventory := memory.NewInventoryReservationAdapter(map[string]int{"CHAIR-001": 5, "DESK-001": 5})
	pricingPolicy := pricing.NewFixedPricingPolicy()
	approvalPolicy := approval.NewCategoryApprovalPolicy()
	paymentGateway := payment.NewAcceptAllGateway()
	shipmentClock := timeadapter.NewFixedClock(time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC))

	_ = customerRepo.Save(domain.Customer{ID: "customer-001", Active: true})
	_ = productRepo.Save(domain.Product{SKU: "CHAIR-001", Name: "Office Chair", Category: "Standard", BasePrice: 10000, Available: true, ReturnWindowDays: 30})
	_ = productRepo.Save(domain.Product{SKU: "DESK-001", Name: "Executive Desk", Category: "Standard", BasePrice: 50000, Available: true, ReturnWindowDays: 30})

	createQuote := NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)
	submitQuote := NewSubmitQuoteUseCase(quoteRepo, approvalPolicy)
	convertQuote := NewConvertQuoteToOrderUseCase(quoteRepo, orderRepo, inventory)
	capturePayment := NewCapturePaymentUseCase(orderRepo, paymentGateway)
	createShipment := NewCreateShipmentUseCase(orderRepo, shipmentRepo, inventory, shipmentClock)
	getOrder := NewGetOrderUseCase(orderRepo)
	listOrders := NewListOrdersUseCase(orderRepo)

	firstQuote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(firstQuote.ID, "CHAIR-001", 2)
	_, _ = submitQuote.Execute(firstQuote.ID)
	firstOrder, _ := convertQuote.Execute(firstQuote.ID)

	secondQuote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(secondQuote.ID, "DESK-001", 1)
	_, _ = submitQuote.Execute(secondQuote.ID)
	secondOrder, _ := convertQuote.Execute(secondQuote.ID)
	_, _ = capturePayment.Execute(secondOrder.ID)
	_, _ = createShipment.Execute(secondOrder.ID)

	loadedOrder, err := getOrder.Execute(firstOrder.ID)
	if err != nil {
		t.Fatalf("expected get order to succeed, got %v", err)
	}

	if loadedOrder.Status != domain.OrderStatusReadyForPayment {
		t.Fatalf("expected loaded order status %s, got %s", domain.OrderStatusReadyForPayment, loadedOrder.Status)
	}

	readyForPayment, err := listOrders.Execute(domain.OrderStatusReadyForPayment)
	if err != nil {
		t.Fatalf("expected list ready-for-payment orders to succeed, got %v", err)
	}

	if len(readyForPayment) != 1 || readyForPayment[0].ID != firstOrder.ID {
		t.Fatalf("expected one ready-for-payment order %s, got %+v", firstOrder.ID, readyForPayment)
	}

	shippedOrders, err := listOrders.Execute(domain.OrderStatusShipped)
	if err != nil {
		t.Fatalf("expected list shipped orders to succeed, got %v", err)
	}

	if len(shippedOrders) != 1 || shippedOrders[0].ID != secondOrder.ID {
		t.Fatalf("expected one shipped order %s, got %+v", secondOrder.ID, shippedOrders)
	}
}
