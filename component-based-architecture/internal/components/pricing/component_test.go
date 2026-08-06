package pricing

import (
	"testing"

	"component-based-architecture/internal/components/plugins"
	"component-based-architecture/internal/components/products"
)

func TestUnitPriceAppliesSeasonalPricingWhenPluginEnabled(t *testing.T) {
	pluginComponent := plugins.NewComponent()
	if _, err := pluginComponent.RegisterPricingPlugin(plugins.RegisterPricingPluginCommand{PluginID: "seasonal-pricing"}); err != nil {
		t.Fatal(err)
	}
	if _, err := pluginComponent.EnablePlugin(plugins.EnablePluginCommand{PluginID: "seasonal-pricing"}); err != nil {
		t.Fatal(err)
	}
	price, err := NewComponent(pluginComponent).UnitPrice(products.ProductForQuote{UnitPrice: 15000})
	if err != nil {
		t.Fatal(err)
	}
	if price != 14250 {
		t.Fatalf("price = %d, want 14250", price)
	}
}

func TestUnitPriceUsesProductPriceWhenPluginDisabled(t *testing.T) {
	price, err := NewComponent(plugins.NewComponent()).UnitPrice(products.ProductForQuote{UnitPrice: 15000})
	if err != nil {
		t.Fatal(err)
	}
	if price != 15000 {
		t.Fatalf("price = %d, want 15000", price)
	}
}
