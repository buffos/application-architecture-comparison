package application

import "onion-architecture/internal/domain"

type RegisterPricingPluginCommand struct {
	Name string
}

type RegisterPricingPluginResult struct {
	Name    string
	Type    string
	Enabled bool
}

type RegisterPricingPluginService struct {
	plugins PluginRepository
}

func NewRegisterPricingPluginService(plugins PluginRepository) RegisterPricingPluginService {
	return RegisterPricingPluginService{plugins: plugins}
}

func (s RegisterPricingPluginService) Execute(command RegisterPricingPluginCommand) (RegisterPricingPluginResult, error) {
	plugin, err := domain.NewPluginRegistration(command.Name, "pricing")
	if err != nil {
		return RegisterPricingPluginResult{}, err
	}

	if err := s.plugins.Save(plugin); err != nil {
		return RegisterPricingPluginResult{}, err
	}

	return RegisterPricingPluginResult{
		Name:    plugin.Name,
		Type:    plugin.Type,
		Enabled: plugin.Enabled,
	}, nil
}
