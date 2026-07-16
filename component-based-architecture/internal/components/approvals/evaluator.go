package approvals

// Evaluator is the public policy contract provided by this component.
type Evaluator interface {
	RequiresApproval(submission QuoteSubmission) bool
}

type QuoteSubmission struct {
	Lines []QuoteSubmissionLine
}

type QuoteSubmissionLine struct {
	ProductCategory string
}
