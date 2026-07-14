package approvals

import "microkernel-architecture/internal/kernel"

type Service struct{}

func NewService() Service {
	return Service{}
}

func (s Service) RequiresApproval(submission kernel.QuoteSubmission) bool {
	for _, line := range submission.Lines {
		if line.ProductCategory == "CustomBuild" {
			return true
		}
	}

	return false
}
