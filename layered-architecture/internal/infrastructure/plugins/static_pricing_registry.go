package plugins

import "layered-architecture/internal/application"

type StaticPricingPluginRegistry struct {
	plugins []application.PricingPlugin
}

func NewStaticPricingPluginRegistry(plugins []application.PricingPlugin) StaticPricingPluginRegistry {
	return StaticPricingPluginRegistry{plugins: plugins}
}

func (r StaticPricingPluginRegistry) EnabledPricingPlugins() []application.PricingPlugin {
	return r.plugins
}
