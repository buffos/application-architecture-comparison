package approvals

// Component owns approval policy for this lesson.
type Component struct{}

func NewComponent() *Component {
	return &Component{}
}

func (c *Component) RequiresApproval(submission QuoteSubmission) bool {
	for _, line := range submission.Lines {
		if line.ProductCategory == "CustomBuild" {
			return true
		}
	}
	return false
}

var _ Evaluator = (*Component)(nil)
