package pricing

import "onion-architecture/internal/domain"

type pluginRepository interface {
	List() ([]domain.PluginRegistration, error)
}

type PluginPolicy struct {
	base    FixedPolicy
	plugins pluginRepository
}

func NewPluginPolicy(base FixedPolicy, plugins pluginRepository) PluginPolicy {
	return PluginPolicy{
		base:    base,
		plugins: plugins,
	}
}

func (p PluginPolicy) Adjust(product domain.Product) (domain.Product, error) {
	adjusted, err := p.base.Adjust(product)
	if err != nil {
		return domain.Product{}, err
	}

	plugins, err := p.plugins.List()
	if err != nil {
		return domain.Product{}, err
	}

	for _, plugin := range plugins {
		if plugin.Type != "pricing" || !plugin.Enabled {
			continue
		}

		switch plugin.Name {
		case "seasonal-pricing":
			adjusted.UnitPrice = adjusted.UnitPrice * 95 / 100
		}
	}

	return adjusted, nil
}
