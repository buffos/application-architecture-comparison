package application

import (
	"sort"

	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type OrdersAwaitingApprovalRow struct {
	QuoteID     string
	CustomerID  string
	LineCount   int
	TotalAmount int
}

type GetOrdersAwaitingApprovalReportUseCase struct {
	quotes ports.QuoteRepository
}

func NewGetOrdersAwaitingApprovalReportUseCase(quotes ports.QuoteRepository) GetOrdersAwaitingApprovalReportUseCase {
	return GetOrdersAwaitingApprovalReportUseCase{quotes: quotes}
}

func (uc GetOrdersAwaitingApprovalReportUseCase) Execute() ([]OrdersAwaitingApprovalRow, error) {
	quotes, err := uc.quotes.ListByStatus(domain.QuoteStatusPendingApproval)
	if err != nil {
		return nil, err
	}

	rows := make([]OrdersAwaitingApprovalRow, 0, len(quotes))
	for _, quote := range quotes {
		total := 0
		for _, line := range quote.Lines {
			total += line.LineTotal
		}

		rows = append(rows, OrdersAwaitingApprovalRow{
			QuoteID:     quote.ID,
			CustomerID:  quote.CustomerID,
			LineCount:   len(quote.Lines),
			TotalAmount: total,
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].QuoteID < rows[j].QuoteID
	})

	return rows, nil
}
