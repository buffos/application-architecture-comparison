package plugins

import (
	"errors"
	"testing"
)

func TestRegisterEnableAndListPricingPlugin(t *testing.T) {
	c := NewComponent()
	registered, err := c.RegisterPricingPlugin(RegisterPricingPluginCommand{PluginID: "seasonal-pricing"})
	if err != nil {
		t.Fatal(err)
	}
	if registered.Enabled || registered.Type != PluginTypePricing {
		t.Fatalf("unexpected registration %+v", registered)
	}
	enabled, err := c.EnablePlugin(EnablePluginCommand{PluginID: "seasonal-pricing"})
	if err != nil {
		t.Fatal(err)
	}
	if !enabled.Enabled {
		t.Fatalf("expected enabled plugin, got %+v", enabled)
	}
	list := c.ListPlugins()
	if len(list) != 1 || !list[0].Enabled {
		t.Fatalf("unexpected plugin list %+v", list)
	}
}

func TestEnablePluginRequiresRegistration(t *testing.T) {
	_, err := NewComponent().EnablePlugin(EnablePluginCommand{PluginID: "missing"})
	if !errors.Is(err, ErrPluginNotFound) {
		t.Fatalf("got %v", err)
	}
}
