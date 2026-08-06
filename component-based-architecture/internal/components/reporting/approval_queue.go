package reporting

import "component-based-architecture/internal/components/quotes"

type ApprovalQueueRow struct {
	QuoteID     string
	CustomerID  string
	LineCount   int
	TotalAmount int
}

type OrdersAwaitingApprovalReport struct {
	Rows []ApprovalQueueRow
}

type ApprovalQueueComponent struct {
	quotes quotes.QuoteLookup
}

func NewApprovalQueueComponent(quotes quotes.QuoteLookup) *ApprovalQueueComponent {
	return &ApprovalQueueComponent{quotes: quotes}
}

func (c *ApprovalQueueComponent) OrdersAwaitingApprovalReport() OrdersAwaitingApprovalReport {
	pending := c.quotes.ListQuotes(quotes.ListQuotesQuery{Status: quotes.QuoteStatusPendingApproval})
	report := OrdersAwaitingApprovalReport{Rows: make([]ApprovalQueueRow, 0, len(pending))}
	for _, quote := range pending {
		report.Rows = append(report.Rows, ApprovalQueueRow{QuoteID: quote.QuoteID, CustomerID: quote.CustomerID, LineCount: quote.LineCount, TotalAmount: quote.TotalAmount})
	}
	return report
}
