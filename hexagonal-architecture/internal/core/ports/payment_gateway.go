package ports

import "hexagonal-architecture/internal/core/domain"

type PaymentResult string

const PaymentResultAccepted PaymentResult = "Accepted"
const PaymentResultManualReview PaymentResult = "ManualReview"
const PaymentResultFailed PaymentResult = "Failed"

type PaymentGateway interface {
	Capture(order domain.Order) (PaymentResult, error)
}
