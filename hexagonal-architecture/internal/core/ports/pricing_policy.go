package ports

import "hexagonal-architecture/internal/core/domain"

type PricingPolicy interface {
	Price(product domain.Product, quantity int) (int, error)
}
