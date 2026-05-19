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

func TestShippedOrderCanRequestReturnAndBeRefunded(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	orderRepo := memory.NewOrderRepository()
	shipmentRepo := memory.NewShipmentRepository()
	returnRepo := memory.NewReturnRequestRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	inventory := memory.NewInventoryReservationAdapter(map[string]int{
		"CHAIR-001": 5,
	})
	pricingPolicy := pricing.NewFixedPricingPolicy()
	approvalPolicy := approval.NewCategoryApprovalPolicy()
	paymentGateway := payment.NewAcceptAllGateway()
	refundGateway := refund.NewAcceptAllGateway()
	returnPolicy := returnpolicy.NewWindowPolicy()
	shipmentClock := timeadapter.NewFixedClock(time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC))
	returnClock := timeadapter.NewFixedClock(time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC))
	idempotency := memory.NewIdempotencyStore()

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
	requestReturn := NewRequestReturnUseCase(orderRepo, returnRepo, returnClock)
	acceptReturn := NewAcceptReturnUseCase(returnRepo, returnPolicy, idempotency)
	completeRefund := NewCompleteRefundUseCase(returnRepo, refundGateway, inventory, idempotency)

	quote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(quote.ID, "CHAIR-001", 2)
	_, _ = submitQuote.Execute(quote.ID)
	order, _ := convertQuote.Execute(quote.ID)
	_, _ = capturePayment.Execute(order.ID)
	_, _ = createShipment.Execute(order.ID)

	request, err := requestReturn.Execute(order.ID, "Damaged", "warehouse-clerk-1")
	if err != nil {
		t.Fatalf("expected return request to succeed, got %v", err)
	}

	if request.Status != domain.ReturnStatusRequested {
		t.Fatalf("expected return status %s, got %s", domain.ReturnStatusRequested, request.Status)
	}

	if request.RequestedBy != "warehouse-clerk-1" {
		t.Fatalf("expected requestedBy warehouse-clerk-1, got %s", request.RequestedBy)
	}

	accepted, err := acceptReturn.Execute(request.ID, "warehouse-clerk-1", "return-accept-001")
	if err != nil {
		t.Fatalf("expected return acceptance to succeed, got %v", err)
	}

	if accepted.Status != domain.ReturnStatusAccepted {
		t.Fatalf("expected return status %s, got %s", domain.ReturnStatusAccepted, accepted.Status)
	}

	if accepted.ReviewedBy != "warehouse-clerk-1" {
		t.Fatalf("expected reviewedBy warehouse-clerk-1, got %s", accepted.ReviewedBy)
	}

	refunded, err := completeRefund.Execute(request.ID, "manager-1", "refund-001")
	if err != nil {
		t.Fatalf("expected refund to succeed, got %v", err)
	}

	if refunded.Status != domain.ReturnStatusRefunded {
		t.Fatalf("expected return status %s, got %s", domain.ReturnStatusRefunded, refunded.Status)
	}

	if refunded.ProcessedBy != "manager-1" {
		t.Fatalf("expected processedBy manager-1, got %s", refunded.ProcessedBy)
	}

	if inventory.Available("CHAIR-001") != 3 {
		t.Fatalf("expected restocked stock 3, got %d", inventory.Available("CHAIR-001"))
	}
}

func TestReturnCanBeRejectedAndThenCannotBeRefunded(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	orderRepo := memory.NewOrderRepository()
	shipmentRepo := memory.NewShipmentRepository()
	returnRepo := memory.NewReturnRequestRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	inventory := memory.NewInventoryReservationAdapter(map[string]int{
		"CHAIR-001": 5,
	})
	pricingPolicy := pricing.NewFixedPricingPolicy()
	approvalPolicy := approval.NewCategoryApprovalPolicy()
	paymentGateway := payment.NewAcceptAllGateway()
	refundGateway := refund.NewAcceptAllGateway()
	returnPolicy := returnpolicy.NewWindowPolicy()
	shipmentClock := timeadapter.NewFixedClock(time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC))
	returnClock := timeadapter.NewFixedClock(time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC))
	idempotency := memory.NewIdempotencyStore()

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
	requestReturn := NewRequestReturnUseCase(orderRepo, returnRepo, returnClock)
	acceptReturn := NewAcceptReturnUseCase(returnRepo, returnPolicy, idempotency)
	rejectReturn := NewRejectReturnUseCase(returnRepo)
	completeRefund := NewCompleteRefundUseCase(returnRepo, refundGateway, inventory, idempotency)

	quote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(quote.ID, "CHAIR-001", 2)
	_, _ = submitQuote.Execute(quote.ID)
	order, _ := convertQuote.Execute(quote.ID)
	_, _ = capturePayment.Execute(order.ID)
	_, _ = createShipment.Execute(order.ID)

	request, err := requestReturn.Execute(order.ID, "Changed mind", "warehouse-clerk-1")
	if err != nil {
		t.Fatalf("expected return request to succeed, got %v", err)
	}

	rejected, err := rejectReturn.Execute(request.ID, "warehouse-clerk-2", "Outside policy")
	if err != nil {
		t.Fatalf("expected return rejection to succeed, got %v", err)
	}

	if rejected.Status != domain.ReturnStatusRejected {
		t.Fatalf("expected return status %s, got %s", domain.ReturnStatusRejected, rejected.Status)
	}

	if rejected.ReviewedBy != "warehouse-clerk-2" {
		t.Fatalf("expected reviewedBy warehouse-clerk-2, got %s", rejected.ReviewedBy)
	}

	if rejected.ReviewNote != "Outside policy" {
		t.Fatalf("expected review note Outside policy, got %s", rejected.ReviewNote)
	}

	_, err = acceptReturn.Execute(request.ID, "warehouse-clerk-3", "return-accept-002")
	if err != domain.ErrReturnNotEligible && err != domain.ErrReturnReviewNotAllowed {
		t.Fatalf("expected review denial or already-reviewed error, got %v", err)
	}

	_, err = completeRefund.Execute(request.ID, "manager-1", "refund-002")
	if err != domain.ErrReturnRefundNotAllowed {
		t.Fatalf("expected %v, got %v", domain.ErrReturnRefundNotAllowed, err)
	}

	if inventory.Available("CHAIR-001") != 1 {
		t.Fatalf("expected stock to remain 1 after rejected return, got %d", inventory.Available("CHAIR-001"))
	}
}

func TestReturnAcceptanceCanBeBlockedByPolicy(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	orderRepo := memory.NewOrderRepository()
	shipmentRepo := memory.NewShipmentRepository()
	returnRepo := memory.NewReturnRequestRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	inventory := memory.NewInventoryReservationAdapter(map[string]int{
		"CHAIR-001": 5,
	})
	pricingPolicy := pricing.NewFixedPricingPolicy()
	approvalPolicy := approval.NewCategoryApprovalPolicy()
	paymentGateway := payment.NewAcceptAllGateway()
	returnPolicy := returnpolicy.NewWindowPolicy()
	shipmentClock := timeadapter.NewFixedClock(time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC))
	returnClock := timeadapter.NewFixedClock(time.Date(2026, 6, 5, 9, 0, 0, 0, time.UTC))
	idempotency := memory.NewIdempotencyStore()

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
	requestReturn := NewRequestReturnUseCase(orderRepo, returnRepo, returnClock)
	acceptReturn := NewAcceptReturnUseCase(returnRepo, returnPolicy, idempotency)

	quote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(quote.ID, "CHAIR-001", 2)
	_, _ = submitQuote.Execute(quote.ID)
	order, _ := convertQuote.Execute(quote.ID)
	_, _ = capturePayment.Execute(order.ID)
	_, _ = createShipment.Execute(order.ID)

	request, err := requestReturn.Execute(order.ID, "Changed mind", "warehouse-clerk-1")
	if err != nil {
		t.Fatalf("expected return request to succeed, got %v", err)
	}

	_, err = acceptReturn.Execute(request.ID, "warehouse-clerk-1", "return-accept-003")
	if err != domain.ErrReturnNotEligible {
		t.Fatalf("expected %v, got %v", domain.ErrReturnNotEligible, err)
	}

	storedRequest, err := returnRepo.FindByID(request.ID)
	if err != nil {
		t.Fatalf("expected return lookup to succeed, got %v", err)
	}

	if storedRequest.Status != domain.ReturnStatusRequested {
		t.Fatalf("expected return status %s, got %s", domain.ReturnStatusRequested, storedRequest.Status)
	}
}

func TestReturnRequestRequiresShippedOrder(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	orderRepo := memory.NewOrderRepository()
	returnRepo := memory.NewReturnRequestRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	inventory := memory.NewInventoryReservationAdapter(map[string]int{
		"CHAIR-001": 5,
	})
	pricingPolicy := pricing.NewFixedPricingPolicy()
	approvalPolicy := approval.NewCategoryApprovalPolicy()
	returnClock := timeadapter.NewFixedClock(time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC))

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
	requestReturn := NewRequestReturnUseCase(orderRepo, returnRepo, returnClock)

	quote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(quote.ID, "CHAIR-001", 2)
	_, _ = submitQuote.Execute(quote.ID)
	order, _ := convertQuote.Execute(quote.ID)

	_, err := requestReturn.Execute(order.ID, "Changed mind", "warehouse-clerk-1")
	if err != domain.ErrReturnNotEligible {
		t.Fatalf("expected %v, got %v", domain.ErrReturnNotEligible, err)
	}
}

func TestClearanceItemCannotBeReturned(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	orderRepo := memory.NewOrderRepository()
	shipmentRepo := memory.NewShipmentRepository()
	returnRepo := memory.NewReturnRequestRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	inventory := memory.NewInventoryReservationAdapter(map[string]int{
		"LAMP-001": 5,
	})
	pricingPolicy := pricing.NewFixedPricingPolicy()
	approvalPolicy := approval.NewCategoryApprovalPolicy()
	paymentGateway := payment.NewAcceptAllGateway()
	shipmentClock := timeadapter.NewFixedClock(time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC))
	returnClock := timeadapter.NewFixedClock(time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC))

	_ = customerRepo.Save(domain.Customer{ID: "customer-001", Active: true})
	_ = productRepo.Save(domain.Product{
		SKU:              "LAMP-001",
		Name:             "Clearance Lamp",
		Category:         "Clearance",
		BasePrice:        4000,
		Available:        true,
		ReturnWindowDays: 30,
	})

	createQuote := NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)
	submitQuote := NewSubmitQuoteUseCase(quoteRepo, approvalPolicy)
	convertQuote := NewConvertQuoteToOrderUseCase(quoteRepo, orderRepo, inventory)
	capturePayment := NewCapturePaymentUseCase(orderRepo, paymentGateway)
	createShipment := NewCreateShipmentUseCase(orderRepo, shipmentRepo, inventory, shipmentClock)
	requestReturn := NewRequestReturnUseCase(orderRepo, returnRepo, returnClock)

	quote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(quote.ID, "LAMP-001", 1)
	_, _ = submitQuote.Execute(quote.ID)
	order, _ := convertQuote.Execute(quote.ID)
	_, _ = capturePayment.Execute(order.ID)
	_, _ = createShipment.Execute(order.ID)

	_, err := requestReturn.Execute(order.ID, "No longer needed", "warehouse-clerk-1")
	if err != domain.ErrReturnNotEligible {
		t.Fatalf("expected %v, got %v", domain.ErrReturnNotEligible, err)
	}
}

func TestReturnActorsAreRequired(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	orderRepo := memory.NewOrderRepository()
	shipmentRepo := memory.NewShipmentRepository()
	returnRepo := memory.NewReturnRequestRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	inventory := memory.NewInventoryReservationAdapter(map[string]int{
		"CHAIR-001": 5,
	})
	pricingPolicy := pricing.NewFixedPricingPolicy()
	approvalPolicy := approval.NewCategoryApprovalPolicy()
	paymentGateway := payment.NewAcceptAllGateway()
	returnPolicy := returnpolicy.NewWindowPolicy()
	shipmentClock := timeadapter.NewFixedClock(time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC))
	returnClock := timeadapter.NewFixedClock(time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC))
	idempotency := memory.NewIdempotencyStore()

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
	requestReturn := NewRequestReturnUseCase(orderRepo, returnRepo, returnClock)
	acceptReturn := NewAcceptReturnUseCase(returnRepo, returnPolicy, idempotency)
	rejectReturn := NewRejectReturnUseCase(returnRepo)
	completeRefund := NewCompleteRefundUseCase(returnRepo, refund.NewAcceptAllGateway(), inventory, idempotency)

	quote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(quote.ID, "CHAIR-001", 2)
	_, _ = submitQuote.Execute(quote.ID)
	order, _ := convertQuote.Execute(quote.ID)
	_, _ = capturePayment.Execute(order.ID)
	_, _ = createShipment.Execute(order.ID)

	_, err := requestReturn.Execute(order.ID, "Damaged", "")
	if err != domain.ErrActorRequired {
		t.Fatalf("expected %v, got %v", domain.ErrActorRequired, err)
	}

	request, err := requestReturn.Execute(order.ID, "Damaged", "warehouse-clerk-1")
	if err != nil {
		t.Fatalf("expected return request to succeed, got %v", err)
	}

	_, err = acceptReturn.Execute(request.ID, "", "return-accept-004")
	if err != domain.ErrActorRequired {
		t.Fatalf("expected %v, got %v", domain.ErrActorRequired, err)
	}

	_, err = rejectReturn.Execute(request.ID, "", "Outside window")
	if err != domain.ErrActorRequired {
		t.Fatalf("expected %v, got %v", domain.ErrActorRequired, err)
	}

	accepted, err := acceptReturn.Execute(request.ID, "warehouse-clerk-1", "return-accept-005")
	if err != nil {
		t.Fatalf("expected return acceptance to succeed, got %v", err)
	}

	_, err = completeRefund.Execute(accepted.ID, "", "refund-003")
	if err != domain.ErrActorRequired {
		t.Fatalf("expected %v, got %v", domain.ErrActorRequired, err)
	}
}

func TestAcceptReturnIsIdempotent(t *testing.T) {
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
	acceptReturn := NewAcceptReturnUseCase(returnRepo, returnPolicy, idempotency)

	quote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(quote.ID, "CHAIR-001", 2)
	_, _ = submitQuote.Execute(quote.ID)
	order, _ := convertQuote.Execute(quote.ID)
	_, _ = capturePayment.Execute(order.ID)
	_, _ = createShipment.Execute(order.ID)
	request, _ := requestReturn.Execute(order.ID, "Damaged", "warehouse-clerk-1")

	first, err := acceptReturn.Execute(request.ID, "warehouse-clerk-1", "return-accept-010")
	if err != nil {
		t.Fatalf("expected first accept to succeed, got %v", err)
	}

	second, err := acceptReturn.Execute(request.ID, "warehouse-clerk-9", "return-accept-010")
	if err != nil {
		t.Fatalf("expected second accept with same key to succeed, got %v", err)
	}

	if second.Status != domain.ReturnStatusAccepted || second.ReviewedBy != first.ReviewedBy {
		t.Fatalf("expected idempotent accept result to match first result, got status=%s reviewedBy=%s", second.Status, second.ReviewedBy)
	}
}

func TestCompleteRefundIsIdempotent(t *testing.T) {
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
	returnPolicy := returnpolicy.NewWindowPolicy()
	refundGateway := refund.NewAcceptAllGateway()
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
	acceptReturn := NewAcceptReturnUseCase(returnRepo, returnPolicy, idempotency)
	completeRefund := NewCompleteRefundUseCase(returnRepo, refundGateway, inventory, idempotency)

	quote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(quote.ID, "CHAIR-001", 2)
	_, _ = submitQuote.Execute(quote.ID)
	order, _ := convertQuote.Execute(quote.ID)
	_, _ = capturePayment.Execute(order.ID)
	_, _ = createShipment.Execute(order.ID)
	request, _ := requestReturn.Execute(order.ID, "Damaged", "warehouse-clerk-1")
	_, _ = acceptReturn.Execute(request.ID, "warehouse-clerk-1", "return-accept-011")

	first, err := completeRefund.Execute(request.ID, "manager-1", "refund-010")
	if err != nil {
		t.Fatalf("expected first refund to succeed, got %v", err)
	}

	second, err := completeRefund.Execute(request.ID, "manager-9", "refund-010")
	if err != nil {
		t.Fatalf("expected second refund with same key to succeed, got %v", err)
	}

	if second.Status != domain.ReturnStatusRefunded || second.ProcessedBy != first.ProcessedBy {
		t.Fatalf("expected idempotent refund result to match first result, got status=%s processedBy=%s", second.Status, second.ProcessedBy)
	}

	if inventory.Available("CHAIR-001") != 3 {
		t.Fatalf("expected stock to be restocked once to 3, got %d", inventory.Available("CHAIR-001"))
	}
}
