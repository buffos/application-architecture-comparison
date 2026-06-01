package memory

import (
	"sort"
	"sync"

	"clean-architecture/internal/entities"
)

type PluginGateway struct {
	mu      sync.RWMutex
	plugins map[string]entities.PluginRegistration
}

func NewPluginGateway() *PluginGateway {
	return &PluginGateway{
		plugins: make(map[string]entities.PluginRegistration),
	}
}

func (g *PluginGateway) Save(plugin entities.PluginRegistration) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.plugins[plugin.Name] = plugin
	return nil
}

func (g *PluginGateway) FindByName(name string) (entities.PluginRegistration, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	plugin, ok := g.plugins[name]
	if !ok {
		return entities.PluginRegistration{}, entities.ErrPluginNotFound
	}

	return plugin, nil
}

func (g *PluginGateway) List() ([]entities.PluginRegistration, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	plugins := make([]entities.PluginRegistration, 0, len(g.plugins))
	for _, plugin := range g.plugins {
		plugins = append(plugins, plugin)
	}

	sort.Slice(plugins, func(i int, j int) bool {
		return plugins[i].Name < plugins[j].Name
	})

	return plugins, nil
}
