package returneligibility

import "microkernel-architecture/internal/kernel"

type Service struct{}

func NewService() Service {
	return Service{}
}

func (s Service) Allows(review kernel.ReturnEligibilityReview) bool {
	return review.Reason != "outside return window"
}
