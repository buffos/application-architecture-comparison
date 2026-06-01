package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type PluginListItemViewModel struct {
	Name    string
	Enabled bool
}

type ListPluginsViewModel struct {
	Message string
	Count   int
	Plugins []PluginListItemViewModel
}

type ListPluginsPresenter struct {
	viewModel ListPluginsViewModel
}

func NewListPluginsPresenter() *ListPluginsPresenter {
	return &ListPluginsPresenter{}
}

func (p *ListPluginsPresenter) Present(output usecases.ListPluginsOutput) error {
	items := make([]PluginListItemViewModel, 0, len(output.Plugins))
	for _, plugin := range output.Plugins {
		items = append(items, PluginListItemViewModel{Name: plugin.Name, Enabled: plugin.Enabled})
	}

	p.viewModel = ListPluginsViewModel{
		Message: fmt.Sprintf("listed plugins: count=%d", output.Count),
		Count:   output.Count,
		Plugins: items,
	}
	return nil
}

func (p *ListPluginsPresenter) ViewModel() ListPluginsViewModel {
	return p.viewModel
}
