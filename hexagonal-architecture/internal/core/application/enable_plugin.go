package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type EnablePluginUseCase struct {
	plugins ports.PluginRepository
}

func NewEnablePluginUseCase(plugins ports.PluginRepository) EnablePluginUseCase {
	return EnablePluginUseCase{plugins: plugins}
}

func (uc EnablePluginUseCase) Execute(key string) (domain.PluginRegistration, error) {
	plugin, err := uc.plugins.FindByKey(key)
	if err != nil {
		return domain.PluginRegistration{}, err
	}

	if err := plugin.Enable(); err != nil {
		return domain.PluginRegistration{}, err
	}

	if err := uc.plugins.Save(plugin); err != nil {
		return domain.PluginRegistration{}, err
	}

	return plugin, nil
}
