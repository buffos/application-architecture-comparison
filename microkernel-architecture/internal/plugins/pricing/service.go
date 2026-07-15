package pricing

import "microkernel-architecture/internal/kernel"

type Service struct{}

func NewService() Service {
	return Service{}
}

func (s Service) UnitPriceForQuote(input kernel.QuotePricingInput) (int, error) {
	return input.UnitPrice, nil
}
