package application

type EnablePluginCommand struct {
	Name string
}

type EnablePluginResult struct {
	Name    string
	Type    string
	Enabled bool
}

type EnablePluginService struct {
	plugins PluginRepository
}

func NewEnablePluginService(plugins PluginRepository) EnablePluginService {
	return EnablePluginService{plugins: plugins}
}

func (s EnablePluginService) Execute(command EnablePluginCommand) (EnablePluginResult, error) {
	plugin, err := s.plugins.FindByName(command.Name)
	if err != nil {
		return EnablePluginResult{}, err
	}

	plugin.Enable()

	if err := s.plugins.Save(plugin); err != nil {
		return EnablePluginResult{}, err
	}

	return EnablePluginResult{
		Name:    plugin.Name,
		Type:    plugin.Type,
		Enabled: plugin.Enabled,
	}, nil
}
