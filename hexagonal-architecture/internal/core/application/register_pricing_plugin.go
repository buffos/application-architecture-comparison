package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type RegisterPricingPluginUseCase struct {
	plugins ports.PluginRepository
}

func NewRegisterPricingPluginUseCase(plugins ports.PluginRepository) RegisterPricingPluginUseCase {
	return RegisterPricingPluginUseCase{plugins: plugins}
}

func (uc RegisterPricingPluginUseCase) Execute(key string, discountPercent int) (domain.PluginRegistration, error) {
	plugin, err := domain.NewPricingPlugin(key, discountPercent)
	if err != nil {
		return domain.PluginRegistration{}, err
	}

	if err := uc.plugins.Save(plugin); err != nil {
		return domain.PluginRegistration{}, err
	}

	return plugin, nil
}
