package application

import (
	"testing"

	"layered-architecture/internal/domain"
	"layered-architecture/internal/infrastructure/memory"
)

func TestReportingServiceShowsPendingApprovalAndLowStock(t *testing.T) {
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	quoteRepo := memory.NewQuoteRepository()
	stockRepo := memory.NewStockRecordRepository()
	orderRepo := memory.NewOrderRepository()

	customerService := NewCustomerService(customerRepo)
	catalogService := NewCatalogService(productRepo)
	inventoryService := NewInventoryService(productRepo, stockRepo)
	quoteService := NewQuoteService(quoteRepo, customerRepo, productRepo, NoopPricingPluginRegistry{})
	orderService := NewOrderService(orderRepo, quoteRepo, stockRepo)
	reportingService := NewReportingQueryService(quoteRepo, orderRepo, stockRepo)

	customer, _ := customerService.CreateCustomer("Acme Corp", "Preferred", "Invoice30")
	standardProduct, _ := catalogService.CreateProduct("CHAIR-001", "Office Chair", "Standard", 10000, true)
	customProduct, _ := catalogService.CreateProduct("DESK-001", "Executive Desk", "CustomBuild", 50000, true)
	_, _ = inventoryService.ReceiveStock(standardProduct.SKU, 2)
	_, _ = inventoryService.ReceiveStock(customProduct.SKU, 5)

	approvedQuote, _ := quoteService.CreateDraftQuote(customer.ID)
	_, _ = quoteService.AddQuoteLine(approvedQuote.ID, standardProduct.SKU, 2)
	_, _ = quoteService.SubmitQuote(approvedQuote.ID)
	_, _ = orderService.ConvertQuoteToOrder(approvedQuote.ID)

	pendingQuote, _ := quoteService.CreateDraftQuote(customer.ID)
	_, _ = quoteService.AddQuoteLine(pendingQuote.ID, customProduct.SKU, 1)
	_, _ = quoteService.SubmitQuote(pendingQuote.ID)

	lowStockItems, err := reportingService.GetLowStockItems(0)
	if err != nil {
		t.Fatalf("expected low stock query to succeed, got %v", err)
	}

	if len(lowStockItems) != 1 || lowStockItems[0].SKU != standardProduct.SKU {
		t.Fatalf("expected one low stock item for %s, got %+v", standardProduct.SKU, lowStockItems)
	}

	awaitingApproval, err := reportingService.GetOrdersAwaitingApproval()
	if err != nil {
		t.Fatalf("expected awaiting approval query to succeed, got %v", err)
	}

	if len(awaitingApproval) != 1 || awaitingApproval[0].CurrentStatus != domain.QuoteStatusPendingApproval {
		t.Fatalf("expected one pending approval quote, got %+v", awaitingApproval)
	}

	conversionReport, err := reportingService.GetQuoteConversionReport()
	if err != nil {
		t.Fatalf("expected conversion report query to succeed, got %v", err)
	}

	if conversionReport.TotalQuotes != 2 || conversionReport.ConvertedQuotes != 1 {
		t.Fatalf("expected total=2 converted=1, got %+v", conversionReport)
	}
}
