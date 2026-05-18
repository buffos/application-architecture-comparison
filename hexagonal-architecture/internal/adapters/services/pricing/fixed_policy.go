package pricing

import "hexagonal-architecture/internal/core/domain"

type FixedPricingPolicy struct{}

func NewFixedPricingPolicy() FixedPricingPolicy {
	return FixedPricingPolicy{}
}

func (FixedPricingPolicy) Price(product domain.Product, quantity int) (int, error) {
	return product.BasePrice, nil
}
