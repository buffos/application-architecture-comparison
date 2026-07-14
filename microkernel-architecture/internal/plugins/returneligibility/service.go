package returneligibility

import "microkernel-architecture/internal/kernel"

type Service struct{}

func NewService() Service {
	return Service{}
}

func (s Service) Allows(review kernel.ReturnEligibilityReview) bool {
	if review.ShippedAt.IsZero() || review.RequestedAt.Before(review.ShippedAt) {
		return false
	}

	for _, line := range review.Lines {
		deadline := review.ShippedAt.AddDate(0, 0, line.ReturnWindowDays)
		if review.RequestedAt.After(deadline) {
			return false
		}
	}

	return true
}
