package console

import (
	"fmt"
	"strings"

	"layered-architecture/internal/application"
)

type QuoteHandler struct {
	service application.QuoteService
}

func NewQuoteHandler(service application.QuoteService) QuoteHandler {
	return QuoteHandler{service: service}
}

func (h QuoteHandler) RunDemo() (string, error) {
	createdQuote, err := h.service.CreateDraftQuote("customer-001")
	if err != nil {
		return "", err
	}

	quoteWithLine, err := h.service.AddQuoteLine(createdQuote.ID, "Office Chair", 2)
	if err != nil {
		return "", err
	}

	submittedQuote, err := h.service.SubmitQuote(createdQuote.ID)
	if err != nil {
		return "", err
	}

	loadedQuote, err := h.service.GetQuote(createdQuote.ID)
	if err != nil {
		return "", err
	}

	lines := []string{
		fmt.Sprintf("created draft quote: id=%s customer=%s status=%s", createdQuote.ID, createdQuote.CustomerID, createdQuote.Status),
		fmt.Sprintf("added quote line: id=%s lines=%d status=%s", quoteWithLine.ID, len(quoteWithLine.Lines), quoteWithLine.Status),
		fmt.Sprintf("submitted quote: id=%s lines=%d status=%s", submittedQuote.ID, len(submittedQuote.Lines), submittedQuote.Status),
		fmt.Sprintf("loaded draft quote: id=%s customer=%s status=%s", loadedQuote.ID, loadedQuote.CustomerID, loadedQuote.Status),
	}

	return strings.Join(lines, "\n"), nil
}
