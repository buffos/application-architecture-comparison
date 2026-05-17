package console

import (
	"fmt"
	"strings"

	"layered-architecture/internal/application"
)

type QuoteHandler struct {
	customerService  application.CustomerService
	catalogService   application.CatalogService
	inventoryService application.InventoryService
	quoteService     application.QuoteService
	orderService     application.OrderService
}

func NewQuoteHandler(customerService application.CustomerService, catalogService application.CatalogService, inventoryService application.InventoryService, quoteService application.QuoteService, orderService application.OrderService) QuoteHandler {
	return QuoteHandler{
		customerService:  customerService,
		catalogService:   catalogService,
		inventoryService: inventoryService,
		quoteService:     quoteService,
		orderService:     orderService,
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

	loadedQuote, err := h.quoteService.GetQuote(createdQuote.ID)
	if err != nil {
		return "", err
	}

	order, err := h.orderService.ConvertQuoteToOrder(createdQuote.ID)
	if err != nil {
		return "", err
	}

	loadedOrder, err := h.orderService.GetOrder(order.ID)
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
		fmt.Sprintf("loaded quote: id=%s customer=%s status=%s", loadedQuote.ID, loadedQuote.CustomerID, loadedQuote.Status),
		fmt.Sprintf("converted order: id=%s sourceQuote=%s status=%s", order.ID, order.SourceQuoteID, order.Status),
		fmt.Sprintf("loaded order: id=%s customer=%s lines=%d status=%s", loadedOrder.ID, loadedOrder.CustomerID, len(loadedOrder.Lines), loadedOrder.Status),
	}

	return strings.Join(lines, "\n"), nil
}
