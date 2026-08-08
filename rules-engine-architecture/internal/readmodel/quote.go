package readmodel

import (
	"sort"

	"rules-engine-architecture/internal/engine"
)

// EvaluatedQuote is the read-side input produced by an application evaluation
// workflow. The projector does not invoke the Engine.
type EvaluatedQuote struct {
	Memory   *engine.WorkingMemory
	Decision engine.PolicyDecision
}

type QuoteSummary struct {
	ID              string
	CustomerID      string
	Status          string
	SubtotalCents   int64
	DiscountPercent int
	Outcome         engine.DecisionOutcome
	RequiredReviews []engine.ReviewRequirement
}

func ProjectQuoteList(evaluated []EvaluatedQuote) []QuoteSummary {
	summaries := make([]QuoteSummary, 0, len(evaluated))
	for _, quote := range evaluated {
		summaries = append(summaries, QuoteSummary{
			ID:              quote.Memory.Quote.ID,
			CustomerID:      quote.Memory.Quote.CustomerID,
			Status:          quote.Memory.Quote.Status,
			SubtotalCents:   quoteSubtotal(quote.Memory),
			DiscountPercent: quote.Memory.Quote.DiscountPercent,
			Outcome:         quote.Decision.Outcome,
			RequiredReviews: append([]engine.ReviewRequirement(nil), quote.Decision.RequiredReviews...),
		})
	}

	sort.Slice(summaries, func(left, right int) bool {
		return summaries[left].ID < summaries[right].ID
	})
	return summaries
}

func quoteSubtotal(memory *engine.WorkingMemory) int64 {
	var subtotal int64
	for _, line := range memory.Quote.Lines {
		subtotal += int64(line.Quantity) * line.UnitPriceCents
	}
	return subtotal
}
