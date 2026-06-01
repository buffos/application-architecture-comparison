package memory

import (
	"sync"

	"onion-architecture/internal/domain"
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

	r.plugins[plugin.Name] = plugin
	return nil
}

func (r *PluginRepository) FindByName(name string) (domain.PluginRegistration, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugin, ok := r.plugins[name]
	if !ok {
		return domain.PluginRegistration{}, domain.ErrPluginNotFound
	}

	return plugin, nil
}

func (r *PluginRepository) List() ([]domain.PluginRegistration, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]domain.PluginRegistration, 0, len(r.plugins))
	for _, plugin := range r.plugins {
		result = append(result, plugin)
	}

	return result, nil
}
