package plugins

import "errors"

const PluginTypePricing = "pricing"

var (
	ErrPluginIDRequired     = errors.New("plugin id is required")
	ErrPluginAlreadyExists  = errors.New("plugin already exists")
	ErrPluginNotFound       = errors.New("plugin not found")
	ErrPluginNotPricingType = errors.New("plugin is not a pricing plugin")
)

type Plugin struct {
	ID      string
	Type    string
	Enabled bool
}

type PluginDetails struct {
	PluginID string
	Type     string
	Enabled  bool
}

type RegisterPricingPluginCommand struct{ PluginID string }
type RegisterPricingPluginResult struct {
	PluginID string
	Type     string
	Enabled  bool
}
type EnablePluginCommand struct{ PluginID string }
type EnablePluginResult struct {
	PluginID string
	Type     string
	Enabled  bool
}

type Reader interface {
	IsEnabled(pluginID string) (bool, error)
	ListPlugins() []PluginDetails
}

type Component struct {
	plugins map[string]Plugin
}

func NewComponent() *Component { return &Component{plugins: map[string]Plugin{}} }

func (c *Component) RegisterPricingPlugin(command RegisterPricingPluginCommand) (RegisterPricingPluginResult, error) {
	if command.PluginID == "" {
		return RegisterPricingPluginResult{}, ErrPluginIDRequired
	}
	if _, exists := c.plugins[command.PluginID]; exists {
		return RegisterPricingPluginResult{}, ErrPluginAlreadyExists
	}
	plugin := Plugin{ID: command.PluginID, Type: PluginTypePricing}
	c.plugins[plugin.ID] = plugin
	return RegisterPricingPluginResult{PluginID: plugin.ID, Type: plugin.Type, Enabled: plugin.Enabled}, nil
}

func (c *Component) EnablePlugin(command EnablePluginCommand) (EnablePluginResult, error) {
	plugin, ok := c.plugins[command.PluginID]
	if !ok {
		return EnablePluginResult{}, ErrPluginNotFound
	}
	if plugin.Type != PluginTypePricing {
		return EnablePluginResult{}, ErrPluginNotPricingType
	}
	plugin.Enabled = true
	c.plugins[plugin.ID] = plugin
	return EnablePluginResult{PluginID: plugin.ID, Type: plugin.Type, Enabled: plugin.Enabled}, nil
}

func (c *Component) IsEnabled(pluginID string) (bool, error) {
	plugin, ok := c.plugins[pluginID]
	if !ok {
		return false, ErrPluginNotFound
	}
	return plugin.Enabled, nil
}

func (c *Component) ListPlugins() []PluginDetails {
	list := make([]PluginDetails, 0, len(c.plugins))
	for _, plugin := range c.plugins {
		list = append(list, PluginDetails{PluginID: plugin.ID, Type: plugin.Type, Enabled: plugin.Enabled})
	}
	return list
}

var _ Reader = (*Component)(nil)
