package pricing

import "clean-architecture/internal/entities"

type PluginRepository interface {
	List() ([]entities.PluginRegistration, error)
}

type PluginPolicy struct {
	plugins PluginRepository
}

func NewPluginPolicy(plugins PluginRepository) PluginPolicy {
	return PluginPolicy{plugins: plugins}
}

func (p PluginPolicy) AdjustUnitPrice(product entities.Product, quantity int) (int, error) {
	_ = quantity
	price := product.BasePrice

	plugins, err := p.plugins.List()
	if err != nil {
		return 0, err
	}

	for _, plugin := range plugins {
		if !plugin.Enabled {
			continue
		}

		switch plugin.Name {
		case "seasonal-pricing":
			price = price * 95 / 100
		}
	}

	return price, nil
}
