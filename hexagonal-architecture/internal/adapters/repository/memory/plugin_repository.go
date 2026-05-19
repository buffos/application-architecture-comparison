package memory

import (
	"sync"

	"hexagonal-architecture/internal/core/domain"
)

type PluginRepository struct {
	mu      sync.RWMutex
	plugins map[string]domain.PluginRegistration
}

func NewPluginRepository() *PluginRepository {
	return &PluginRepository{
		plugins: make(map[string]domain.PluginRegistration),
	}
}

func (r *PluginRepository) Save(plugin domain.PluginRegistration) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.plugins[plugin.Key] = plugin
	return nil
}

func (r *PluginRepository) FindByKey(key string) (domain.PluginRegistration, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugin, ok := r.plugins[key]
	if !ok {
		return domain.PluginRegistration{}, domain.ErrPluginNotFound
	}

	return plugin, nil
}

func (r *PluginRepository) List() ([]domain.PluginRegistration, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugins := make([]domain.PluginRegistration, 0, len(r.plugins))
	for _, plugin := range r.plugins {
		plugins = append(plugins, plugin)
	}

	return plugins, nil
}
