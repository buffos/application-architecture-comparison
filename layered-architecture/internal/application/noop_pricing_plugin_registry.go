package application

type NoopPricingPluginRegistry struct{}

func (NoopPricingPluginRegistry) EnabledPricingPlugins() []PricingPlugin {
	return nil
}
