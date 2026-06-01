package application

type OrdersAwaitingApprovalRow struct {
	QuoteID     string
	CustomerID  string
	LineCount   int
	TotalAmount int
}

type OrdersAwaitingApprovalReportService struct {
	quotes QuoteFinder
}

func NewOrdersAwaitingApprovalReportService(quotes QuoteFinder) OrdersAwaitingApprovalReportService {
	return OrdersAwaitingApprovalReportService{
		quotes: quotes,
	}
}

func (s OrdersAwaitingApprovalReportService) Execute() ([]OrdersAwaitingApprovalRow, error) {
	quotes, err := s.quotes.ListByStatus("PendingApproval")
	if err != nil {
		return nil, err
	}

	result := make([]OrdersAwaitingApprovalRow, 0, len(quotes))
	for _, quote := range quotes {
		totalAmount := 0
		for _, line := range quote.Lines {
			totalAmount += line.UnitPrice * line.Quantity
		}

		result = append(result, OrdersAwaitingApprovalRow{
			QuoteID:     quote.ID,
			CustomerID:  quote.CustomerID,
			LineCount:   len(quote.Lines),
			TotalAmount: totalAmount,
		})
	}

	return result, nil
}
