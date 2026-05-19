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
	"hexagonal-architecture/internal/adapters/services/returnpolicy"
	timeadapter "hexagonal-architecture/internal/adapters/services/time"
	"hexagonal-architecture/internal/core/application"
	"hexagonal-architecture/internal/core/domain"
)

func TestReportHandlerReturnsQuoteConversion(t *testing.T) {
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
	_ = productRepo.Save(domain.Product{SKU: "CHAIR-001", Name: "Office Chair", Category: "Standard", BasePrice: 10000, Available: true, ReturnWindowDays: 30})

	createQuote := application.NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := application.NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)
	submitQuote := application.NewSubmitQuoteUseCase(quoteRepo, approvalPolicy)
	approvalQueueReport := application.NewGetOrdersAwaitingApprovalReportUseCase(quoteRepo)
	convertQuote := application.NewConvertQuoteToOrderUseCase(quoteRepo, orderRepo, inventory)
	reportUseCase := application.NewGetQuoteConversionReportUseCase(quoteRepo, orderRepo)
	returnRateReport := application.NewGetReturnRateByCategoryReportUseCase(orderRepo, returnRepo)
	topDiscountedReport := application.NewGetTopDiscountedProductsReportUseCase(quoteRepo)
	lowStockReport := application.NewGetLowStockItemsReportUseCase(inventory)
	handler := NewReportHandler(approvalQueueReport, reportUseCase, returnRateReport, topDiscountedReport, lowStockReport)

	quote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(quote.ID, "CHAIR-001", 1)
	_, _ = submitQuote.Execute(quote.ID)
	_, _ = convertQuote.Execute(quote.ID)

	request := httptest.NewRequest(http.MethodGet, "/reports/quote-conversion", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	body := recorder.Body.String()
	if !strings.Contains(body, `"totalQuotes":1`) || !strings.Contains(body, `"convertedQuotes":1`) {
		t.Fatalf("expected quote conversion report in body, got %s", body)
	}
}

func TestReportHandlerReturnsReturnRateByCategory(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	orderRepo := memory.NewOrderRepository()
	shipmentRepo := memory.NewShipmentRepository()
	returnRepo := memory.NewReturnRequestRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	inventory := memory.NewInventoryReservationAdapter(map[string]int{
		"CHAIR-001": 5,
		"DESK-001":  5,
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

	createQuote := application.NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := application.NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)
	submitQuote := application.NewSubmitQuoteUseCase(quoteRepo, approvalPolicy)
	approveQuote := application.NewApproveQuoteUseCase(quoteRepo)
	convertQuote := application.NewConvertQuoteToOrderUseCase(quoteRepo, orderRepo, inventory)
	capturePayment := application.NewCapturePaymentUseCase(orderRepo, paymentGateway)
	createShipment := application.NewCreateShipmentUseCase(orderRepo, shipmentRepo, inventory, shipmentClock)
	requestReturn := application.NewRequestReturnUseCase(orderRepo, returnRepo, returnClock)
	acceptReturn := application.NewAcceptReturnUseCase(orderRepo, returnRepo, returnPolicy, idempotency)
	approvalQueueReport := application.NewGetOrdersAwaitingApprovalReportUseCase(quoteRepo)
	quoteConversion := application.NewGetQuoteConversionReportUseCase(quoteRepo, orderRepo)
	returnRateReport := application.NewGetReturnRateByCategoryReportUseCase(orderRepo, returnRepo)
	topDiscountedReport := application.NewGetTopDiscountedProductsReportUseCase(quoteRepo)
	lowStockReport := application.NewGetLowStockItemsReportUseCase(inventory)
	handler := NewReportHandler(approvalQueueReport, quoteConversion, returnRateReport, topDiscountedReport, lowStockReport)

	standardQuote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(standardQuote.ID, "CHAIR-001", 2)
	_, _ = submitQuote.Execute(standardQuote.ID)
	standardOrder, _ := convertQuote.Execute(standardQuote.ID)
	_, _ = capturePayment.Execute(standardOrder.ID)
	_, _ = createShipment.Execute(standardOrder.ID)
	standardReturn, _ := requestReturn.Execute(standardOrder.ID, "Damaged", "warehouse-clerk-1")
	_, _ = acceptReturn.Execute(standardReturn.ID, "warehouse-clerk-1", "return-accept-201")

	customQuote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(customQuote.ID, "DESK-001", 1)
	_, _ = submitQuote.Execute(customQuote.ID)
	_, _ = approveQuote.Execute(customQuote.ID)
	customOrder, _ := convertQuote.Execute(customQuote.ID)
	_, _ = capturePayment.Execute(customOrder.ID)
	_, _ = createShipment.Execute(customOrder.ID)

	request := httptest.NewRequest(http.MethodGet, "/reports/return-rate-by-category", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	body := recorder.Body.String()
	if !strings.Contains(body, `"category":"CustomBuild"`) || !strings.Contains(body, `"category":"Standard"`) {
		t.Fatalf("expected category rows in body, got %s", body)
	}
	if !strings.Contains(body, `"returnQuantity":2`) || !strings.Contains(body, `"returnRate":1`) {
		t.Fatalf("expected returned standard quantity and rate in body, got %s", body)
	}
}

func TestReportHandlerReturnsTopDiscountedProducts(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	orderRepo := memory.NewOrderRepository()
	returnRepo := memory.NewReturnRequestRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	pricingPolicy := httpDiscountedPricingPolicy{
		adjustedBySKU: map[string]int{
			"CHAIR-001": 9000,
			"DESK-001":  45000,
		},
	}

	_ = customerRepo.Save(domain.Customer{ID: "customer-001", Active: true})
	_ = productRepo.Save(domain.Product{SKU: "CHAIR-001", Name: "Office Chair", Category: "Standard", BasePrice: 10000, Available: true, ReturnWindowDays: 30})
	_ = productRepo.Save(domain.Product{SKU: "DESK-001", Name: "Executive Desk", Category: "CustomBuild", BasePrice: 50000, Available: true, ReturnWindowDays: 30})

	createQuote := application.NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := application.NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)
	approvalQueueReport := application.NewGetOrdersAwaitingApprovalReportUseCase(quoteRepo)
	quoteConversion := application.NewGetQuoteConversionReportUseCase(quoteRepo, orderRepo)
	returnRateReport := application.NewGetReturnRateByCategoryReportUseCase(orderRepo, returnRepo)
	topDiscountedReport := application.NewGetTopDiscountedProductsReportUseCase(quoteRepo)
	lowStockReport := application.NewGetLowStockItemsReportUseCase(memory.NewInventoryReservationAdapter(map[string]int{}))
	handler := NewReportHandler(approvalQueueReport, quoteConversion, returnRateReport, topDiscountedReport, lowStockReport)

	quote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(quote.ID, "CHAIR-001", 2)
	_, _ = addQuoteLine.Execute(quote.ID, "DESK-001", 1)

	request := httptest.NewRequest(http.MethodGet, "/reports/top-discounted-products", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	body := recorder.Body.String()
	if !strings.Contains(body, `"sku":"DESK-001"`) || !strings.Contains(body, `"sku":"CHAIR-001"`) {
		t.Fatalf("expected discounted products in body, got %s", body)
	}
	if !strings.Contains(body, `"totalDiscountAmount":5000`) {
		t.Fatalf("expected top discount amount in body, got %s", body)
	}
}

func TestReportHandlerReturnsLowStockItems(t *testing.T) {
	inventory := memory.NewInventoryReservationAdapter(map[string]int{
		"CHAIR-001": 2,
		"DESK-001":  6,
		"LAMP-001":  1,
	})
	inventory.SetReorderThreshold("CHAIR-001", 3)
	inventory.SetReorderThreshold("DESK-001", 5)
	inventory.SetReorderThreshold("LAMP-001", 1)

	quoteConversion := application.NewGetQuoteConversionReportUseCase(memory.NewQuoteRepository(), memory.NewOrderRepository())
	returnRateReport := application.NewGetReturnRateByCategoryReportUseCase(memory.NewOrderRepository(), memory.NewReturnRequestRepository())
	topDiscountedReport := application.NewGetTopDiscountedProductsReportUseCase(memory.NewQuoteRepository())
	lowStockReport := application.NewGetLowStockItemsReportUseCase(inventory)
	handler := NewReportHandler(application.NewGetOrdersAwaitingApprovalReportUseCase(memory.NewQuoteRepository()), quoteConversion, returnRateReport, topDiscountedReport, lowStockReport)

	request := httptest.NewRequest(http.MethodGet, "/reports/low-stock-items", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	body := recorder.Body.String()
	if !strings.Contains(body, `"sku":"LAMP-001"`) || !strings.Contains(body, `"sku":"CHAIR-001"`) {
		t.Fatalf("expected low-stock items in body, got %s", body)
	}
	if !strings.Contains(body, `"reorderThreshold":3`) {
		t.Fatalf("expected reorder threshold in body, got %s", body)
	}
}

func TestReportHandlerReturnsOrdersAwaitingApproval(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	pricingPolicy := pricing.NewFixedPricingPolicy()
	approvalPolicy := approval.NewCategoryApprovalPolicy()

	_ = customerRepo.Save(domain.Customer{ID: "customer-001", Active: true})
	_ = productRepo.Save(domain.Product{SKU: "DESK-001", Name: "Executive Desk", Category: "CustomBuild", BasePrice: 50000, Available: true, ReturnWindowDays: 30})

	createQuote := application.NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := application.NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)
	submitQuote := application.NewSubmitQuoteUseCase(quoteRepo, approvalPolicy)
	approvalQueueReport := application.NewGetOrdersAwaitingApprovalReportUseCase(quoteRepo)
	handler := NewReportHandler(
		approvalQueueReport,
		application.NewGetQuoteConversionReportUseCase(quoteRepo, memory.NewOrderRepository()),
		application.NewGetReturnRateByCategoryReportUseCase(memory.NewOrderRepository(), memory.NewReturnRequestRepository()),
		application.NewGetTopDiscountedProductsReportUseCase(quoteRepo),
		application.NewGetLowStockItemsReportUseCase(memory.NewInventoryReservationAdapter(map[string]int{})),
	)

	quote, _ := createQuote.Execute("customer-001")
	_, _ = addQuoteLine.Execute(quote.ID, "DESK-001", 1)
	_, _ = submitQuote.Execute(quote.ID)

	request := httptest.NewRequest(http.MethodGet, "/reports/orders-awaiting-approval", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	body := recorder.Body.String()
	if !strings.Contains(body, `"quoteId":"`+quote.ID+`"`) || !strings.Contains(body, `"customerId":"customer-001"`) {
		t.Fatalf("expected awaiting-approval queue item in body, got %s", body)
	}
	if !strings.Contains(body, `"totalAmount":50000`) {
		t.Fatalf("expected total amount in body, got %s", body)
	}
}

type httpDiscountedPricingPolicy struct {
	adjustedBySKU map[string]int
}

func (p httpDiscountedPricingPolicy) Price(product domain.Product, quantity int) (int, error) {
	if adjusted, ok := p.adjustedBySKU[product.SKU]; ok {
		return adjusted, nil
	}

	return product.BasePrice, nil
}
