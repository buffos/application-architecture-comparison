package ports

import "hexagonal-architecture/internal/core/domain"

type PluginRepository interface {
	Save(plugin domain.PluginRegistration) error
	FindByKey(key string) (domain.PluginRegistration, error)
	List() ([]domain.PluginRegistration, error)
}
