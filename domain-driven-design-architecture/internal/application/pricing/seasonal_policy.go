package pricing

import "domain-driven-design-architecture/internal/domain/quoting"

type SeasonalPolicy struct {
	base            quoting.PricingPolicy
	discountPercent int64
}

func NewSeasonalPolicy(base quoting.PricingPolicy, discountPercent int64) SeasonalPolicy {
	if base == nil {
		base = quoting.TierPricingPolicy{}
	}
	return SeasonalPolicy{base: base, discountPercent: discountPercent}
}

func (p SeasonalPolicy) Price(product quoting.ProductPricingInput, tier quoting.CustomerPricingTier) (quoting.Money, error) {
	basePrice, err := p.base.Price(product, tier)
	if err != nil {
		return quoting.Money{}, err
	}
	adjustedCents := basePrice.Cents() * (100 - p.discountPercent) / 100
	return quoting.NewMoney(adjustedCents, basePrice.Currency())
}

var _ quoting.PricingPolicy = SeasonalPolicy{}
