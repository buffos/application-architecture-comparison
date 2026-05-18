package cli

import (
	"fmt"

	"hexagonal-architecture/internal/core/application"
)

type QuoteHandler struct {
	createQuote    application.CreateDraftQuoteUseCase
	addQuoteLine   application.AddQuoteLineUseCase
	submitQuote    application.SubmitQuoteUseCase
	convertQuote   application.ConvertQuoteToOrderUseCase
	capturePayment application.CapturePaymentUseCase
	createShipment application.CreateShipmentUseCase
	getQuote       application.GetQuoteUseCase
}

func NewQuoteHandler(createQuote application.CreateDraftQuoteUseCase, addQuoteLine application.AddQuoteLineUseCase, submitQuote application.SubmitQuoteUseCase, convertQuote application.ConvertQuoteToOrderUseCase, capturePayment application.CapturePaymentUseCase, createShipment application.CreateShipmentUseCase, getQuote application.GetQuoteUseCase) QuoteHandler {
	return QuoteHandler{
		createQuote:    createQuote,
		addQuoteLine:   addQuoteLine,
		submitQuote:    submitQuote,
		convertQuote:   convertQuote,
		capturePayment: capturePayment,
		createShipment: createShipment,
		getQuote:       getQuote,
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

	order, err := h.convertQuote.Execute(quote.ID)
	if err != nil {
		return "", err
	}

	paidOrder, err := h.capturePayment.Execute(order.ID)
	if err != nil {
		return "", err
	}

	shipment, err := h.createShipment.Execute(order.ID)
	if err != nil {
		return "", err
	}

	loadedQuote, err := h.getQuote.Execute(quote.ID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("created draft quote: id=%s customer=%s status=%s\nadded quote line: id=%s lines=%d status=%s\nsubmitted quote: id=%s lines=%d status=%s\nconverted order: id=%s sourceQuote=%s status=%s\ncaptured payment: id=%s status=%s payment=%s\ncreated shipment: id=%s order=%s status=%s lines=%d\nloaded draft quote: id=%s customer=%s lines=%d status=%s", quote.ID, quote.CustomerID, quote.Status, quoteWithLine.ID, len(quoteWithLine.Lines), quoteWithLine.Status, submittedQuote.ID, len(submittedQuote.Lines), submittedQuote.Status, order.ID, order.SourceQuoteID, order.Status, paidOrder.ID, paidOrder.Status, paidOrder.PaymentStatus, shipment.ID, shipment.OrderID, shipment.Status, len(shipment.Lines), loadedQuote.ID, loadedQuote.CustomerID, len(loadedQuote.Lines), loadedQuote.Status), nil
}
