package pricing

import "onion-architecture/internal/domain"

type FixedPolicy struct{}

func NewFixedPolicy() FixedPolicy {
	return FixedPolicy{}
}

func (p FixedPolicy) Adjust(product domain.Product) (domain.Product, error) {
	return product, nil
}
