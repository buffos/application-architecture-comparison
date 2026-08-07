package scripts

import (
	"errors"
	"sort"

	"transaction-script-architecture/internal/data"
)

const PluginTypePricing = "pricing"

var (
	ErrPluginKeyRequired  = errors.New("plugin key is required")
	ErrPluginTypeRequired = errors.New("plugin type is required")
	ErrPluginNotFound     = errors.New("plugin not found")
)

func RegisterPlugin(store *data.Store, key string, pluginType string, version string, config map[string]string) (data.PluginRegistration, error) {
	if store == nil {
		return data.PluginRegistration{}, ErrStoreRequired
	}
	if key == "" {
		return data.PluginRegistration{}, ErrPluginKeyRequired
	}
	if pluginType == "" {
		return data.PluginRegistration{}, ErrPluginTypeRequired
	}

	registration := data.PluginRegistration{
		Key:     key,
		Type:    pluginType,
		Version: version,
		Config:  clonePluginConfig(config),
	}
	store.Plugins[key] = registration
	return registration, nil
}

func EnablePlugin(store *data.Store, key string) (data.PluginRegistration, error) {
	return setPluginEnabled(store, key, true)
}

func DisablePlugin(store *data.Store, key string) (data.PluginRegistration, error) {
	return setPluginEnabled(store, key, false)
}

func setPluginEnabled(store *data.Store, key string, enabled bool) (data.PluginRegistration, error) {
	if store == nil {
		return data.PluginRegistration{}, ErrStoreRequired
	}
	if key == "" {
		return data.PluginRegistration{}, ErrPluginKeyRequired
	}

	plugin, ok := store.Plugins[key]
	if !ok {
		return data.PluginRegistration{}, ErrPluginNotFound
	}
	plugin.Enabled = enabled
	store.Plugins[key] = plugin
	return plugin, nil
}

func ListPlugins(store *data.Store) ([]data.PluginRegistration, error) {
	if store == nil {
		return nil, ErrStoreRequired
	}

	plugins := make([]data.PluginRegistration, 0, len(store.Plugins))
	for _, plugin := range store.Plugins {
		plugin.Config = clonePluginConfig(plugin.Config)
		plugins = append(plugins, plugin)
	}
	sort.Slice(plugins, func(i, j int) bool {
		return plugins[i].Key < plugins[j].Key
	})
	return plugins, nil
}

func clonePluginConfig(config map[string]string) map[string]string {
	if config == nil {
		return nil
	}
	clone := make(map[string]string, len(config))
	for key, value := range config {
		clone[key] = value
	}
	return clone
}
