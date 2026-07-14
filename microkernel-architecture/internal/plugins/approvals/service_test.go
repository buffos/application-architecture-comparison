package approvals

import (
	"testing"

	"microkernel-architecture/internal/kernel"
)

func TestRequiresApprovalForCustomBuild(t *testing.T) {
	service := NewService()

	if !service.RequiresApproval(kernel.QuoteSubmission{
		Lines: []kernel.QuoteSubmissionLine{
			{ProductCategory: "CustomBuild"},
		},
	}) {
		t.Fatalf("expected custom build submission to require approval")
	}
}

func TestDoesNotRequireApprovalForStandard(t *testing.T) {
	service := NewService()

	if service.RequiresApproval(kernel.QuoteSubmission{
		Lines: []kernel.QuoteSubmissionLine{
			{ProductCategory: "Standard"},
		},
	}) {
		t.Fatalf("expected standard submission not to require approval")
	}
}
