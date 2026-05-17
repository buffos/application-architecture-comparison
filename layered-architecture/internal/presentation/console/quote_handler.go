package console

import (
	"fmt"
	"strings"

	"layered-architecture/internal/application"
)

type QuoteHandler struct {
	customerService application.CustomerService
	catalogService  application.CatalogService
	quoteService    application.QuoteService
}

func NewQuoteHandler(customerService application.CustomerService, catalogService application.CatalogService, quoteService application.QuoteService) QuoteHandler {
	return QuoteHandler{
		customerService: customerService,
		catalogService:  catalogService,
		quoteService:    quoteService,
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

	lines := []string{
		fmt.Sprintf("created customer: id=%s tier=%s paymentTerms=%s", customer.ID, customer.Tier, customer.PaymentTerms),
		fmt.Sprintf("created product: sku=%s name=%s category=%s", product.SKU, product.Name, product.Category),
		fmt.Sprintf("created draft quote: id=%s customer=%s status=%s", createdQuote.ID, createdQuote.CustomerID, createdQuote.Status),
		fmt.Sprintf("added quote line: id=%s sku=%s name=%s lines=%d status=%s", quoteWithLine.ID, quoteWithLine.Lines[0].SKU, quoteWithLine.Lines[0].ProductNameSnapshot, len(quoteWithLine.Lines), quoteWithLine.Status),
		fmt.Sprintf("submitted quote: id=%s lines=%d status=%s", submittedQuote.ID, len(submittedQuote.Lines), submittedQuote.Status),
		fmt.Sprintf("loaded quote: id=%s customer=%s status=%s", loadedQuote.ID, loadedQuote.CustomerID, loadedQuote.Status),
	}

	return strings.Join(lines, "\n"), nil
}
