package scripts

import (
	"reflect"
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestEvaluateQuoteApprovalRequiresReviewForCustomBuild(t *testing.T) {
	quote := data.Quote{
		Status: data.QuoteStatusDraft,
		Lines:  []data.QuoteLine{{ProductCategory: "CustomBuild"}},
	}

	got := EvaluateQuoteApproval(quote)

	if !got.Required {
		t.Fatal("Required = false, want true")
	}
	if !reflect.DeepEqual(got.Reasons, []string{ApprovalReasonCustomBuild}) {
		t.Fatalf("Reasons = %#v, want %#v", got.Reasons, []string{ApprovalReasonCustomBuild})
	}
	if quote.Status != data.QuoteStatusDraft {
		t.Fatalf("quote status = %q, want %q", quote.Status, data.QuoteStatusDraft)
	}
}

func TestEvaluateQuoteApprovalDoesNotRequireReviewForStandardQuote(t *testing.T) {
	got := EvaluateQuoteApproval(data.Quote{
		Lines: []data.QuoteLine{{ProductCategory: "Standard"}},
	})

	if got.Required {
		t.Fatal("Required = true, want false")
	}
	if len(got.Reasons) != 0 {
		t.Fatalf("Reasons = %#v, want empty", got.Reasons)
	}
}
