package memory

import (
	"sync"

	"modular-monolith/internal/modules/plugins"
)

type PluginRepository struct {
	mu      sync.RWMutex
	plugins map[string]plugins.PluginRegistration
}

func NewPluginRepository() *PluginRepository {
	return &PluginRepository{
		plugins: make(map[string]plugins.PluginRegistration),
	}
}

func (r *PluginRepository) Save(plugin plugins.PluginRegistration) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.plugins[plugin.ID] = plugin
	return nil
}

func (r *PluginRepository) FindByID(id string) (plugins.PluginRegistration, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugin, ok := r.plugins[id]
	if !ok {
		return plugins.PluginRegistration{}, plugins.ErrPluginNotFound
	}
	return plugin, nil
}

func (r *PluginRepository) List() ([]plugins.PluginRegistration, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]plugins.PluginRegistration, 0, len(r.plugins))
	for _, plugin := range r.plugins {
		list = append(list, plugin)
	}
	return list, nil
}
