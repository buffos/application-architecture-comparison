package cli

import (
	"fmt"

	"hexagonal-architecture/internal/core/application"
)

type QuoteHandler struct {
	createQuote  application.CreateDraftQuoteUseCase
	addQuoteLine application.AddQuoteLineUseCase
	submitQuote  application.SubmitQuoteUseCase
	getQuote     application.GetQuoteUseCase
}

func NewQuoteHandler(createQuote application.CreateDraftQuoteUseCase, addQuoteLine application.AddQuoteLineUseCase, submitQuote application.SubmitQuoteUseCase, getQuote application.GetQuoteUseCase) QuoteHandler {
	return QuoteHandler{
		createQuote:  createQuote,
		addQuoteLine: addQuoteLine,
		submitQuote:  submitQuote,
		getQuote:     getQuote,
	}
}

func (h QuoteHandler) RunDemo() (string, error) {
	quote, err := h.createQuote.Execute("customer-001")
	if err != nil {
		return "", err
	}

	quoteWithLine, err := h.addQuoteLine.Execute(quote.ID, "CHAIR-001", 2)
	if err != nil {
		return "", err
	}

	submittedQuote, err := h.submitQuote.Execute(quote.ID)
	if err != nil {
		return "", err
	}

	loadedQuote, err := h.getQuote.Execute(quote.ID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("created draft quote: id=%s customer=%s status=%s\nadded quote line: id=%s lines=%d status=%s\nsubmitted quote: id=%s lines=%d status=%s\nloaded draft quote: id=%s customer=%s lines=%d status=%s", quote.ID, quote.CustomerID, quote.Status, quoteWithLine.ID, len(quoteWithLine.Lines), quoteWithLine.Status, submittedQuote.ID, len(submittedQuote.Lines), submittedQuote.Status, loadedQuote.ID, loadedQuote.CustomerID, len(loadedQuote.Lines), loadedQuote.Status), nil
}
