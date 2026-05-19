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

func TestGetAndListShipments(t *testing.T) {
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
	getShipment := NewGetShipmentUseCase(shipmentRepo)
	listShipments := NewListShipmentsUseCase(shipmentRepo)

	firstQuote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(firstQuote.ID, "CHAIR-001", 2)
	_, _ = submitQuote.Execute(firstQuote.ID)
	firstOrder, _ := convertQuote.Execute(firstQuote.ID)
	_, _ = capturePayment.Execute(firstOrder.ID)
	firstShipment, _ := createShipment.Execute(firstOrder.ID)

	secondQuote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(secondQuote.ID, "DESK-001", 1)
	_, _ = submitQuote.Execute(secondQuote.ID)
	secondOrder, _ := convertQuote.Execute(secondQuote.ID)
	_, _ = capturePayment.Execute(secondOrder.ID)
	secondShipment, _ := createShipment.Execute(secondOrder.ID)

	loadedShipment, err := getShipment.Execute(firstShipment.ID)
	if err != nil {
		t.Fatalf("expected get shipment to succeed, got %v", err)
	}

	if loadedShipment.OrderID != firstOrder.ID {
		t.Fatalf("expected loaded shipment order id %s, got %s", firstOrder.ID, loadedShipment.OrderID)
	}

	firstOrderShipments, err := listShipments.Execute(firstOrder.ID)
	if err != nil {
		t.Fatalf("expected list shipments by order to succeed, got %v", err)
	}

	if len(firstOrderShipments) != 1 || firstOrderShipments[0].ID != firstShipment.ID {
		t.Fatalf("expected one shipment %s for first order, got %+v", firstShipment.ID, firstOrderShipments)
	}

	secondOrderShipments, err := listShipments.Execute(secondOrder.ID)
	if err != nil {
		t.Fatalf("expected list shipments by order to succeed, got %v", err)
	}

	if len(secondOrderShipments) != 1 || secondOrderShipments[0].ID != secondShipment.ID {
		t.Fatalf("expected one shipment %s for second order, got %+v", secondShipment.ID, secondOrderShipments)
	}
}
