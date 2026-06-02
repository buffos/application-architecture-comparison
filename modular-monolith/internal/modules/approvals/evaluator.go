package approvals

type Evaluator interface {
	RequiresApproval(submission QuoteSubmission) bool
}

type QuoteSubmission struct {
	Lines []QuoteSubmissionLine
}

type QuoteSubmissionLine struct {
	ProductCategory string
}
