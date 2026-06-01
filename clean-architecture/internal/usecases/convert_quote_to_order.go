package usecases

import "clean-architecture/internal/entities"

type ConvertQuoteToOrderInput struct {
	QuoteID string
}

type ConvertQuoteToOrderOutput struct {
	OrderID       string
	SourceQuoteID string
	Status        string
	Lines         int
}

type ConvertQuoteToOrderInputBoundary interface {
	Execute(input ConvertQuoteToOrderInput) error
}

type ConvertQuoteToOrderOutputBoundary interface {
	Present(output ConvertQuoteToOrderOutput) error
}

type OrderWriter interface {
	Save(order entities.Order) error
}

type ConvertQuoteToOrderInteractor struct {
	quotes QuoteReader
	orders OrderWriter
	output ConvertQuoteToOrderOutputBoundary
}

func NewConvertQuoteToOrderInteractor(quotes QuoteReader, orders OrderWriter, output ConvertQuoteToOrderOutputBoundary) ConvertQuoteToOrderInteractor {
	return ConvertQuoteToOrderInteractor{
		quotes: quotes,
		orders: orders,
		output: output,
	}
}

func (uc ConvertQuoteToOrderInteractor) Execute(input ConvertQuoteToOrderInput) error {
	quote, err := uc.quotes.FindByID(input.QuoteID)
	if err != nil {
		return err
	}

	order, err := entities.NewOrderFromApprovedQuote(quote)
	if err != nil {
		return err
	}

	if err := uc.orders.Save(order); err != nil {
		return err
	}

	return uc.output.Present(ConvertQuoteToOrderOutput{
		OrderID:       order.ID,
		SourceQuoteID: order.SourceQuoteID,
		Status:        order.Status,
		Lines:         len(order.Lines),
	})
}
