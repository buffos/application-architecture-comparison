package application

type PricingPluginInput struct {
	SKU       string
	Category  string
	Quantity  int
	BasePrice int
}

type PricingAdjustment struct {
	Label         string
	AdjustedPrice int
}

type PricingPlugin interface {
	Key() string
	Adjust(input PricingPluginInput) (PricingAdjustment, bool)
}

type PricingPluginRegistry interface {
	EnabledPricingPlugins() []PricingPlugin
}
