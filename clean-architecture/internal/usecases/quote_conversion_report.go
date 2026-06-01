package usecases

import "clean-architecture/internal/entities"

type QuoteConversionReportInput struct{}

type QuoteConversionReportOutput struct {
	TotalQuotes     int
	ApprovedQuotes  int
	ConvertedQuotes int
	ConversionRate  float64
}

type QuoteConversionReportInputBoundary interface {
	Execute(input QuoteConversionReportInput) error
}

type QuoteConversionReportOutputBoundary interface {
	Present(output QuoteConversionReportOutput) error
}

type QuoteReportReader interface {
	ListByStatus(status string) ([]entities.Quote, error)
}

type OrderReportReader interface {
	ListByStatus(status string) ([]entities.Order, error)
}

type QuoteConversionReportInteractor struct {
	quotes QuoteReportReader
	orders OrderReportReader
	output QuoteConversionReportOutputBoundary
}

func NewQuoteConversionReportInteractor(quotes QuoteReportReader, orders OrderReportReader, output QuoteConversionReportOutputBoundary) QuoteConversionReportInteractor {
	return QuoteConversionReportInteractor{
		quotes: quotes,
		orders: orders,
		output: output,
	}
}

func (uc QuoteConversionReportInteractor) Execute(input QuoteConversionReportInput) error {
	_ = input

	draftQuotes, err := uc.quotes.ListByStatus(entities.QuoteStatusDraft)
	if err != nil {
		return err
	}

	pendingApprovalQuotes, err := uc.quotes.ListByStatus(entities.QuoteStatusPendingApproval)
	if err != nil {
		return err
	}

	approvedQuotes, err := uc.quotes.ListByStatus(entities.QuoteStatusApproved)
	if err != nil {
		return err
	}

	pendingPaymentOrders, err := uc.orders.ListByStatus(entities.OrderStatusPendingPayment)
	if err != nil {
		return err
	}

	paidOrders, err := uc.orders.ListByStatus(entities.OrderStatusPaid)
	if err != nil {
		return err
	}

	shippedOrders, err := uc.orders.ListByStatus(entities.OrderStatusShipped)
	if err != nil {
		return err
	}

	cancelledOrders, err := uc.orders.ListByStatus(entities.OrderStatusCancelled)
	if err != nil {
		return err
	}

	totalQuotes := len(draftQuotes) + len(pendingApprovalQuotes) + len(approvedQuotes)
	convertedQuotes := len(pendingPaymentOrders) + len(paidOrders) + len(shippedOrders) + len(cancelledOrders)

	var conversionRate float64
	if totalQuotes > 0 {
		conversionRate = float64(convertedQuotes) / float64(totalQuotes)
	}

	return uc.output.Present(QuoteConversionReportOutput{
		TotalQuotes:     totalQuotes,
		ApprovedQuotes:  len(approvedQuotes),
		ConvertedQuotes: convertedQuotes,
		ConversionRate:  conversionRate,
	})
}
