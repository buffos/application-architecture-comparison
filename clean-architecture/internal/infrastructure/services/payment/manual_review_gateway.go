package payment

import (
	"clean-architecture/internal/entities"
	"clean-architecture/internal/usecases"
)

type ManualReviewGateway struct {
	threshold int
}

func NewManualReviewGateway(threshold int) ManualReviewGateway {
	return ManualReviewGateway{threshold: threshold}
}

func (g ManualReviewGateway) Capture(order entities.Order) (string, error) {
	total := 0
	for _, line := range order.Lines {
		total += line.LineTotal
	}

	if total >= g.threshold {
		return usecases.PaymentCaptureReview, nil
	}

	return usecases.PaymentCaptureApproved, nil
}
