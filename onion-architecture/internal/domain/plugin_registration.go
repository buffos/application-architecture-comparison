package domain

import "errors"

var ErrPluginNameRequired = errors.New("plugin name is required")
var ErrPluginNotFound = errors.New("plugin not found")

type PluginRegistration struct {
	Name    string
	Type    string
	Enabled bool
}

func NewPluginRegistration(name string, pluginType string) (PluginRegistration, error) {
	if name == "" {
		return PluginRegistration{}, ErrPluginNameRequired
	}

	return PluginRegistration{
		Name: name,
		Type: pluginType,
	}, nil
}

func (p *PluginRegistration) Enable() {
	p.Enabled = true
}
