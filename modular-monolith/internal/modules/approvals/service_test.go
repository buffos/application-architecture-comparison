package approvals

import "testing"

func TestRequiresApprovalForCustomBuild(t *testing.T) {
	service := NewService()

	requiresApproval := service.RequiresApproval(QuoteSubmission{
		Lines: []QuoteSubmissionLine{
			{ProductCategory: "CustomBuild"},
		},
	})

	if !requiresApproval {
		t.Fatalf("expected custom build quote to require approval")
	}
}

func TestDoesNotRequireApprovalForStandardQuote(t *testing.T) {
	service := NewService()

	requiresApproval := service.RequiresApproval(QuoteSubmission{
		Lines: []QuoteSubmissionLine{
			{ProductCategory: "Standard"},
		},
	})

	if requiresApproval {
		t.Fatalf("expected standard quote not to require approval")
	}
}
