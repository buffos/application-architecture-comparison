package plugins

import "layered-architecture/internal/application"

type SeasonalDiscountPlugin struct {
	key             string
	discountPercent int
}

func NewSeasonalDiscountPlugin(key string, discountPercent int) SeasonalDiscountPlugin {
	return SeasonalDiscountPlugin{key: key, discountPercent: discountPercent}
}

func (p SeasonalDiscountPlugin) Key() string {
	return p.key
}

func (p SeasonalDiscountPlugin) Adjust(input application.PricingPluginInput) (application.PricingAdjustment, bool) {
	if input.BasePrice <= 0 || p.discountPercent <= 0 {
		return application.PricingAdjustment{}, false
	}

	adjusted := input.BasePrice - ((input.BasePrice * p.discountPercent) / 100)
	return application.PricingAdjustment{
		Label:         p.key,
		AdjustedPrice: adjusted,
	}, true
}
