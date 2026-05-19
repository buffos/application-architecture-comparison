package pricing

import (
	"sort"

	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type PluginAwarePolicy struct {
	base    ports.PricingPolicy
	plugins ports.PluginRepository
}

func NewPluginAwarePolicy(base ports.PricingPolicy, plugins ports.PluginRepository) PluginAwarePolicy {
	return PluginAwarePolicy{
		base:    base,
		plugins: plugins,
	}
}

func (p PluginAwarePolicy) Price(product domain.Product, quantity int) (int, error) {
	price, err := p.base.Price(product, quantity)
	if err != nil {
		return 0, err
	}

	plugins, err := p.plugins.List()
	if err != nil {
		return 0, err
	}

	sort.Slice(plugins, func(i, j int) bool {
		return plugins[i].Key < plugins[j].Key
	})

	for _, plugin := range plugins {
		if plugin.Type != domain.PluginTypePricing || plugin.Status != domain.PluginStatusEnabled {
			continue
		}

		price -= (price * plugin.DiscountPercent) / 100
	}

	return price, nil
}
