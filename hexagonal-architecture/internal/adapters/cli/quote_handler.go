package cli

import (
	"fmt"

	"hexagonal-architecture/internal/core/application"
)

type QuoteHandler struct {
	createQuote application.CreateDraftQuoteUseCase
}

func NewQuoteHandler(createQuote application.CreateDraftQuoteUseCase) QuoteHandler {
	return QuoteHandler{createQuote: createQuote}
}

func (h QuoteHandler) RunDemo() (string, error) {
	quote, err := h.createQuote.Execute("customer-001")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("created draft quote: id=%s customer=%s status=%s", quote.ID, quote.CustomerID, quote.Status), nil
}
