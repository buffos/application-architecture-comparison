package application

type PluginDetails struct {
	Name    string
	Type    string
	Enabled bool
}

type ListPluginsService struct {
	plugins PluginRepository
}

func NewListPluginsService(plugins PluginRepository) ListPluginsService {
	return ListPluginsService{plugins: plugins}
}

func (s ListPluginsService) Execute() ([]PluginDetails, error) {
	plugins, err := s.plugins.List()
	if err != nil {
		return nil, err
	}

	result := make([]PluginDetails, 0, len(plugins))
	for _, plugin := range plugins {
		result = append(result, PluginDetails{
			Name:    plugin.Name,
			Type:    plugin.Type,
			Enabled: plugin.Enabled,
		})
	}

	return result, nil
}
