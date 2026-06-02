package approvals

type Service struct{}

func NewService() Service {
	return Service{}
}

func (s Service) RequiresApproval(submission QuoteSubmission) bool {
	for _, line := range submission.Lines {
		if line.ProductCategory == "CustomBuild" {
			return true
		}
	}

	return false
}
