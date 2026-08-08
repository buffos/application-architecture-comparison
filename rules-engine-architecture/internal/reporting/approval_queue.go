package reporting

import (
	"sort"

	"rules-engine-architecture/internal/engine"
	"rules-engine-architecture/internal/readmodel"
)

type ApprovalQueueRow struct {
	QuoteID         string
	CustomerID      string
	RequiredReviews []engine.ReviewRequirement
	Reasons         []string
}

func BuildOrdersAwaitingApprovalReport(
	evaluated []readmodel.EvaluatedQuote,
) []ApprovalQueueRow {
	rows := make([]ApprovalQueueRow, 0)
	for _, quote := range evaluated {
		if !containsReview(quote.Decision.RequiredReviews, engine.ReviewManagerApproval) {
			continue
		}

		reasons := make([]string, 0, len(quote.Decision.Reasons))
		for _, reason := range quote.Decision.Reasons {
			reasons = append(reasons, reason.Message)
		}
		rows = append(rows, ApprovalQueueRow{
			QuoteID:         quote.Memory.Quote.ID,
			CustomerID:      quote.Memory.Quote.CustomerID,
			RequiredReviews: append([]engine.ReviewRequirement(nil), quote.Decision.RequiredReviews...),
			Reasons:         reasons,
		})
	}

	sort.Slice(rows, func(left, right int) bool {
		return rows[left].QuoteID < rows[right].QuoteID
	})
	return rows
}

func containsReview(reviews []engine.ReviewRequirement, wanted engine.ReviewRequirement) bool {
	for _, review := range reviews {
		if review == wanted {
			return true
		}
	}
	return false
}
