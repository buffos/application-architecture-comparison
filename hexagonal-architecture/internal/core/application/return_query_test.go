package application

import (
	"testing"
	"time"

	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/adapters/services/approval"
	"hexagonal-architecture/internal/adapters/services/payment"
	"hexagonal-architecture/internal/adapters/services/pricing"
	"hexagonal-architecture/internal/adapters/services/refund"
	"hexagonal-architecture/internal/adapters/services/returnpolicy"
	timeadapter "hexagonal-architecture/internal/adapters/services/time"
	"hexagonal-architecture/internal/core/domain"
)

func TestGetAndListReturnRequests(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	orderRepo := memory.NewOrderRepository()
	shipmentRepo := memory.NewShipmentRepository()
	returnRepo := memory.NewReturnRequestRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	inventory := memory.NewInventoryReservationAdapter(map[string]int{"CHAIR-001": 5, "DESK-001": 5})
	idempotency := memory.NewIdempotencyStore()
	pricingPolicy := pricing.NewFixedPricingPolicy()
	approvalPolicy := approval.NewCategoryApprovalPolicy()
	paymentGateway := payment.NewAcceptAllGateway()
	refundGateway := refund.NewAcceptAllGateway()
	returnPolicy := returnpolicy.NewWindowPolicy()
	shipmentClock := timeadapter.NewFixedClock(time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC))
	returnClock := timeadapter.NewFixedClock(time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC))

	_ = customerRepo.Save(domain.Customer{ID: "customer-001", Active: true})
	_ = productRepo.Save(domain.Product{SKU: "CHAIR-001", Name: "Office Chair", Category: "Standard", BasePrice: 10000, Available: true, ReturnWindowDays: 30})
	_ = productRepo.Save(domain.Product{SKU: "DESK-001", Name: "Executive Desk", Category: "Standard", BasePrice: 50000, Available: true, ReturnWindowDays: 30})

	createQuote := NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)
	submitQuote := NewSubmitQuoteUseCase(quoteRepo, approvalPolicy)
	convertQuote := NewConvertQuoteToOrderUseCase(quoteRepo, orderRepo, inventory)
	capturePayment := NewCapturePaymentUseCase(orderRepo, paymentGateway)
	createShipment := NewCreateShipmentUseCase(orderRepo, shipmentRepo, inventory, shipmentClock)
	requestReturn := NewRequestReturnUseCase(orderRepo, returnRepo, returnClock)
	acceptReturn := NewAcceptReturnUseCase(orderRepo, returnRepo, returnPolicy, idempotency)
	completeRefund := NewCompleteRefundUseCase(returnRepo, refundGateway, inventory, idempotency)
	getReturn := NewGetReturnRequestUseCase(returnRepo)
	listReturns := NewListReturnRequestsUseCase(returnRepo)

	firstQuote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(firstQuote.ID, "CHAIR-001", 2)
	_, _ = submitQuote.Execute(firstQuote.ID)
	firstOrder, _ := convertQuote.Execute(firstQuote.ID)
	_, _ = capturePayment.Execute(firstOrder.ID)
	_, _ = createShipment.Execute(firstOrder.ID)
	firstReturn, _ := requestReturn.Execute(firstOrder.ID, "Damaged", "warehouse-clerk-1")
	_, _ = acceptReturn.Execute(firstReturn.ID, "warehouse-clerk-1", "return-accept-101")
	_, _ = completeRefund.Execute(firstReturn.ID, "manager-1", "refund-101")

	secondQuote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(secondQuote.ID, "DESK-001", 1)
	_, _ = submitQuote.Execute(secondQuote.ID)
	secondOrder, _ := convertQuote.Execute(secondQuote.ID)
	_, _ = capturePayment.Execute(secondOrder.ID)
	_, _ = createShipment.Execute(secondOrder.ID)
	secondReturn, _ := requestReturn.Execute(secondOrder.ID, "Changed mind", "warehouse-clerk-2")

	loadedReturn, err := getReturn.Execute(firstReturn.ID)
	if err != nil {
		t.Fatalf("expected get return request to succeed, got %v", err)
	}

	if loadedReturn.Status != domain.ReturnStatusRefunded {
		t.Fatalf("expected loaded return status %s, got %s", domain.ReturnStatusRefunded, loadedReturn.Status)
	}

	requestedReturns, err := listReturns.Execute(domain.ReturnStatusRequested)
	if err != nil {
		t.Fatalf("expected list requested returns to succeed, got %v", err)
	}

	if len(requestedReturns) != 1 || requestedReturns[0].ID != secondReturn.ID {
		t.Fatalf("expected one requested return %s, got %+v", secondReturn.ID, requestedReturns)
	}

	refundedReturns, err := listReturns.Execute(domain.ReturnStatusRefunded)
	if err != nil {
		t.Fatalf("expected list refunded returns to succeed, got %v", err)
	}

	if len(refundedReturns) != 1 || refundedReturns[0].ID != firstReturn.ID {
		t.Fatalf("expected one refunded return %s, got %+v", firstReturn.ID, refundedReturns)
	}
}
