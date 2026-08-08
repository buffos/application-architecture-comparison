package readmodel

import (
	"testing"

	"rules-engine-architecture/internal/engine"
)

func TestProjectQuoteListSortsAndSummarizesWithoutInference(t *testing.T) {
	first := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001"},
		engine.QuoteFact{
			ID:              "Q-2002",
			CustomerID:      "CUST-001",
			Status:          "Draft",
			DiscountPercent: 10,
			Lines:           []engine.QuoteLineFact{{Quantity: 2, UnitPriceCents: 12500}},
		},
		nil,
	)
	second := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-002"},
		engine.QuoteFact{
			ID:         "Q-1001",
			CustomerID: "CUST-002",
			Status:     "Submitted",
			Lines:      []engine.QuoteLineFact{{Quantity: 1, UnitPriceCents: 5000}},
		},
		nil,
	)
	secondDecision := engine.PolicyDecision{
		Outcome:         engine.OutcomeNeedsApproval,
		RequiredReviews: []engine.ReviewRequirement{engine.ReviewManagerApproval},
	}

	summaries := ProjectQuoteList([]EvaluatedQuote{
		{Memory: first, Decision: engine.PolicyDecision{Outcome: engine.OutcomeAllowed}},
		{Memory: second, Decision: secondDecision},
	})

	if len(summaries) != 2 {
		t.Fatalf("expected two quote summaries, got %d", len(summaries))
	}
	if summaries[0].ID != "Q-1001" || summaries[1].ID != "Q-2002" {
		t.Fatalf("expected id sorting, got %+v", summaries)
	}
	if summaries[1].SubtotalCents != 25000 {
		t.Fatalf("expected subtotal 25000, got %d", summaries[1].SubtotalCents)
	}
	if len(second.Trace) != 0 || len(first.Trace) != 0 {
		t.Fatal("expected projection not to execute the Rule Engine")
	}
}
