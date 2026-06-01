package usecases

import "clean-architecture/internal/entities"

type PluginListItem struct {
	Name    string
	Enabled bool
}

type ListPluginsInput struct{}

type ListPluginsOutput struct {
	Count   int
	Plugins []PluginListItem
}

type ListPluginsInputBoundary interface {
	Execute(input ListPluginsInput) error
}

type ListPluginsOutputBoundary interface {
	Present(output ListPluginsOutput) error
}

type PluginLister interface {
	List() ([]entities.PluginRegistration, error)
}

type ListPluginsInteractor struct {
	plugins PluginLister
	output  ListPluginsOutputBoundary
}

func NewListPluginsInteractor(plugins PluginLister, output ListPluginsOutputBoundary) ListPluginsInteractor {
	return ListPluginsInteractor{plugins: plugins, output: output}
}

func (uc ListPluginsInteractor) Execute(input ListPluginsInput) error {
	_ = input
	plugins, err := uc.plugins.List()
	if err != nil {
		return err
	}

	items := make([]PluginListItem, 0, len(plugins))
	for _, plugin := range plugins {
		items = append(items, PluginListItem{Name: plugin.Name, Enabled: plugin.Enabled})
	}

	return uc.output.Present(ListPluginsOutput{Count: len(items), Plugins: items})
}
