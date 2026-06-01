package application

import "onion-architecture/internal/domain"

type PluginRepository interface {
	Save(plugin domain.PluginRegistration) error
	FindByName(name string) (domain.PluginRegistration, error)
	List() ([]domain.PluginRegistration, error)
}
