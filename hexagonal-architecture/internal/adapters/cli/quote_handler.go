package cli

import (
	"fmt"

	"hexagonal-architecture/internal/core/application"
)

type QuoteHandler struct {
	createQuote application.CreateDraftQuoteUseCase
	getQuote    application.GetQuoteUseCase
}

func NewQuoteHandler(createQuote application.CreateDraftQuoteUseCase, getQuote application.GetQuoteUseCase) QuoteHandler {
	return QuoteHandler{
		createQuote: createQuote,
		getQuote:    getQuote,
	}
}

func (h QuoteHandler) RunDemo() (string, error) {
	quote, err := h.createQuote.Execute("customer-001")
	if err != nil {
		return "", err
	}

	loadedQuote, err := h.getQuote.Execute(quote.ID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("created draft quote: id=%s customer=%s status=%s\nloaded draft quote: id=%s customer=%s status=%s", quote.ID, quote.CustomerID, quote.Status, loadedQuote.ID, loadedQuote.CustomerID, loadedQuote.Status), nil
}
