package pricing

import "clean-architecture/internal/entities"

type FixedPolicy struct{}

func NewFixedPolicy() FixedPolicy {
	return FixedPolicy{}
}

func (p FixedPolicy) AdjustUnitPrice(product entities.Product, quantity int) (int, error) {
	_ = quantity
	return product.BasePrice, nil
}
