package payment

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type ManualReviewGateway struct{}

func NewManualReviewGateway() ManualReviewGateway {
	return ManualReviewGateway{}
}

func (ManualReviewGateway) Capture(order domain.Order) (ports.PaymentResult, error) {
	return ports.PaymentResultManualReview, nil
}
