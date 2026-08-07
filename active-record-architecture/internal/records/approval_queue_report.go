package records

import "sort"

// ApprovalQueueItem is one quote waiting for review.
type ApprovalQueueItem struct {
	QuoteID    string
	CustomerID string
	Reasons    []string
}

// GetOrdersAwaitingApproval projects pending quote Active Records into a
// deterministic approval queue without mutating the quotes.
func GetOrdersAwaitingApproval(db *Database) ([]ApprovalQueueItem, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}

	quotes, err := ListQuotes(db, QuoteStatusPendingApproval, "")
	if err != nil {
		return nil, err
	}
	items := make([]ApprovalQueueItem, 0, len(quotes))
	for _, quote := range quotes {
		decision := quote.EvaluateApproval()
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
