package seasonalpricing

import "microkernel-architecture/internal/kernel"

type Service struct {
	base            kernel.QuotePricer
	discountPercent int
}

func NewService(base kernel.QuotePricer, discountPercent int) Service {
	return Service{
		base:            base,
		discountPercent: discountPercent,
	}
}

func (s Service) UnitPriceForQuote(input kernel.QuotePricingInput) (int, error) {
	price, err := s.base.UnitPriceForQuote(input)
	if err != nil {
		return 0, err
	}

	if input.ProductCategory != "CustomBuild" || s.discountPercent <= 0 {
		return price, nil
	}

	discounted := price - (price*s.discountPercent)/100
	return discounted, nil
}
