package ports

import "hexagonal-architecture/internal/core/domain"

type PaymentGateway interface {
	Capture(order domain.Order) (bool, error)
}
