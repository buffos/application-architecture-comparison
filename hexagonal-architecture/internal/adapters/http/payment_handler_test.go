package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/adapters/services/approval"
	"hexagonal-architecture/internal/adapters/services/payment"
	"hexagonal-architecture/internal/adapters/services/pricing"
	"hexagonal-architecture/internal/core/application"
	"hexagonal-architecture/internal/core/domain"
)

func TestPaymentHandlerCapturesAndApprovesReview(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	orderRepo := memory.NewOrderRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	inventory := memory.NewInventoryReservationAdapter(map[string]int{"DESK-001": 5})
	pricingPolicy := pricing.NewFixedPricingPolicy()
	approvalPolicy := approval.NewCategoryApprovalPolicy()

	_ = customerRepo.Save(domain.Customer{ID: "customer-001", Active: true})
	_ = productRepo.Save(domain.Product{SKU: "DESK-001", Name: "Executive Desk", Category: "CustomBuild", BasePrice: 50000, Available: true, ReturnWindowDays: 30})

	createQuote := application.NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := application.NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)
	submitQuote := application.NewSubmitQuoteUseCase(quoteRepo, approvalPolicy)
	approveQuote := application.NewApproveQuoteUseCase(quoteRepo)
	convertQuote := application.NewConvertQuoteToOrderUseCase(quoteRepo, orderRepo, inventory)

	quote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(quote.ID, "DESK-001", 1)
	_, _ = submitQuote.Execute(quote.ID)
	_, _ = approveQuote.Execute(quote.ID)
	order, _ := convertQuote.Execute(quote.ID)

	handler := NewPaymentHandler(
		application.NewCapturePaymentUseCase(orderRepo, payment.NewManualReviewGateway()),
		application.NewApprovePaymentReviewUseCase(orderRepo),
	)

	captureRequest := httptest.NewRequest(http.MethodPost, "/orders/"+order.ID+"/capture-payment", strings.NewReader(`{}`))
	captureRecorder := httptest.NewRecorder()
	handler.ServeHTTP(captureRecorder, captureRequest)

	if captureRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, captureRecorder.Code)
	}

	if !strings.Contains(captureRecorder.Body.String(), `"status":"PaymentReview"`) || !strings.Contains(captureRecorder.Body.String(), `"paymentStatus":"ManualReview"`) {
		t.Fatalf("expected payment review response, got %s", captureRecorder.Body.String())
	}

	approveRequest := httptest.NewRequest(http.MethodPost, "/orders/"+order.ID+"/approve-payment-review", strings.NewReader(`{"reviewedBy":"manager-1"}`))
	approveRecorder := httptest.NewRecorder()
	handler.ServeHTTP(approveRecorder, approveRequest)

	if approveRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, approveRecorder.Code)
	}

	if !strings.Contains(approveRecorder.Body.String(), `"status":"ReadyForFulfillment"`) || !strings.Contains(approveRecorder.Body.String(), `"paymentStatus":"Accepted"`) {
		t.Fatalf("expected ready-for-fulfillment response, got %s", approveRecorder.Body.String())
	}
}
