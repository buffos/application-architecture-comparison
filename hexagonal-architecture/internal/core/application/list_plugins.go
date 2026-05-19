package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type ListPluginsUseCase struct {
	plugins ports.PluginRepository
}

func NewListPluginsUseCase(plugins ports.PluginRepository) ListPluginsUseCase {
	return ListPluginsUseCase{plugins: plugins}
}

func (uc ListPluginsUseCase) Execute() ([]domain.PluginRegistration, error) {
	return uc.plugins.List()
}
