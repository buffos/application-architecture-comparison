package pricing

import (
	"errors"

	"rich-domain-model-architecture/internal/domain/quoting"
)

var ErrDiscountPercentInvalid = errors.New("discount percent must be between 0 and 100")

// SeasonalPolicy is an application-supplied pricing plugin. It composes a
// domain pricing policy and changes only the resulting Money value.
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

func (policy SeasonalPolicy) Price(product quoting.ProductPricingInput, tier quoting.CustomerPricingTier) (quoting.Money, error) {
	if policy.discountPercent < 0 || policy.discountPercent > 100 {
		return quoting.Money{}, ErrDiscountPercentInvalid
	}
	basePrice, err := policy.base.Price(product, tier)
	if err != nil {
		return quoting.Money{}, err
	}
	adjustedCents := basePrice.Cents() * (100 - policy.discountPercent) / 100
	return quoting.NewMoney(adjustedCents, basePrice.Currency())
}

var _ quoting.PricingPolicy = SeasonalPolicy{}
