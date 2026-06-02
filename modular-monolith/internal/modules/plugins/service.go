package plugins

type RegisterPricingPluginCommand struct {
	PluginID string
}

type RegisterPricingPluginResult struct {
	PluginID string
	Type     string
	Enabled  bool
}

type EnablePluginCommand struct {
	PluginID string
}

type EnablePluginResult struct {
	PluginID string
	Type     string
	Enabled  bool
}

type ListPluginsQuery struct{}

type PluginDetails struct {
	PluginID string
	Type     string
	Enabled  bool
}

type Reader interface {
	IsEnabled(pluginID string) (bool, error)
}

type Service struct {
	plugins Repository
}

func NewService(plugins Repository) Service {
	return Service{plugins: plugins}
}

func (s Service) RegisterPricingPlugin(command RegisterPricingPluginCommand) (RegisterPricingPluginResult, error) {
	if _, err := s.plugins.FindByID(command.PluginID); err == nil {
		return RegisterPricingPluginResult{}, ErrPluginAlreadyExists
	}

	plugin := PluginRegistration{
		ID:      command.PluginID,
		Type:    PluginTypePricing,
		Enabled: false,
	}
	if err := s.plugins.Save(plugin); err != nil {
		return RegisterPricingPluginResult{}, err
	}

	return RegisterPricingPluginResult{
		PluginID: plugin.ID,
		Type:     plugin.Type,
		Enabled:  plugin.Enabled,
	}, nil
}

func (s Service) EnablePlugin(command EnablePluginCommand) (EnablePluginResult, error) {
	plugin, err := s.plugins.FindByID(command.PluginID)
	if err != nil {
		return EnablePluginResult{}, err
	}
	if plugin.Type != PluginTypePricing {
		return EnablePluginResult{}, ErrPluginNotPricingType
	}

	plugin.Enabled = true
	if err := s.plugins.Save(plugin); err != nil {
		return EnablePluginResult{}, err
	}

	return EnablePluginResult{
		PluginID: plugin.ID,
		Type:     plugin.Type,
		Enabled:  plugin.Enabled,
	}, nil
}

func (s Service) ListPlugins(query ListPluginsQuery) ([]PluginDetails, error) {
	plugins, err := s.plugins.List()
	if err != nil {
		return nil, err
	}

	list := make([]PluginDetails, 0, len(plugins))
	for _, plugin := range plugins {
		list = append(list, PluginDetails{
			PluginID: plugin.ID,
			Type:     plugin.Type,
			Enabled:  plugin.Enabled,
		})
	}

	return list, nil
}

func (s Service) IsEnabled(pluginID string) (bool, error) {
	plugin, err := s.plugins.FindByID(pluginID)
	if err != nil {
		return false, err
	}
	return plugin.Enabled, nil
}
