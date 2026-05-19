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

func TestPartialReturnTracksRemainingReturnableQuantity(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	orderRepo := memory.NewOrderRepository()
	shipmentRepo := memory.NewShipmentRepository()
	returnRepo := memory.NewReturnRequestRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	inventory := memory.NewInventoryReservationAdapter(map[string]int{"CHAIR-001": 5})
	pricingPolicy := pricing.NewFixedPricingPolicy()
	approvalPolicy := approval.NewCategoryApprovalPolicy()
	paymentGateway := payment.NewAcceptAllGateway()
	refundGateway := refund.NewAcceptAllGateway()
	returnPolicy := returnpolicy.NewWindowPolicy()
	shipmentClock := timeadapter.NewFixedClock(time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC))
	returnClock := timeadapter.NewFixedClock(time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC))
	idempotency := memory.NewIdempotencyStore()

	_ = customerRepo.Save(domain.Customer{ID: "customer-001", Active: true})
	_ = productRepo.Save(domain.Product{SKU: "CHAIR-001", Name: "Office Chair", Category: "Standard", BasePrice: 10000, Available: true, ReturnWindowDays: 30})

	createQuote := NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)
	submitQuote := NewSubmitQuoteUseCase(quoteRepo, approvalPolicy)
	convertQuote := NewConvertQuoteToOrderUseCase(quoteRepo, orderRepo, inventory)
	capturePayment := NewCapturePaymentUseCase(orderRepo, paymentGateway)
	createShipment := NewCreateShipmentUseCase(orderRepo, shipmentRepo, inventory, shipmentClock)
	requestReturn := NewRequestReturnUseCase(orderRepo, returnRepo, returnClock)
	acceptReturn := NewAcceptReturnUseCase(orderRepo, returnRepo, returnPolicy, idempotency)
	completeRefund := NewCompleteRefundUseCase(returnRepo, refundGateway, inventory, idempotency)

	quote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(quote.ID, "CHAIR-001", 3)
	_, _ = submitQuote.Execute(quote.ID)
	order, _ := convertQuote.Execute(quote.ID)
	_, _ = capturePayment.Execute(order.ID)
	_, _ = createShipment.Execute(order.ID, domain.ShipmentLine{SKU: "CHAIR-001", Quantity: 2})

	request, err := requestReturn.Execute(order.ID, "Damaged", "warehouse-clerk-1", domain.ReturnLineRequest{SKU: "CHAIR-001", Quantity: 1})
	if err != nil {
		t.Fatalf("expected partial return request to succeed, got %v", err)
	}

	if len(request.Lines) != 1 || request.Lines[0].Quantity != 1 {
		t.Fatalf("unexpected partial return lines: %+v", request.Lines)
	}

	_, err = acceptReturn.Execute(request.ID, "warehouse-clerk-1", "return-accept-031")
	if err != nil {
		t.Fatalf("expected partial return acceptance to succeed, got %v", err)
	}

	loadedOrder, err := orderRepo.FindByID(order.ID)
	if err != nil {
		t.Fatalf("expected order lookup to succeed, got %v", err)
	}

	if loadedOrder.Lines[0].ReturnedQuantity != 1 {
		t.Fatalf("expected returned quantity 1, got %d", loadedOrder.Lines[0].ReturnedQuantity)
	}

	_, err = requestReturn.Execute(order.ID, "Damaged again", "warehouse-clerk-1", domain.ReturnLineRequest{SKU: "CHAIR-001", Quantity: 2})
	if err != domain.ErrReturnQuantityExceedsRemaining {
		t.Fatalf("expected %v, got %v", domain.ErrReturnQuantityExceedsRemaining, err)
	}

	_, err = completeRefund.Execute(request.ID, "manager-1", "refund-031")
	if err != nil {
		t.Fatalf("expected partial refund to succeed, got %v", err)
	}

	if inventory.Available("CHAIR-001") != 3 {
		t.Fatalf("expected stock restocked to 3, got %d", inventory.Available("CHAIR-001"))
	}
}
