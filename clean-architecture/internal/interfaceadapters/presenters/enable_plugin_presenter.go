package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type EnablePluginViewModel struct {
	Message string
	Name    string
	Enabled bool
}

type EnablePluginPresenter struct {
	viewModel EnablePluginViewModel
}

func NewEnablePluginPresenter() *EnablePluginPresenter {
	return &EnablePluginPresenter{}
}

func (p *EnablePluginPresenter) Present(output usecases.EnablePluginOutput) error {
	p.viewModel = EnablePluginViewModel{
		Message: fmt.Sprintf("enabled plugin: name=%s enabled=%t", output.Name, output.Enabled),
		Name:    output.Name,
		Enabled: output.Enabled,
	}
	return nil
}

func (p *EnablePluginPresenter) ViewModel() EnablePluginViewModel {
	return p.viewModel
}
