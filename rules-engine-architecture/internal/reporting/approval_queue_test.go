package reporting

import (
	"testing"

	"rules-engine-architecture/internal/engine"
	"rules-engine-architecture/internal/readmodel"
)

func TestBuildOrdersAwaitingApprovalReportFiltersManagerReviews(t *testing.T) {
	managerReview := evaluatedQuote("Q-2002", "CUST-002", []engine.ReviewRequirement{
		engine.ReviewPayment,
		engine.ReviewManagerApproval,
	})
	paymentOnly := evaluatedQuote("Q-1001", "CUST-001", []engine.ReviewRequirement{
		engine.ReviewPayment,
	})

	rows := BuildOrdersAwaitingApprovalReport([]readmodel.EvaluatedQuote{paymentOnly, managerReview})

	if len(rows) != 1 || rows[0].QuoteID != "Q-2002" {
		t.Fatalf("expected only manager review row, got %+v", rows)
	}
	if len(rows[0].Reasons) != 1 || rows[0].Reasons[0] != "manager approval needed" {
		t.Fatalf("expected copied reason, got %+v", rows[0].Reasons)
	}
}

func evaluatedQuote(
	quoteID string,
	customerID string,
	reviews []engine.ReviewRequirement,
) readmodel.EvaluatedQuote {
	memory := engine.NewWorkingMemory(
		engine.CustomerFact{ID: customerID},
		engine.QuoteFact{ID: quoteID, CustomerID: customerID},
		nil,
	)
	return readmodel.EvaluatedQuote{
		Memory: memory,
		Decision: engine.PolicyDecision{
			RequiredReviews: reviews,
			Reasons: []engine.Finding{{
				Message: "manager approval needed",
			}},
		},
	}
}
