package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/adapters/services/approval"
	"hexagonal-architecture/internal/adapters/services/payment"
	"hexagonal-architecture/internal/adapters/services/pricing"
	"hexagonal-architecture/internal/adapters/services/refund"
	"hexagonal-architecture/internal/adapters/services/returnpolicy"
	timeadapter "hexagonal-architecture/internal/adapters/services/time"
	"hexagonal-architecture/internal/core/application"
	"hexagonal-architecture/internal/core/domain"
)

func TestReturnHandlerGetsAndListsReturns(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	orderRepo := memory.NewOrderRepository()
	shipmentRepo := memory.NewShipmentRepository()
	returnRepo := memory.NewReturnRequestRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	inventory := memory.NewInventoryReservationAdapter(map[string]int{"CHAIR-001": 5})
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

	createQuote := application.NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := application.NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)
	submitQuote := application.NewSubmitQuoteUseCase(quoteRepo, approvalPolicy)
	convertQuote := application.NewConvertQuoteToOrderUseCase(quoteRepo, orderRepo, inventory)
	capturePayment := application.NewCapturePaymentUseCase(orderRepo, paymentGateway)
	createShipment := application.NewCreateShipmentUseCase(orderRepo, shipmentRepo, inventory, shipmentClock)
	requestReturn := application.NewRequestReturnUseCase(orderRepo, returnRepo, returnClock)
	acceptReturn := application.NewAcceptReturnUseCase(returnRepo, returnPolicy, idempotency)
	completeRefund := application.NewCompleteRefundUseCase(returnRepo, refundGateway, inventory, idempotency)
	getReturn := application.NewGetReturnRequestUseCase(returnRepo)
	listReturns := application.NewListReturnRequestsUseCase(returnRepo)
	handler := NewReturnHandler(getReturn, listReturns)

	quote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(quote.ID, "CHAIR-001", 2)
	_, _ = submitQuote.Execute(quote.ID)
	order, _ := convertQuote.Execute(quote.ID)
	_, _ = capturePayment.Execute(order.ID)
	_, _ = createShipment.Execute(order.ID)
	request, _ := requestReturn.Execute(order.ID, "Damaged", "warehouse-clerk-1")
	_, _ = acceptReturn.Execute(request.ID, "warehouse-clerk-1", "return-accept-201")
	_, _ = completeRefund.Execute(request.ID, "manager-1", "refund-201")

	getRequest := httptest.NewRequest(http.MethodGet, "/returns/"+request.ID, nil)
	getRecorder := httptest.NewRecorder()
	handler.ServeHTTP(getRecorder, getRequest)

	if getRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getRecorder.Code)
	}

	if !strings.Contains(getRecorder.Body.String(), `"status":"Refunded"`) {
		t.Fatalf("expected refunded return in body, got %s", getRecorder.Body.String())
	}

	listRequest := httptest.NewRequest(http.MethodGet, "/returns?status=Refunded", nil)
	listRecorder := httptest.NewRecorder()
	handler.ServeHTTP(listRecorder, listRequest)

	if listRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, listRecorder.Code)
	}

	body := listRecorder.Body.String()
	if !strings.Contains(body, `"id":"`+request.ID+`"`) || !strings.Contains(body, `"requestedBy":"warehouse-clerk-1"`) {
		t.Fatalf("expected return list response to contain request id and actor metadata, got %s", body)
	}
}
