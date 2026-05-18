package console

import (
	"fmt"
	"strings"

	"layered-architecture/internal/application"
	"layered-architecture/internal/domain"
)

type QuoteHandler struct {
	customerService    application.CustomerService
	catalogService     application.CatalogService
	inventoryService   application.InventoryService
	quoteService       application.QuoteService
	orderService       application.OrderService
	paymentService     application.PaymentService
	fulfillmentService application.FulfillmentService
	returnService      application.ReturnService
	reportingService   application.ReportingQueryService
}

func NewQuoteHandler(customerService application.CustomerService, catalogService application.CatalogService, inventoryService application.InventoryService, quoteService application.QuoteService, orderService application.OrderService, paymentService application.PaymentService, fulfillmentService application.FulfillmentService, returnService application.ReturnService, reportingService application.ReportingQueryService) QuoteHandler {
	return QuoteHandler{
		customerService:    customerService,
		catalogService:     catalogService,
		inventoryService:   inventoryService,
		quoteService:       quoteService,
		orderService:       orderService,
		paymentService:     paymentService,
		fulfillmentService: fulfillmentService,
		returnService:      returnService,
		reportingService:   reportingService,
	}
}

func (h QuoteHandler) RunDemo() (string, error) {
	customer, err := h.customerService.CreateCustomer("Acme Corp", "Preferred", "Invoice30")
	if err != nil {
		return "", err
	}

	product, err := h.catalogService.CreateProduct("CHAIR-001", "Office Chair", "Standard", true)
	if err != nil {
		return "", err
	}

	stock, err := h.inventoryService.ReceiveStock(product.SKU, 10)
	if err != nil {
		return "", err
	}

	createdQuote, err := h.quoteService.CreateDraftQuote(customer.ID)
	if err != nil {
		return "", err
	}

	quoteWithLine, err := h.quoteService.AddQuoteLine(createdQuote.ID, product.SKU, 2)
	if err != nil {
		return "", err
	}

	submittedQuote, err := h.quoteService.SubmitQuote(createdQuote.ID)
	if err != nil {
		return "", err
	}

	quoteReadyForConversion := submittedQuote
	if submittedQuote.Status == domain.QuoteStatusPendingApproval {
		quoteReadyForConversion, err = h.quoteService.ApproveQuote(createdQuote.ID)
		if err != nil {
			return "", err
		}
	}

	loadedQuote, err := h.quoteService.GetQuote(createdQuote.ID)
	if err != nil {
		return "", err
	}

	order, err := h.orderService.ConvertQuoteToOrder(createdQuote.ID)
	if err != nil {
		return "", err
	}

	paidOrder, err := h.paymentService.CapturePayment(order.ID)
	if err != nil {
		return "", err
	}

	shipment, err := h.fulfillmentService.CreateShipment(order.ID)
	if err != nil {
		return "", err
	}

	returnRequest, err := h.returnService.RequestReturn(order.ID, "Damaged")
	if err != nil {
		return "", err
	}

	acceptedReturn, err := h.returnService.AcceptReturn(returnRequest.ID)
	if err != nil {
		return "", err
	}

	loadedOrder, err := h.orderService.GetOrder(order.ID)
	if err != nil {
		return "", err
	}

	cancelledQuote, err := h.quoteService.CreateDraftQuote(customer.ID)
	if err != nil {
		return "", err
	}

	cancelledQuoteWithLine, err := h.quoteService.AddQuoteLine(cancelledQuote.ID, product.SKU, 1)
	if err != nil {
		return "", err
	}

	cancelledQuoteReady, err := h.quoteService.SubmitQuote(cancelledQuote.ID)
	if err != nil {
		return "", err
	}

	cancelledOrder, err := h.orderService.ConvertQuoteToOrder(cancelledQuote.ID)
	if err != nil {
		return "", err
	}

	cancelledOrder, err = h.orderService.CancelOrder(cancelledOrder.ID)
	if err != nil {
		return "", err
	}

	customBuildProduct, err := h.catalogService.CreateProduct("DESK-001", "Executive Desk", "CustomBuild", true)
	if err != nil {
		return "", err
	}

	if _, err := h.inventoryService.ReceiveStock(customBuildProduct.SKU, 4); err != nil {
		return "", err
	}

	pendingQuote, err := h.quoteService.CreateDraftQuote(customer.ID)
	if err != nil {
		return "", err
	}

	pendingQuoteWithLine, err := h.quoteService.AddQuoteLine(pendingQuote.ID, customBuildProduct.SKU, 1)
	if err != nil {
		return "", err
	}

	pendingApprovalQuote, err := h.quoteService.SubmitQuote(pendingQuote.ID)
	if err != nil {
		return "", err
	}

	lowStockItems, err := h.reportingService.GetLowStockItems(3)
	if err != nil {
		return "", err
	}

	awaitingApproval, err := h.reportingService.GetOrdersAwaitingApproval()
	if err != nil {
		return "", err
	}

	conversionReport, err := h.reportingService.GetQuoteConversionReport()
	if err != nil {
		return "", err
	}

	lines := []string{
		fmt.Sprintf("created customer: id=%s tier=%s paymentTerms=%s", customer.ID, customer.Tier, customer.PaymentTerms),
		fmt.Sprintf("created product: sku=%s name=%s category=%s", product.SKU, product.Name, product.Category),
		fmt.Sprintf("received stock: sku=%s onHand=%d reserved=%d available=%d", stock.SKU, stock.OnHand, stock.Reserved, stock.Available()),
		fmt.Sprintf("created draft quote: id=%s customer=%s status=%s", createdQuote.ID, createdQuote.CustomerID, createdQuote.Status),
		fmt.Sprintf("added quote line: id=%s sku=%s name=%s lines=%d status=%s", quoteWithLine.ID, quoteWithLine.Lines[0].SKU, quoteWithLine.Lines[0].ProductNameSnapshot, len(quoteWithLine.Lines), quoteWithLine.Status),
		fmt.Sprintf("submitted quote: id=%s lines=%d status=%s", submittedQuote.ID, len(submittedQuote.Lines), submittedQuote.Status),
		fmt.Sprintf("quote ready for conversion: id=%s status=%s", quoteReadyForConversion.ID, quoteReadyForConversion.Status),
		fmt.Sprintf("loaded quote: id=%s customer=%s status=%s", loadedQuote.ID, loadedQuote.CustomerID, loadedQuote.Status),
		fmt.Sprintf("converted order: id=%s sourceQuote=%s status=%s payment=%s", order.ID, order.SourceQuoteID, order.Status, order.PaymentStatus),
		fmt.Sprintf("captured payment: id=%s status=%s payment=%s", paidOrder.ID, paidOrder.Status, paidOrder.PaymentStatus),
		fmt.Sprintf("created shipment: id=%s order=%s status=%s lines=%d", shipment.ID, shipment.OrderID, shipment.Status, len(shipment.Lines)),
		fmt.Sprintf("requested return: id=%s order=%s status=%s lines=%d", returnRequest.ID, returnRequest.OrderID, returnRequest.Status, len(returnRequest.Lines)),
		fmt.Sprintf("accepted return: id=%s status=%s", acceptedReturn.ID, acceptedReturn.Status),
		fmt.Sprintf("loaded order: id=%s customer=%s lines=%d status=%s payment=%s", loadedOrder.ID, loadedOrder.CustomerID, len(loadedOrder.Lines), loadedOrder.Status, loadedOrder.PaymentStatus),
		fmt.Sprintf("created cancellation quote: id=%s status=%s", cancelledQuote.ID, cancelledQuote.Status),
		fmt.Sprintf("added cancellation quote line: id=%s lines=%d status=%s", cancelledQuoteWithLine.ID, len(cancelledQuoteWithLine.Lines), cancelledQuoteWithLine.Status),
		fmt.Sprintf("submitted cancellation quote: id=%s status=%s", cancelledQuoteReady.ID, cancelledQuoteReady.Status),
		fmt.Sprintf("cancelled order: id=%s status=%s payment=%s", cancelledOrder.ID, cancelledOrder.Status, cancelledOrder.PaymentStatus),
		fmt.Sprintf("created custom build product: sku=%s category=%s", customBuildProduct.SKU, customBuildProduct.Category),
		fmt.Sprintf("created pending approval quote: id=%s status=%s lines=%d", pendingApprovalQuote.ID, pendingApprovalQuote.Status, len(pendingQuoteWithLine.Lines)),
		fmt.Sprintf("low stock report: items=%d", len(lowStockItems)),
		fmt.Sprintf("awaiting approval report: quotes=%d", len(awaitingApproval)),
		fmt.Sprintf("quote conversion report: total=%d converted=%d rate=%.2f", conversionReport.TotalQuotes, conversionReport.ConvertedQuotes, conversionReport.ConversionRate),
	}

	return strings.Join(lines, "\n"), nil
}
