package reports

import "rich-domain-model-architecture/internal/domain/quoting"

type ApprovalQueueRow struct {
	QuoteID    string
	CustomerID string
	LineCount  int
	TotalCents int64
	Currency   string
}

type OrdersAwaitingApprovalReport struct {
	Rows []ApprovalQueueRow
}

func BuildOrdersAwaitingApprovalReport(quotes []quoting.Quote) OrdersAwaitingApprovalReport {
	report := OrdersAwaitingApprovalReport{Rows: make([]ApprovalQueueRow, 0)}
	for _, quote := range quotes {
		if quote.Status() != quoting.QuoteStatusPendingApproval {
			continue
		}
		total, err := quote.Total()
		if err != nil {
			continue
		}
		report.Rows = append(report.Rows, ApprovalQueueRow{QuoteID: string(quote.ID()), CustomerID: string(quote.CustomerID()), LineCount: len(quote.Lines()), TotalCents: total.Cents(), Currency: total.Currency()})
	}
	return report
}
