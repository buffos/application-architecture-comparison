package quoting

type ApprovalReasonCode string

const ApprovalReasonCustomBuild ApprovalReasonCode = "CustomBuildRequiresApproval"

type ApprovalReason struct {
	Code    ApprovalReasonCode
	Message string
}

type ApprovalDecision struct {
	Required bool
	Reasons  []ApprovalReason
}

// QuoteApprovalService evaluates approval policy without changing the Quote
// aggregate. Workflow transitions belong to a later application lesson.
type QuoteApprovalService struct{}

func NewQuoteApprovalService() QuoteApprovalService { return QuoteApprovalService{} }

func (QuoteApprovalService) Evaluate(quote Quote) ApprovalDecision {
	reasons := make([]ApprovalReason, 0)
	for _, line := range quote.Lines() {
		if line.ProductCategory() == ProductCategoryCustomBuild {
			reasons = append(reasons, ApprovalReason{Code: ApprovalReasonCustomBuild, Message: "custom-build lines require approval"})
		}
	}
	return ApprovalDecision{Required: len(reasons) > 0, Reasons: reasons}
}
