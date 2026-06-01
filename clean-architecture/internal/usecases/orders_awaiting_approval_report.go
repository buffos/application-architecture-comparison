package usecases

import "clean-architecture/internal/entities"

type OrdersAwaitingApprovalReportInput struct{}

type ApprovalQueueItem struct {
	QuoteID     string
	CustomerID  string
	LineCount   int
	TotalAmount int
}

type OrdersAwaitingApprovalReportOutput struct {
	Count int
	Items []ApprovalQueueItem
}

type OrdersAwaitingApprovalReportInputBoundary interface {
	Execute(input OrdersAwaitingApprovalReportInput) error
}

type OrdersAwaitingApprovalReportOutputBoundary interface {
	Present(output OrdersAwaitingApprovalReportOutput) error
}

type ApprovalQueueQuoteReader interface {
	ListByStatus(status string) ([]entities.Quote, error)
}

type OrdersAwaitingApprovalReportInteractor struct {
	quotes ApprovalQueueQuoteReader
	output OrdersAwaitingApprovalReportOutputBoundary
}

func NewOrdersAwaitingApprovalReportInteractor(quotes ApprovalQueueQuoteReader, output OrdersAwaitingApprovalReportOutputBoundary) OrdersAwaitingApprovalReportInteractor {
	return OrdersAwaitingApprovalReportInteractor{
		quotes: quotes,
		output: output,
	}
}

func (uc OrdersAwaitingApprovalReportInteractor) Execute(input OrdersAwaitingApprovalReportInput) error {
	_ = input

	quotes, err := uc.quotes.ListByStatus(entities.QuoteStatusPendingApproval)
	if err != nil {
		return err
	}

	items := make([]ApprovalQueueItem, 0, len(quotes))
	for _, quote := range quotes {
		totalAmount := 0
		for _, line := range quote.Lines {
			totalAmount += line.LineTotal
		}

		items = append(items, ApprovalQueueItem{
			QuoteID:     quote.ID,
			CustomerID:  quote.CustomerID,
			LineCount:   len(quote.Lines),
			TotalAmount: totalAmount,
		})
	}

	return uc.output.Present(OrdersAwaitingApprovalReportOutput{
		Count: len(items),
		Items: items,
	})
}
