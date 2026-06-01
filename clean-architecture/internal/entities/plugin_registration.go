package entities

import "errors"

var ErrPluginNameRequired = errors.New("plugin name is required")
var ErrPluginNotFound = errors.New("plugin not found")

type PluginRegistration struct {
	Name    string
	Enabled bool
}

func NewPluginRegistration(name string) (PluginRegistration, error) {
	if name == "" {
		return PluginRegistration{}, ErrPluginNameRequired
	}

	return PluginRegistration{Name: name}, nil
}

func (p *PluginRegistration) Enable() {
	p.Enabled = true
}
