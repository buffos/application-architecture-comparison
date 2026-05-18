package application

import (
	"testing"

	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/adapters/services/approval"
	"hexagonal-architecture/internal/adapters/services/payment"
	"hexagonal-architecture/internal/adapters/services/pricing"
	"hexagonal-architecture/internal/adapters/services/refund"
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

	_ = customerRepo.Save(domain.Customer{ID: "customer-001", Active: true})
	_ = productRepo.Save(domain.Product{
		SKU:       "CHAIR-001",
		Name:      "Office Chair",
		Category:  "Standard",
		BasePrice: 10000,
		Available: true,
	})

	createQuote := NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)
	submitQuote := NewSubmitQuoteUseCase(quoteRepo, approvalPolicy)
	convertQuote := NewConvertQuoteToOrderUseCase(quoteRepo, orderRepo, inventory)
	capturePayment := NewCapturePaymentUseCase(orderRepo, paymentGateway)
	createShipment := NewCreateShipmentUseCase(orderRepo, shipmentRepo, inventory)
	requestReturn := NewRequestReturnUseCase(orderRepo, returnRepo)
	completeRefund := NewCompleteRefundUseCase(returnRepo, refundGateway)

	quote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(quote.ID, "CHAIR-001", 2)
	_, _ = submitQuote.Execute(quote.ID)
	order, _ := convertQuote.Execute(quote.ID)
	_, _ = capturePayment.Execute(order.ID)
	_, _ = createShipment.Execute(order.ID)

	request, err := requestReturn.Execute(order.ID, "Damaged")
	if err != nil {
		t.Fatalf("expected return request to succeed, got %v", err)
	}

	if request.Status != domain.ReturnStatusRequested {
		t.Fatalf("expected return status %s, got %s", domain.ReturnStatusRequested, request.Status)
	}

	refunded, err := completeRefund.Execute(request.ID)
	if err != nil {
		t.Fatalf("expected refund to succeed, got %v", err)
	}

	if refunded.Status != domain.ReturnStatusRefunded {
		t.Fatalf("expected return status %s, got %s", domain.ReturnStatusRefunded, refunded.Status)
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

	_ = customerRepo.Save(domain.Customer{ID: "customer-001", Active: true})
	_ = productRepo.Save(domain.Product{
		SKU:       "CHAIR-001",
		Name:      "Office Chair",
		Category:  "Standard",
		BasePrice: 10000,
		Available: true,
	})

	createQuote := NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)
	submitQuote := NewSubmitQuoteUseCase(quoteRepo, approvalPolicy)
	convertQuote := NewConvertQuoteToOrderUseCase(quoteRepo, orderRepo, inventory)
	requestReturn := NewRequestReturnUseCase(orderRepo, returnRepo)

	quote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(quote.ID, "CHAIR-001", 2)
	_, _ = submitQuote.Execute(quote.ID)
	order, _ := convertQuote.Execute(quote.ID)

	_, err := requestReturn.Execute(order.ID, "Changed mind")
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

	_ = customerRepo.Save(domain.Customer{ID: "customer-001", Active: true})
	_ = productRepo.Save(domain.Product{
		SKU:       "LAMP-001",
		Name:      "Clearance Lamp",
		Category:  "Clearance",
		BasePrice: 4000,
		Available: true,
	})

	createQuote := NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)
	submitQuote := NewSubmitQuoteUseCase(quoteRepo, approvalPolicy)
	convertQuote := NewConvertQuoteToOrderUseCase(quoteRepo, orderRepo, inventory)
	capturePayment := NewCapturePaymentUseCase(orderRepo, paymentGateway)
	createShipment := NewCreateShipmentUseCase(orderRepo, shipmentRepo, inventory)
	requestReturn := NewRequestReturnUseCase(orderRepo, returnRepo)

	quote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(quote.ID, "LAMP-001", 1)
	_, _ = submitQuote.Execute(quote.ID)
	order, _ := convertQuote.Execute(quote.ID)
	_, _ = capturePayment.Execute(order.ID)
	_, _ = createShipment.Execute(order.ID)

	_, err := requestReturn.Execute(order.ID, "No longer needed")
	if err != domain.ErrReturnNotEligible {
		t.Fatalf("expected %v, got %v", domain.ErrReturnNotEligible, err)
	}
}
