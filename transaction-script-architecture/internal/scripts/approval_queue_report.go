package scripts

import (
	"sort"

	"transaction-script-architecture/internal/data"
)

type ApprovalQueueItem struct {
	QuoteID    string
	CustomerID string
	Reasons    []string
}

func GetOrdersAwaitingApproval(store *data.Store) ([]ApprovalQueueItem, error) {
	if store == nil {
		return nil, ErrStoreRequired
	}

	items := make([]ApprovalQueueItem, 0)
	for _, quote := range store.Quotes {
		if quote.Status != data.QuoteStatusPendingApproval {
			continue
		}

		decision := EvaluateQuoteApproval(quote)
		items = append(items, ApprovalQueueItem{
			QuoteID:    quote.ID,
			CustomerID: quote.CustomerID,
			Reasons:    append([]string(nil), decision.Reasons...),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].QuoteID < items[j].QuoteID
	})

	return items, nil
}
