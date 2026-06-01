package payment

import "onion-architecture/internal/domain"

type ManualReviewGateway struct{}

func NewManualReviewGateway() ManualReviewGateway {
	return ManualReviewGateway{}
}

func (g ManualReviewGateway) Capture(order domain.Order) (string, error) {
	return "Review", nil
}
