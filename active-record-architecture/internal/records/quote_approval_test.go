package records

import "testing"

func TestQuoteEvaluateApprovalIsNonMutating(t *testing.T) {
	quote := &Quote{
		Status: QuoteStatusDraft,
		Lines:  []QuoteLine{{ProductCategory: "CustomBuild"}},
	}

	decision := quote.EvaluateApproval()
	if !decision.Required {
		t.Fatal("approval decision required = false, want true")
	}
	if len(decision.Reasons) != 1 || decision.Reasons[0] != ApprovalReasonCustomBuild {
		t.Fatalf("approval reasons = %#v, want [%q]", decision.Reasons, ApprovalReasonCustomBuild)
	}
	if quote.Status != QuoteStatusDraft {
		t.Fatalf("quote status = %q, want %q", quote.Status, QuoteStatusDraft)
	}
}

func TestQuoteEvaluateApprovalReturnsNoFindingsForStandardLine(t *testing.T) {
	quote := &Quote{Lines: []QuoteLine{{ProductCategory: "Standard"}}}

	decision := quote.EvaluateApproval()
	if decision.Required || len(decision.Reasons) != 0 {
		t.Fatalf("approval decision = %#v, want no findings", decision)
	}
}
