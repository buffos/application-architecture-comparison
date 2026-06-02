package plugins

import "testing"

type stubRepository struct {
	plugins map[string]PluginRegistration
}

func (r *stubRepository) Save(plugin PluginRegistration) error {
	if r.plugins == nil {
		r.plugins = make(map[string]PluginRegistration)
	}
	r.plugins[plugin.ID] = plugin
	return nil
}

func (r *stubRepository) FindByID(id string) (PluginRegistration, error) {
	plugin, ok := r.plugins[id]
	if !ok {
		return PluginRegistration{}, ErrPluginNotFound
	}
	return plugin, nil
}

func (r *stubRepository) List() ([]PluginRegistration, error) {
	list := make([]PluginRegistration, 0, len(r.plugins))
	for _, plugin := range r.plugins {
		list = append(list, plugin)
	}
	return list, nil
}

func TestRegisterPricingPluginStoresDisabledPlugin(t *testing.T) {
	repository := &stubRepository{plugins: map[string]PluginRegistration{}}
	service := NewService(repository)

	result, err := service.RegisterPricingPlugin(RegisterPricingPluginCommand{PluginID: "seasonal-pricing"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Type != PluginTypePricing || result.Enabled {
		t.Fatalf("unexpected result %+v", result)
	}
}

func TestEnablePluginMarksPluginEnabled(t *testing.T) {
	repository := &stubRepository{plugins: map[string]PluginRegistration{
		"seasonal-pricing": {ID: "seasonal-pricing", Type: PluginTypePricing},
	}}
	service := NewService(repository)

	result, err := service.EnablePlugin(EnablePluginCommand{PluginID: "seasonal-pricing"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !result.Enabled {
		t.Fatalf("expected plugin enabled")
	}
}
