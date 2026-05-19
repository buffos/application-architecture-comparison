package application

import (
	"testing"
	"time"

	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/adapters/services/approval"
	"hexagonal-architecture/internal/adapters/services/payment"
	"hexagonal-architecture/internal/adapters/services/pricing"
	"hexagonal-architecture/internal/adapters/services/returnpolicy"
	timeadapter "hexagonal-architecture/internal/adapters/services/time"
	"hexagonal-architecture/internal/core/domain"
)

func TestReturnRateByCategoryReport(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	orderRepo := memory.NewOrderRepository()
	shipmentRepo := memory.NewShipmentRepository()
	returnRepo := memory.NewReturnRequestRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	inventory := memory.NewInventoryReservationAdapter(map[string]int{
		"CHAIR-001": 10,
		"DESK-001":  10,
	})
	pricingPolicy := pricing.NewFixedPricingPolicy()
	approvalPolicy := approval.NewCategoryApprovalPolicy()
	paymentGateway := payment.NewAcceptAllGateway()
	returnPolicy := returnpolicy.NewWindowPolicy()
	shipmentClock := timeadapter.NewFixedClock(time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC))
	returnClock := timeadapter.NewFixedClock(time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC))
	idempotency := memory.NewIdempotencyStore()

	_ = customerRepo.Save(domain.Customer{ID: "customer-001", Active: true})
	_ = productRepo.Save(domain.Product{SKU: "CHAIR-001", Name: "Office Chair", Category: "Standard", BasePrice: 10000, Available: true, ReturnWindowDays: 30})
	_ = productRepo.Save(domain.Product{SKU: "DESK-001", Name: "Executive Desk", Category: "CustomBuild", BasePrice: 50000, Available: true, ReturnWindowDays: 30})

	createQuote := NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)
	submitQuote := NewSubmitQuoteUseCase(quoteRepo, approvalPolicy)
	approveQuote := NewApproveQuoteUseCase(quoteRepo)
	convertQuote := NewConvertQuoteToOrderUseCase(quoteRepo, orderRepo, inventory)
	capturePayment := NewCapturePaymentUseCase(orderRepo, paymentGateway)
	createShipment := NewCreateShipmentUseCase(orderRepo, shipmentRepo, inventory, shipmentClock)
	requestReturn := NewRequestReturnUseCase(orderRepo, returnRepo, returnClock)
	acceptReturn := NewAcceptReturnUseCase(orderRepo, returnRepo, returnPolicy, idempotency)
	reportUseCase := NewGetReturnRateByCategoryReportUseCase(orderRepo, returnRepo)

	standardOrder := createShippedOrder(t, createQuote, addQuoteLine, submitQuote, approveQuote, convertQuote, capturePayment, createShipment, "customer-001", "CHAIR-001", 2)
	customOrder := createShippedOrder(t, createQuote, addQuoteLine, submitQuote, approveQuote, convertQuote, capturePayment, createShipment, "customer-001", "DESK-001", 1)

	standardReturn, err := requestReturn.Execute(standardOrder.ID, "Damaged", "warehouse-clerk-1")
	if err != nil {
		t.Fatalf("expected standard return request to succeed, got %v", err)
	}
	_, err = acceptReturn.Execute(standardReturn.ID, "warehouse-clerk-1", "return-accept-101")
	if err != nil {
		t.Fatalf("expected standard return acceptance to succeed, got %v", err)
	}

	customReturn, err := requestReturn.Execute(customOrder.ID, "Wrong spec", "warehouse-clerk-1")
	if err != nil {
		t.Fatalf("expected custom return request to succeed, got %v", err)
	}

	report, err := reportUseCase.Execute()
	if err != nil {
		t.Fatalf("expected report to succeed, got %v", err)
	}

	if len(report) != 2 {
		t.Fatalf("expected 2 report rows, got %d", len(report))
	}

	if report[0].Category != "CustomBuild" || report[0].ShippedQuantity != 1 || report[0].ReturnQuantity != 0 || report[0].ReturnRate != 0 {
		t.Fatalf("unexpected custom build row: %+v", report[0])
	}

	if report[1].Category != "Standard" || report[1].ShippedQuantity != 2 || report[1].ReturnQuantity != 2 || report[1].ReturnRate != 1 {
		t.Fatalf("unexpected standard row: %+v", report[1])
	}

	storedCustomReturn, err := returnRepo.FindByID(customReturn.ID)
	if err != nil {
		t.Fatalf("expected custom return lookup to succeed, got %v", err)
	}

	if storedCustomReturn.Status != domain.ReturnStatusRequested {
		t.Fatalf("expected requested return to be excluded from numerator, got status %s", storedCustomReturn.Status)
	}
}

func createShippedOrder(
	t *testing.T,
	createQuote CreateDraftQuoteUseCase,
	addQuoteLine AddQuoteLineUseCase,
	submitQuote SubmitQuoteUseCase,
	approveQuote ApproveQuoteUseCase,
	convertQuote ConvertQuoteToOrderUseCase,
	capturePayment CapturePaymentUseCase,
	createShipment CreateShipmentUseCase,
	customerID string,
	sku string,
	quantity int,
) domain.Order {
	t.Helper()

	quote, err := createQuote.Execute(customerID)
	if err != nil {
		t.Fatalf("expected quote creation to succeed, got %v", err)
	}

	_, err = addQuoteLine.Execute(quote.ID, sku, quantity)
	if err != nil {
		t.Fatalf("expected quote line add to succeed, got %v", err)
	}

	quote, err = submitQuote.Execute(quote.ID)
	if err != nil {
		t.Fatalf("expected quote submission to succeed, got %v", err)
	}

	if quote.Status == domain.QuoteStatusPendingApproval {
		quote, err = approveQuote.Execute(quote.ID)
		if err != nil {
			t.Fatalf("expected quote approval to succeed, got %v", err)
		}
	}

	order, err := convertQuote.Execute(quote.ID)
	if err != nil {
		t.Fatalf("expected quote conversion to succeed, got %v", err)
	}

	_, err = capturePayment.Execute(order.ID)
	if err != nil {
		t.Fatalf("expected payment capture to succeed, got %v", err)
	}

	_, err = createShipment.Execute(order.ID)
	if err != nil {
		t.Fatalf("expected shipment creation to succeed, got %v", err)
	}

	return order
}
